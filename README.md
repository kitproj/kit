# Kit

[![CodeQL](https://github.com/kitproj/kit/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/codeql-analysis.yml)
[![Go](https://github.com/kitproj/kit/actions/workflows/go.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/go.yml)
[![goreleaser](https://github.com/kitproj/kit/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/goreleaser.yml)

Kit is a workflow engine for software development.

Kit combines both task execution (like Foreman), container management (like Docker Compose), Kubernetes resource
management (
like Tilt, Skaffold), and a focus on local development (like Garden) in a single, easy-to-use binary.

It works seamlessly with both local-dev and cloud-dev environments, such as Codespaces and Gitpod.

## Install

Like `jq`, `kit` is a small standalone binary. You can download it from
the [releases page](https://github.com/kitproj/kit/releases/latest).

If you're on MacOS, you can use `brew`:

```bash
brew tap kitproj/kit --custom-remote https://github.com/kitproj/kit
brew install kit
```

Otherwise, you can use `curl`:

```bash
curl -q https://raw.githubusercontent.com/kitproj/kit/main/install.sh | sh
```

## Usage

Workflows are described by a directed acyclic graph (DAG) of tasks.

Create a [`tasks.yaml`](tasks.yaml) file, e.g.:

```yaml
apiVersion: kit/v1
spec:
  tasks:
    - name: build
      command: go build .
```

Start:

```bash
kit build
```

### Services

A task can be a **service** by specifying its ports:

```yaml
- name: service
  command: go run .
  ports: [ 8080 ]
```

The ports will be forwarded from the host to the service. A service will be restarted if it does not start-up (i.e. it is listening on the port).

### Dependencies

Tasks can depend on other tasks:

```yaml
- name: build
  command: go build .
- name: test
  command: go test .
  dependencies: [ build ]
```

Tasks will only be started if the dependencies have completed successfully, or if the task is a service, it is running and listening on its port.

### Tasks

#### Host Task

A **host task** runs on the host machine. It is defined by a `command`:

```yaml
- name: build
  command: go build .
```

Once a task completes successfully, any downstream tasks are started. If it is unsuccessful, Kit will exit

Unlike a plain task, if a service does not start-up (i.e. it is listening on the port), it will be restarted. You can
specify a probe to determine if the  service is running correctly:

```yaml
- name: service
  command: go run .
  ports: [ 8080 ]
  readinessProbe:
    httpGet:
      path: /healthz 
```

### Shell Task

A **shell task** is just a host task that runs in a shell:

```yaml
- name: shell
  sh: |
    set -eux
    echo "Hello, world!"
```

### Container Task

A **container task** runs in a container. It is defined by an `image`:

```yaml
- name: mysql
  image: mysql
  ports: [ 3306:3306 ]
```

The ports will be forwarded from the host to the container.

If the image is a path to a directory containing Dockerfile, it will be built and run automatically:

```yaml
- name: kafka
  image: ./src/images/kafka
```

### Kubernetes Task

A **Kubernetes task** deploys manifests to a Kubernetes cluster, it is defined by `manifests`:

```yaml
- name: deploy
  namespace: default
  manifests:
    - manifests/
    - service.yaml
  ports: [ 80:8080 ]
```

The ports will be forwarded from the Kubernetes cluster to the host.

### No-op Task

A **no-op task** is a task that does nothing, commonly a task named `up` is provided and that depends on all other
tasks:

```yaml
- name: up
  dependencies: [ deploy ]
```

### Environment Variables

Task can have **environment variables**:

```yaml
- name: foo
  command: go run .
  env:
    # simple key-value pair
    - FOO=1
    # value from a file
    - name: BAR
      valueFrom:
        file: bar.txt
  # environment variables from a file
  envfile: .env
```

### Watches

Task can be **automatically re-run** when a file changes:

```yaml
- name: build
  command: go build .
  watch: src/
```

### Targets

If a task produces an output, you can avoid repeating work by specifying the **task target**:

```yaml
- name: build
  command: go build .
  target: bin/app
```

The task will be skipped if the target is newer that the watched sources (just like Make).

### Mutexes and Semaphores

Task can have **mutexes** and **semaphores** to control concurrency:

If you want to prevent two tasks from running at the same time, you can use a mutex:

```yaml
tasks:
  - name: foo
    mutex: my-mutex
  - name: bar
    mutex: my-mutex
```

If you want to limit the number of tasks that can run at the same time, you can use a semaphore:

```yaml
# only two can run at the same time
semaphores:
  my-semaphore: 2
tasks:
  - name: foo
    semaphore: my-semaphore
  - name: bar
    semaphore: my-semaphore
```

### Logging

Sometimes a task logs too much, you can send logs to a file:

```yaml 
- name: build
  command: go build .
  log: build.log
```

## Documentation

- [Usage](docs/USAGE.md) - how to use the various features of kit
- [Examples](docs/examples) - examples of how to use kit, e.g. with MySQL, or Kafka
- [Reference](docs/reference) - reference documentation for the various types in kit

