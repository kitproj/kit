# Examples

## Container

The `image` field can be either:

1. An conventional image tag. E.g. `ubunutu`.
2. A path to a a directory containing contain a `Dockerfile`, e.g. `.foo`.

If it is a path to a directory containing `Dockerfile`, that file is built, and tagged.

```yaml
    # conventional image? run in Docker
    - name: baz
      image: httpd
    # path image? build and run in Docker
    - name: qux
      image: demo/qux
```

Any container with the same name as the container name in the YAML is stopped and re-created whenever the process
starts.

## Host Process

If `image` field is omitted, the value of `command` is used to start the process on the host:

```yaml
    # no image? this is a host process
    - name: foo
      command: go run ./demo/foo 
```
## Noop

If `image` field is omitted and `command` is omitted, the task does nothing. This is used if you want to start several tasks, and conventionally you'd name the task `up`.

```yaml
    # no image or command? this is a noop
    - name: foo
```

## Parameters

You can specify environment variables for the task:

```yaml
    - name: foo
      command: echo $FOO
      env:
        FOO: bar
```

Environment variables specified in your shell environment are automatically passed to the task.

```bash
env FOO=qux kit up
```

Would print `qux` instead of `bar`.

### Auto Rebuild and Restart

You can specify a set of files to watch for changes that result in a re-build:

```yaml
  - watch: demo/bar
    name: bar
```        

## Liveness Probe

If the process is not alive (i.e. "dead"), then it is killed and restarted. Just like Kubernetes.

## Quitting

Enter Ctrl+C to send a `SIGTERM` to the process. Each sub-process is then gets sent `SIGTERM`. If they do not exit
_within 3s, then they get a `SIGKILL`. 

You can kill the tool using `kill` for another terminal. If you `kill -9`, then the sub-process will keep
running and you must manually clean up.

## Killing One Task

* To kill a host process: `kill $(lsof -ti:$host_port)`
* To kill a container : `docker kill $name`.

## Prebuild Patterns for Cloud Development

In most cases you will probably only have 1 top node in your command dependency graph (often named `up`). When developing in cloud workspaces (such as Codespaces, GitPod, etc.) it is common for teams to make use of "prebuilds" where longer running start-up tasks like dependency fetching are done in advance on every new commit so that when users start up a workspace these tasks can be pre-cached. In these cases it is recommended to have a task in your task list that can be run on prebuilds even if that task is not in your primary dependency graph. For example, if you have a java service that you need to run, it might make sense to have a separate `dependency-pull` task that is run as part of the prebuild `kit dependency-pull` separate from your primary `kit up` task

```yaml
  tasks:
    - name: dependency-pull
      command: "mvn clean install"
    - name: up
      command: "mvn spring-boot:run"
```