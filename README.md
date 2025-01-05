# Kit

[![CodeQL](https://github.com/kitproj/kit/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/codeql-analysis.yml)
[![Go](https://github.com/kitproj/kit/actions/workflows/go.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/go.yml)
[![goreleaser](https://github.com/kitproj/kit/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/goreleaser.yml)

## Why

Make the dev loop crazy fast.

## What

Kit is a software development tool designed to turbo-charge the software development process.

Kit combines both task execution (like Foreman), container management (like Docker Compose), Kubernetes-like features (
like Tilt, Skaffold), and a focus on local development (like Garden) in a single, easy-to-use binary.

It is designed to work seamlessly with both local-dev and cloud-dev environments, such as Codespaces and Gitpod.

[![Watch the video](https://img.youtube.com/vi/IafQwT1rYOU/hqdefault.jpg)](https://youtu.be/IafQwT1rYOU)

Key features of Kit include:

* **Local testing**: Kit is designed for local testing, allowing developers to test their code on their local machines
  before pushing it to a test environment or production. This speeds up the testing process and helps developers catch
  and fix bugs more quickly.
* **First-class container support**: Kit downloads and runs containers, allowing developers to test their code in a
  containerized environment without having to set up a separate container management system.
* **First-class Kubernetes support**: Kit can deploy manifests to a cluster and automatically port-forward to the
  cluster.
* **Advanced DAG architecture**: Kit's directed acyclic graph (DAG) structure allows for optimized parallel processing,
  reducing the time required for testing and speeding up the development process.
* **Probes**: You can specify liveness probes for your tasks to see if they're working, automatically restarting them
  when they go wrong. You can also specify readiness probes for your tasks to see if they're ready.
* **Dependency management**: You can specify dependencies between tasks, so when upstream tasks become successful or
  ready, downstream tasks are automatically started.
* **Comprehensive container management**: Kit can manage both host processes and containers, providing a comprehensive
  view of the entire software system and quickly identifying any issues or bugs that arise.
* **Automatic rebuilding and restarting**: Kit can automatically rebuild and restart applications in response to changes
  in the source code or configuration files, allowing developers to test their code quickly and efficiently.
* **Parameterizable** Kit allows you to run tasks with different parameters.
* **Flexible integration and extensibility**: Kit is designed to be highly flexible and extensible, with support for a
  wide range of programming languages, frameworks, and tools. It can be easily integrated with existing systems and
  customized to meet specific needs.
* **Terminal output**: Tasks run concurrently and their status is muxed into a single terminal window, so you're not
  overwhelmed by pages of terminal output.
* **Log capture**: Logs are captured so you can look at them anytime.

Kit was written with extensive help from AI.

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

We do not support `go install`.

## Usage

Apps are described by a directed acyclic graph (DAG) of tasks.

Create a [`tasks.yaml`](tasks.yaml) file, e.g.:

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: my-proj
spec:
  tasks:
    - command: go build .
```

There are several kinds of task:

A **host task** runs on the host machine. It is defined by a `command`:

```yaml
- name: build
  command: go build .
```

Once a task completes successfully, any downstream tasks are started. If it is unsuccessful, Kit will exit

A task may finish, or run indefinitely. If it runs indefinitely, it is a **service** and  therefore must specify its ports:

```yaml
- name: service
  command: go run .
  ports: [ 8080 ]  
```

Unlike a plain task, if a service fails, it will be restarted. You can specify a `livenessProbe` to determine if the
service is running correctly:

```yaml
- name: service
  command: go run .
  ports: [ 8080 ]
  livenessProbe:
    httpGet:
      path: /health   
```

A **shell task** is just a host task that runs in a shell:

```yaml
- name: shell
  sh: |
    set -eux
    echo "Hello, world!"
```

A **container task** runs in a container. It is defined by an `image`:

```yaml
- name: mysql
  image: mysql
  ports: [ 3306:3306 ]
```

If the image is a path to a Dockerfile, it will be built:

```yaml
- name: kafka
  image: ./src/images/kafka
```

A **Kubernetes task** deploys manifests to a Kubernetes cluster:

```yaml
- name: deploy
  namespace: default
  manifests:
    - manifests/
    - service.yaml
  ports: [ 80:8080 ]
```

Start:

```bash
kit up
```

## Documentation

- [Usage](docs/USAGE.md) - how to use the various features of kit
- [Examples](docs/examples) - examples of how to use kit, e.g. with MySQL, or Kafka
- [Reference](docs/reference) - reference documentation for the various types in kit

