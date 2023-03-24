# Timestream

Timestream is a fully-managed, scalable time-series database service developed by Amazon Web Services (AWS) that enables users to store, process, and analyze time-series data, such as logs, sensor data, and industrial telemetry, with high precision and accuracy, and built-in analytics and visualization tools.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: motoserver/moto
    name: timestream
    ports: "5000"
```

