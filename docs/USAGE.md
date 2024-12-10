# How to

## How to run a container from an image.

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

## How to run a process on the host

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

## How to parameterize your tasks

You can specify environment variables for the task:

```yaml
    - name: foo
      command: echo $FOO
      env:
        - FOO: bar
```

Environment variables specified in your shell environment are automatically passed to the task.

```bash
env FOO=qux kit up
```

Would print `qux` instead of `bar`, but bear in mind that the shell environment variables takes precedence

You can set-up environment variables for all tasks:

```yaml
spec:
    env:
        - FOO: bar
    envfile: .env
```

Or per task:

```yaml
    - name: foo
      command: sh -c 'echo $FOO'
      env:
        - FOO: bar
```

## How to automatically re-build and restart a task

You can specify a set of files to watch for changes that result in a re-build:

```yaml
  - watch: demo/bar
    name: bar
```        

You can disable watches:

```bash
WATCH=0 kit up
```

OR

```bash
kit -W up
```

## How to automatically kill a broken task

If the process is not alive (i.e. "dead"), then it is killed and restarted. Just like Kubernetes.

```yaml
- name: foo
  image: httpd
  ports: 80:8080
  livenessProbe: http://:8080/health
```

## How to run a shell script

You might want to put a multi-line a shell script. You can do this using the YAML `|` character: 

```yaml
  - name: foo
    command:
      - sh
      - -c
      - |
        echo "hello"
        echo "world"
```

## How to wait for forked tasks

Kit's job is to manage your tasks for you. As a result, if you fork a task in a script, and the script exists, kit will terminate that forked task.

If you fork tasks, you can add a wait:

```yaml
  - name: foo
    command:
      - sh
      - -c
      - |
        sleep 60 &
        wait
```

## Quitting

Enter Ctrl+C to send a `SIGTERM` to the process. Each sub-process is then gets sent `SIGTERM`. If they do not exit
_within 3s, then they get a `SIGKILL`. 

You can kill the tool using `kill` for another terminal. If you `kill -9`, then the sub-process will keep
running and you must manually clean up.

## How to kill a single taks

* To kill a host process: `kill $(lsof -ti:$host_port)`
* To kill a container : `docker kill $name`.

## How to pre-build for Cloud Development environments

In most cases you will probably only have 1 top node in your command dependency graph (often named `up`). When developing in cloud workspaces (such as Codespaces, GitPod, etc.) it is common for teams to make use of "prebuilds" where longer running start-up tasks like dependency fetching are done in advance on every new commit so that when users start up a workspace these tasks can be pre-cached. In these cases it is recommended to have a task in your task list that can be run on prebuilds even if that task is not in your primary dependency graph. For example, if you have a java service that you need to run, it might make sense to have a separate `pre-up` task that is run as part of the prebuild `kit pre-up` separate from your primary `kit up` task

```yaml
  tasks:
    - name: pre-up
      command: "mvn dependency:go-offline"
    - name: up
      command: "mvn spring-boot:run"
      dependencies: pre-up
```

## How to just start downstream tasks?

Have a task named `downstream`, e.g.

```yaml
  tasks:
    - name: kafka
      image: ghcr.io/kitproj/kafka
    - name: sim
      image: ghcr.io/kitproj/sim
    - name: downstream
      dependencies: kafka sim
    - name: up
      command: "mvn spring-boot:run"
      dependencies: downstream
```

## How to prevent two tasks from running at the same time

### Mutexes

If you want to prevent two tasks from running at the same time, you can use a mutex:

```yaml
tasks:
- name: foo
  mutex: my-mutex
- name: bar
  mutex: my-mutex
```

### Semaphores

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
