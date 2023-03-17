# Postgres

The postgres container image provides a preconfigured PostgreSQL database server for easy deployment in containerized environments.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - env:
    - POSTGRES_PASSWORD=password
    image: postgres
    name: postgres
    ports: "5432"
    volumeMounts:
    - mountPath: /var/lib/postgresql/data
      name: postgres.data
  volumes:
  - hostPath:
      path: volumes/postgres/data
    name: postgres.data
```

