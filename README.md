# Joy

This is a tool to enable local development of containerized applications. It uses conventional Kubernetes Pod YAML, but
allows
you to run the process on your host (like `Procfile`), in Docker using just-in-time build (like Docker Compose).

Only high-level status is shown to your terminal:

```
▓ foo        [dead    ]  2022/12/21 17:48:19 listening on 8080
▓ bar        [unready ]  2022/12/21 17:48:20 port=9090
▓ baz        [excluded]  
▓ qux        [running ]  y[Thu Dec 22 01:48:21.957864 2022] [core:notice] [pid 1:tid 281
```

## Install

```bash
brew tap alexec/joy --custom-remote https://github.com/alexec/joy
brew install joy
```

## Usage

Describe you application in a [`joy.yaml`](joy.yaml) file using Kubernetes pod syntax, then start:

```bash
joy -h
````

Logs are stored in `./logs`.

### Container Process

The `image` field can be either:

1. An conventional image tag. E.g. `ubunutu`.
2. A path to a `Dockerfile`, e.g. `foo/Dockerfile`.

If it is a path to a `Dockerfile`, that file is built, and tagged with the container name.

```yaml
    # conventional image? run in Docker
    - name: baz
      image: httpd
    # path image? build and run in Docker
    - name: qux
      image: ./demo/qux/Dockerfile
```

Any container with the same name as the container name in the YAML is stopped and re-created whenever the process
starts.

### Host Process

The `image` field is empty. The value of `command` is used to start the process.

```yaml
    # no image? this is a host process
    - name: foo
      command: [ go, run, ./demo/foo ]
```

### Init Containers

Init containers are started before the main containers. They are allowed to run to completion before the main containers
are started. Useful if you want to do some set-up or build steps.

### Liveness Probe

If the process is not alive (i.e. "dead"), then it is killed and restarted. Just like Kubernetes.

### Quitting

Enter Ctrl+C to send a `SIGTERM` to the process. Each sub-process is then gets sent `SIGTERM`. If they do not exit
within 30s, then they get a `SIGKILL`. You may wish to reduce this number:

```yaml
spec:
  terminationGracePeriodSeconds: 3
```

You can kill the tool using `kill` for another terminal. If you `kill -9`, then the sub-process will keep
running and you must manually clean up.

## References

- [Containers from scratch](https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909)
