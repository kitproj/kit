# Kit

[![CodeQL](https://github.com/kitproj/kit/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/codeql-analysis.yml)
[![Go](https://github.com/kitproj/kit/actions/workflows/go.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/go.yml)
[![goreleaser](https://github.com/kitproj/kit/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/kitproj/kit/actions/workflows/goreleaser.yml)

## Why

Make the dev loop crazy fast.

## What

Kit is a software development tool designed to turbo-charge the software development process. 

Kit combines both task execution (like Foreman), container management (like Docker Compose), Kubernetes-like features (like Tilt, Skaffold), and a focus on local development (like Garden) in a single, easy-to-use binary.

It is designed to work seamlessly with both local-dev and cloud-dev environments, such as Codespaces and Gitpod.

[![Watch the video](https://img.youtube.com/vi/IafQwT1rYOU/hqdefault.jpg)](https://youtu.be/IafQwT1rYOU)


Key features of Kit include:

* **Local testing**: Kit is designed for local testing, allowing developers to test their code on their local machines before pushing it to a test environment or production. This speeds up the testing process and helps developers catch and fix bugs more quickly.
* **First-class container support**: Kit downloads and runs containers, allowing developers to test their code in a containerized environment without having to set up a separate container management system.
* **First-class Kubernetes support**: Kit can deploy manifests to a cluster and automatically port-forward to the cluster.
* **Advanced DAG architecture**: Kit's directed acyclic graph (DAG) structure allows for optimized parallel processing, reducing the time required for testing and speeding up the development process.
* **Probes**: You can specify liveness probes for your tasks to see if they're working, automatically restarting them when they go wrong. You can also specify readiness probes for your tasks to see if they're ready.
* **Dependency management**: You can specify dependencies between tasks, so when upstream tasks become successful or ready, downstream tasks are automatically started.
* **Comprehensive container management**: Kit can manage both host processes and containers, providing a comprehensive view of the entire software system and quickly identifying any issues or bugs that arise.
* **Automatic rebuilding and restarting**: Kit can automatically rebuild and restart applications in response to changes in the source code or configuration files, allowing developers to test their code quickly and efficiently.
* **Parameterizable** Kit allows you to run tasks with different parameters.
* **Flexible integration and extensibility**: Kit is designed to be highly flexible and extensible, with support for a wide range of programming languages, frameworks, and tools. It can be easily integrated with existing systems and customized to meet specific needs.
* **Terminal output**: Tasks run concurrently and their status is muxed into a single terminal window, so you're not overwhelmed by pages of terminal output.
* **Log capture**: Logs are captured so you can look at them anytime.

Kit was written with extensive help from AI.

## Install

Like `jq`, `kit` is a small standalone binary. You can download it from the [releases page](https://github.com/kitproj/kit/releases/latest). 

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

Apps are described by a DAG, for example:


```mermaid
---
title: Example of an app
---
flowchart LR
    api(name: api\ncommand: java -jar target/api.jar\nworkingDir: ./api\nports: 8080):::host --> build-api(name: build-api\ncommand: mvn package\nworkingDir: ./api\nwatch: ./api):::host
    api --> mysql(name: mysql\n image: mysql:latest):::container
    processor(name: processor\ncommand: ./processor):::host --> build-processor(name: build-processor\ncommand: go build ./processor\nwatch: ./processor):::host
    processor --> kafka(name: kafka\nimage: ./src/images/kafka):::container
    processor --> object-storage(name: object-storage\nimage: minio:latest):::container
    processor --> api
    ui(name: ui\ncommand: yarn start\nworkingDir: ./ui\nports: 4000):::host --> build-ui(name: build-ui\ncommand: yarn install\nworkingDir: ./ui):::host
    ui --> api
```

Create a [`tasks.yaml`](tasks.yaml) file, e.g.:

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: my-proj
spec:
  tasks:
    - command: go build -v .
      name: build-bar
      watch: demo/bar/main.go
      workingDir: demo/bar
    - command: ./demo/bar/bar
      dependencies: build-bar
      env:
        - PORT=9090
      name: bar
      ports: "9090"
    - dependencies: bar
      name: up
```

Start:

```bash
kit up
```

You'll see something like this:

![screenshot](screenshot.png)

Logs are stored in `./logs`.

## Documentation

- [Usage](docs/USAGE.md) - how to use the various features of kit
- [Examples](docs/examples) - examples of how to use kit, e.g. with MySQL, or Kafka
- [Reference](docs/reference) - reference documentation for the various types in kit

