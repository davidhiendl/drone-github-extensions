# Deploy via Helm

## Example values.yaml

A minimal values file for helm with only the absolutely required variables.

```yaml
config:
  DRONE_SECRET: "xxxxxxxx"

ingress:
  enabled: true
  hosts:
    - host: github-drone-extensions.example.com
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls:
    - secretName: ingress-tls-github-drone-extensions
      hosts:
        - github-drone-extensions.example.com
```
