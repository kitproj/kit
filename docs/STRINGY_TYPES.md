# Stringy Types

This project uses **stringy types** to reduce the verbosity of the YAML configuration. A stringy type is any type that can be represented as a string. For example, a `command` can be represented as a string, or as an object or array.

Example:

```yaml
- command:
    - go
    - run
    - ./demo/foo
  name: foo
```

In this example, `command` can be written as a space-separated string:

```yaml
- command: go run ./demo/foo
  name: foo
```

If the command contains spaces, use quotes:

```yaml
- command:
    - sh
    - -c
    - echo 1
  name: foo
```

Can be written as

```yaml
- command: sh -c  "echo 1"
  name: foo
```

## Task

```yaml
- name: foo
  command: 
    - go
  args:
    - build
    - foo
  watch:
    - foo
  dependencies:
    - bar
```

```yaml
- name: foo
  command: go build foo
  watch: foo
  dependencies: bar
```

## EnvVars

```yaml
- name: FOO
  value: 1
- name: BAR
  value: 2
```

```
- FOO=1
- BAR=2
```

## Ports

```yaml
- containerPort: 8080
  hostPort: 8081
- containerPort: 9090
  hostPort: 9091
```

```
- 8080:8081 9090:9091
```

## Probe

```yaml
- http:
    - proto: HTTPS
    - port: 8080
  initialDelaySeconds: 5
```

```
- https://:8080?initialDelaySeconds=5
```