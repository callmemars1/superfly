package k8s

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	AppsNamespace = "superfly-apps"
)

type Client struct {
	clientset *kubernetes.Clientset
}

// NewClient creates a new Kubernetes client
func NewClient(inCluster bool, kubeconfig string) (*Client, error) {
	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s clientset: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

// EnsureNamespace creates the apps namespace if it doesn't exist
func (c *Client) EnsureNamespace(ctx context.Context) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: AppsNamespace,
		},
	}

	_, err := c.clientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	return nil
}

// ApplyDeployment creates or updates a Deployment
func (c *Client) ApplyDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	deploymentsClient := c.clientset.AppsV1().Deployments(AppsNamespace)

	existing, err := deploymentsClient.Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new deployment
			_, err = deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create deployment: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// Update existing deployment
	deployment.ResourceVersion = existing.ResourceVersion
	_, err = deploymentsClient.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

// ApplyService creates or updates a Service
func (c *Client) ApplyService(ctx context.Context, service *corev1.Service) error {
	servicesClient := c.clientset.CoreV1().Services(AppsNamespace)

	existing, err := servicesClient.Get(ctx, service.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new service
			_, err = servicesClient.Create(ctx, service, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get service: %w", err)
	}

	// Update existing service (preserve ClusterIP)
	service.ResourceVersion = existing.ResourceVersion
	service.Spec.ClusterIP = existing.Spec.ClusterIP
	_, err = servicesClient.Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	return nil
}

// ApplyIngress creates or updates an Ingress
func (c *Client) ApplyIngress(ctx context.Context, ingress *networkingv1.Ingress) error {
	ingressClient := c.clientset.NetworkingV1().Ingresses(AppsNamespace)

	existing, err := ingressClient.Get(ctx, ingress.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new ingress
			_, err = ingressClient.Create(ctx, ingress, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create ingress: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get ingress: %w", err)
	}

	// Update existing ingress
	ingress.ResourceVersion = existing.ResourceVersion
	_, err = ingressClient.Update(ctx, ingress, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ingress: %w", err)
	}

	return nil
}

// DeleteDeployment deletes a Deployment
func (c *Client) DeleteDeployment(ctx context.Context, name string) error {
	err := c.clientset.AppsV1().Deployments(AppsNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}
	return nil
}

// DeleteService deletes a Service
func (c *Client) DeleteService(ctx context.Context, name string) error {
	err := c.clientset.CoreV1().Services(AppsNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}

// DeleteIngress deletes an Ingress
func (c *Client) DeleteIngress(ctx context.Context, name string) error {
	err := c.clientset.NetworkingV1().Ingresses(AppsNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete ingress: %w", err)
	}
	return nil
}

// GetDeploymentStatus gets the status of a deployment
func (c *Client) GetDeploymentStatus(ctx context.Context, name string) (*appsv1.DeploymentStatus, error) {
	deployment, err := c.clientset.AppsV1().Deployments(AppsNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return &deployment.Status, nil
}

// WaitForDeployment waits for a deployment to be ready
func (c *Client) WaitForDeployment(ctx context.Context, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		status, err := c.GetDeploymentStatus(ctx, name)
		if err != nil {
			return err
		}

		// Check if deployment is ready
		if status.ReadyReplicas > 0 && status.ReadyReplicas == status.Replicas {
			return nil
		}

		// Wait before checking again
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			// Continue loop
		}
	}

	return fmt.Errorf("timeout waiting for deployment to be ready")
}

// RestartDeployment restarts a deployment by updating its restart annotation
func (c *Client) RestartDeployment(ctx context.Context, name string) error {
	deployment, err := c.clientset.AppsV1().Deployments(AppsNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["superfly.dev/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = c.clientset.AppsV1().Deployments(AppsNamespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to restart deployment: %w", err)
	}

	return nil
}
