# Influxdb

The influxdb container image provides a pre-configured instance of the InfluxDB database that can be easily deployed for time-series data storage and analysis.

```yaml
tasks:
  "":
    image: influxdb
volumes:
- hostPath:
    path: volumes/influxdb/influxdb2
  name: influxdb.influxdb2
- hostPath:
    path: volumes/influxdb/influxdb2
  name: influxdb.influxdb2
```

