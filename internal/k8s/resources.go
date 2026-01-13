package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type AppSpec struct {
	Name            string
	Slug            string
	Image           string
	Port            int32
	Replicas        int32
	CPULimit        string
	MemoryLimit     string
	Domain          string
	HealthCheckPath string
}

// BuildDeployment creates a Deployment manifest for an app
func BuildDeployment(spec AppSpec) *appsv1.Deployment {
	labels := map[string]string{
		"app":              spec.Slug,
		"superfly.dev/app": spec.Slug,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Slug,
			Namespace: AppsNamespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						// Prometheus scraping annotations
						"prometheus.io/scrape": "true",
						"prometheus.io/port":   "9090",
						"prometheus.io/path":   "/metrics",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "app",
							Image: spec.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: spec.Port,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(spec.CPULimit),
									corev1.ResourceMemory: resource.MustParse(spec.MemoryLimit),
								},
								Requests: corev1.ResourceList{
									// Requests are 50% of limits
									corev1.ResourceCPU:    resource.MustParse(halveResource(spec.CPULimit)),
									corev1.ResourceMemory: resource.MustParse(halveResource(spec.MemoryLimit)),
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: spec.HealthCheckPath,
										Port: intstr.FromInt32(spec.Port),
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								TimeoutSeconds:      5,
								FailureThreshold:    3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: spec.HealthCheckPath,
										Port: intstr.FromInt32(spec.Port),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
								TimeoutSeconds:      3,
								FailureThreshold:    3,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: 0},
					MaxSurge:       &intstr.IntOrString{Type: intstr.Int, IntVal: 1},
				},
			},
		},
	}
}

// BuildService creates a Service manifest for an app
func BuildService(spec AppSpec) *corev1.Service {
	labels := map[string]string{
		"app":              spec.Slug,
		"superfly.dev/app": spec.Slug,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Slug,
			Namespace: AppsNamespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt32(spec.Port),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

// BuildIngress creates an Ingress manifest for an app
func BuildIngress(spec AppSpec) *networkingv1.Ingress {
	labels := map[string]string{
		"app":              spec.Slug,
		"superfly.dev/app": spec.Slug,
	}

	pathTypePrefix := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Slug,
			Namespace: AppsNamespace,
			Labels:    labels,
			Annotations: map[string]string{
				"cert-manager.io/cluster-issuer":          "letsencrypt-prod",
				"traefik.ingress.kubernetes.io/router.tls": "true",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: spec.Domain,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathTypePrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: spec.Slug,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Add TLS configuration if domain is provided
	if spec.Domain != "" {
		ingress.Spec.TLS = []networkingv1.IngressTLS{
			{
				Hosts:      []string{spec.Domain},
				SecretName: spec.Slug + "-tls",
			},
		}
	}

	return ingress
}

// halveResource returns half of a resource string (e.g., "1000m" -> "500m")
func halveResource(resourceStr string) string {
	q := resource.MustParse(resourceStr)
	q.Sub(q)
	q.Add(resource.MustParse(resourceStr))
	
	// Get value in milli units
	milliValue := q.MilliValue()
	halved := milliValue / 2
	
	// Return as string
	if resourceStr[len(resourceStr)-1] == 'i' {
		// Memory (e.g., "512Mi")
		return resource.NewQuantity(halved/(1024*1024), resource.BinarySI).String()
	}
	// CPU (e.g., "500m")
	return resource.NewMilliQuantity(halved, resource.DecimalSI).String()
}
