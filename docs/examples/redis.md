# Redis

The Redis container image is a lightweight, open-source, and in-memory data structure store used as a cache, database, and message broker.

```yaml
tasks:
  "":
    image: redis
volumes:
- hostPath:
    path: volumes/redis/data
  name: redis.data
```

