# Pulsar

The apachepulsar/pulsar container image provides a scalable and distributed messaging system for enterprise-grade applications.

```yaml
tasks:
  "":
    command: /pulsar/bin/pulsar standalone
    image: apachepulsar/pulsar
    ports: 6650 8080
```

