# Redis

The Redis container image is a lightweight, open-source, and in-memory data structure store used as a cache, database, and message broker.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: redis
    name: redis
    ports: "6379"
    volumeMounts:
    - mountPath: /data
      name: redis.data
  volumes:
  - hostPath:
      path: volumes/redis/data
    name: redis.data
```

