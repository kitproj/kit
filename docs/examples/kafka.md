# Kafka  with Kraft

[Help](https://github.com/kitproj/kafka-docker)

Standalone Kafka image that uses Kraft.

```yaml
tasks:
  "":
    image: ghcr.io/kitproj/kafka
volumes:
- hostPath:
    path: volumes/kafka/kraft-combined-logs
  name: kafka.kraft-combined-logs
```

Licence(s): MIT

