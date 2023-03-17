# Nats

The nats container image provides a lightweight messaging system for distributed applications running in containerized environments.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: nats
    name: nats
    ports: 4222 6222 8222
```

