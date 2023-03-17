# Memcached

The memcached container image provides an efficient in-memory caching system for key-value pairs.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: memcached
    name: memcached
    ports: "11211"
```

