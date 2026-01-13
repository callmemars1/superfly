# DNS Setup for Superfly

Your Superfly installation: **superfly.smartynov.com**

## DNS Configuration Required

To access your apps via their domains, you need to configure DNS records.

### Your Server Details

**Domain**: `superfly.smartynov.com`  
**Server IP**: [YOUR_SERVER_IP_HERE]

---

## DNS Records Setup

### 1. Main Superfly Dashboard (Optional)

If you want to access Superfly's dashboard at `superfly.smartynov.com`:

```
Type: A
Name: superfly
Value: YOUR_SERVER_IP
TTL: 3600
```

### 2. Wildcard for Apps (Recommended)

To allow any subdomain to work (e.g., `app1.superfly.smartynov.com`, `api.superfly.smartynov.com`):

```
Type: A
Name: *.superfly
Value: YOUR_SERVER_IP
TTL: 3600
```

This allows you to deploy apps to any subdomain without adding new DNS records each time.

### 3. Individual App Records (Alternative)

If you don't use wildcard, add a record for each app:

```
Type: A
Name: app1.superfly
Value: YOUR_SERVER_IP
TTL: 3600

Type: A
Name: api.superfly
Value: YOUR_SERVER_IP
TTL: 3600
```

---

## Example DNS Configuration

Assuming your server IP is `123.45.67.89`:

| Type | Name           | Value        | TTL  |
|------|----------------|--------------|------|
| A    | superfly       | 123.45.67.89 | 3600 |
| A    | *.superfly     | 123.45.67.89 | 3600 |

**Result**: All of these will work:
- `https://superfly.smartynov.com` (main domain)
- `https://app1.superfly.smartynov.com` (any app)
- `https://api.superfly.smartynov.com` (any app)
- `https://anything.superfly.smartynov.com` (any app)

---

## Verification

### Check DNS Propagation

```bash
# Check if DNS is resolving
dig superfly.smartynov.com

# Should return your server IP
dig app1.superfly.smartynov.com

# Check from external DNS
nslookup superfly.smartynov.com 8.8.8.8
```

### Test Connectivity

```bash
# Ping your domain
ping superfly.smartynov.com

# Test HTTP (before deploying apps)
curl http://superfly.smartynov.com
```

---

## Deploy Your First App

Once DNS is configured, deploy an app:

```bash
curl -X POST http://localhost:8080/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My First App",
    "image": "nginx:alpine",
    "port": 80,
    "domain": "app.superfly.smartynov.com"
  }'
```

Wait 30-60 seconds for:
1. Pod to start
2. Ingress to be configured
3. Let's Encrypt certificate to be issued

Then access:
```bash
curl https://app.superfly.smartynov.com
```

---

## Common DNS Providers

### Cloudflare

1. Go to your domain's DNS settings
2. Add A record:
   - Type: `A`
   - Name: `superfly`
   - IPv4 address: `YOUR_SERVER_IP`
   - Proxy status: DNS only (gray cloud)
   - TTL: Auto
3. Add wildcard A record:
   - Type: `A`
   - Name: `*.superfly`
   - IPv4 address: `YOUR_SERVER_IP`
   - Proxy status: DNS only (gray cloud)
   - TTL: Auto

**Important**: Use "DNS only" (gray cloud), not "Proxied" (orange cloud), for Let's Encrypt to work properly.

### DigitalOcean

1. Go to Networking → Domains → smartynov.com
2. Add A record:
   - Hostname: `superfly`
   - Will Direct To: `YOUR_SERVER_IP`
   - TTL: 3600
3. Add wildcard:
   - Hostname: `*.superfly`
   - Will Direct To: `YOUR_SERVER_IP`
   - TTL: 3600

### Namecheap

1. Go to Advanced DNS
2. Add New Record:
   - Type: `A Record`
   - Host: `superfly`
   - Value: `YOUR_SERVER_IP`
   - TTL: Automatic
3. Add wildcard:
   - Type: `A Record`
   - Host: `*.superfly`
   - Value: `YOUR_SERVER_IP`
   - TTL: Automatic

### AWS Route 53

1. Go to Hosted Zones → smartynov.com
2. Create Record Set:
   - Name: `superfly.smartynov.com`
   - Type: `A - IPv4 address`
   - Value: `YOUR_SERVER_IP`
   - TTL: 300
3. Create wildcard:
   - Name: `*.superfly.smartynov.com`
   - Type: `A - IPv4 address`
   - Value: `YOUR_SERVER_IP`
   - TTL: 300

---

## DNS Propagation Time

- **Typical**: 5-30 minutes
- **Maximum**: Up to 48 hours (rare)

Check propagation status:
- https://www.whatsmydns.net/#A/superfly.smartynov.com
- https://dnschecker.org/

---

## Troubleshooting

### DNS not resolving

```bash
# Check if DNS record exists
dig superfly.smartynov.com

# If empty, DNS not configured yet
# If shows IP, DNS is working
```

### Certificate not issuing

```bash
# Check certificate status
kubectl get certificate -n superfly-apps

# Describe certificate
kubectl describe certificate APP-NAME-tls -n superfly-apps

# Check cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager -f
```

**Common issues**:
- Ports 80/443 not accessible from internet
- Firewall blocking traffic
- DNS not resolving to correct IP
- Cloudflare proxy enabled (should be DNS only)

### App not accessible

```bash
# Check ingress
kubectl get ingress -n superfly-apps

# Describe ingress
kubectl describe ingress APP-NAME -n superfly-apps

# Check Traefik logs
kubectl logs -n traefik deployment/traefik -f
```

---

## Example App Domains

With your setup, you can use:

- `blog.superfly.smartynov.com` - Personal blog
- `api.superfly.smartynov.com` - Backend API
- `app.superfly.smartynov.com` - Web application
- `db.superfly.smartynov.com` - Database (internal only)
- `cache.superfly.smartynov.com` - Redis cache
- `test.superfly.smartynov.com` - Testing environment

All automatically get HTTPS! 🎉

---

## Production Checklist

Before deploying production apps:

- [ ] DNS configured (A records)
- [ ] DNS propagation complete (check with dig)
- [ ] Ports 80/443 open on firewall
- [ ] Server accessible from internet
- [ ] cert-manager working (check with test app)
- [ ] Let's Encrypt production issuer configured
- [ ] SSL certificates issuing successfully

---

## Next Steps

1. Configure DNS records above
2. Wait for propagation (5-30 mins)
3. Verify with `dig superfly.smartynov.com`
4. Deploy your first app
5. Check certificate status
6. Access via HTTPS!

**Need help?** Check [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) for troubleshooting.
