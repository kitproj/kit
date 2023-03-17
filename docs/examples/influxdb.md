# Influxdb

The influxdb container image provides a pre-configured instance of the InfluxDB database that can be easily deployed for time-series data storage and analysis.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: influxdb
    name: influxdb
    ports: "8086"
    volumeMounts:
    - mountPath: /etc/influxdb2
      name: influxdb.influxdb2
    - mountPath: /var/lib/influxdb2
      name: influxdb.influxdb2
  volumes:
  - hostPath:
      path: volumes/influxdb/influxdb2
    name: influxdb.influxdb2
  - hostPath:
      path: volumes/influxdb/influxdb2
    name: influxdb.influxdb2
```

