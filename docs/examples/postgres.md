# Postgres

The postgres container image provides a preconfigured PostgreSQL database server for easy deployment in containerized environments.

```yaml
tasks:
  "":
    env:
      POSTGRES_PASSWORD: password
    image: postgres
volumes:
- hostPath:
    path: volumes/postgres/data
  name: postgres.data
```

