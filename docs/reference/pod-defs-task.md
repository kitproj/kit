# Task Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/Tasks/items
```

A task is a container or a command to run.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## items Type

`object` ([Task](pod-defs-task.md))

# items Properties

| Property                            | Type      | Required | Nullable       | Defined by                                                                                                                                                |
| :---------------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)                       | `string`  | Required | cannot be null | [Untitled schema](pod-defs-task-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/name")                       |
| [log](#log)                         | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-log.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/log")                         |
| [image](#image)                     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-image.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/image")                     |
| [imagePullPolicy](#imagepullpolicy) | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-imagepullpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/imagePullPolicy") |
| [livenessProbe](#livenessprobe)     | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/livenessProbe")                             |
| [readinessProbe](#readinessprobe)   | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/readinessProbe")                            |
| [command](#command)                 | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/command")                                 |
| [args](#args)                       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/args")                                    |
| [sh](#sh)                           | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-sh.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/sh")                           |
| [manifests](#manifests)             | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/manifests")                               |
| [workingDir](#workingdir)           | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-workingdir.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/workingDir")           |
| [user](#user)                       | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-user.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/user")                       |
| [env](#env)                         | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-envvars.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/env")                                     |
| [envfile](#envfile)                 | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-envfile.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/envfile")                                 |
| [ports](#ports)                     | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-ports.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/ports")                                     |
| [volumeMounts](#volumemounts)       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-volumemounts.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/volumeMounts")       |
| [tty](#tty)                         | `boolean` | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-tty.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/tty")                         |
| [watch](#watch)                     | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/watch")                                   |
| [mutex](#mutex)                     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-mutex.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/mutex")                     |
| [semaphore](#semaphore)             | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-semaphore.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/semaphore")             |
| [dependencies](#dependencies)       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/dependencies")                            |
| [targets](#targets)                 | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/targets")                                 |
| [restartPolicy](#restartpolicy)     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-restartpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/restartPolicy")     |

## name

The name of the task, must be unique

`name`

*   is required

*   Type: `string` ([name](pod-defs-task-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/name")

### name Type

`string` ([name](pod-defs-task-properties-name.md))

## log

Where to log the output of the task. E.g. if the task is verbose. Defaults to /dev/stdout. Maybe a file, or /dev/null.

`log`

*   is optional

*   Type: `string` ([log](pod-defs-task-properties-log.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-log.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/log")

### log Type

`string` ([log](pod-defs-task-properties-log.md))

## image

Either the container image to run, or a directory containing a Dockerfile. If omitted, the process runs on the host.

`image`

*   is optional

*   Type: `string` ([image](pod-defs-task-properties-image.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-image.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/image")

### image Type

`string` ([image](pod-defs-task-properties-image.md))

## imagePullPolicy

Pull policy, e.g. Always, Never, IfNotPresent

`imagePullPolicy`

*   is optional

*   Type: `string` ([imagePullPolicy](pod-defs-task-properties-imagepullpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-imagepullpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/imagePullPolicy")

### imagePullPolicy Type

`string` ([imagePullPolicy](pod-defs-task-properties-imagepullpolicy.md))

## livenessProbe

A probe to check if the task is alive, it will be restarted if not.

`livenessProbe`

*   is optional

*   Type: `object` ([Probe](pod-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/livenessProbe")

### livenessProbe Type

`object` ([Probe](pod-defs-probe.md))

## readinessProbe

A probe to check if the task is alive, it will be restarted if not.

`readinessProbe`

*   is optional

*   Type: `object` ([Probe](pod-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/readinessProbe")

### readinessProbe Type

`object` ([Probe](pod-defs-probe.md))

## command

The command to run in the container or on the host. If both the image and the command are omitted, this is a noop.

`command`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/command")

### command Type

`string[]`

## args

The arguments to pass to the command

`args`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/args")

### args Type

`string[]`

## sh

The shell script to run, instead of the command

`sh`

*   is optional

*   Type: `string` ([sh](pod-defs-task-properties-sh.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-sh.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/sh")

### sh Type

`string` ([sh](pod-defs-task-properties-sh.md))

## manifests

A directories or files of Kubernetes manifests to apply. Once running the task will wait for the resources to be ready.

`manifests`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/manifests")

### manifests Type

`string[]`

## workingDir

The working directory in the container or on the host

`workingDir`

*   is optional

*   Type: `string` ([workingDir](pod-defs-task-properties-workingdir.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-workingdir.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/workingDir")

### workingDir Type

`string` ([workingDir](pod-defs-task-properties-workingdir.md))

## user

The user to run the task as.

`user`

*   is optional

*   Type: `string` ([user](pod-defs-task-properties-user.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-user.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/user")

### user Type

`string` ([user](pod-defs-task-properties-user.md))

## env

A list of environment variables.

`env`

*   is optional

*   Type: `object[]` ([EnvVar](pod-defs-envvar.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvars.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/env")

### env Type

`object[]` ([EnvVar](pod-defs-envvar.md))

## envfile

Environment file (e.g. .env) to use

`envfile`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envfile.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/envfile")

### envfile Type

`string[]`

## ports

A list of ports to expose.

`ports`

*   is optional

*   Type: `object[]` ([Port](pod-defs-port.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-ports.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/ports")

### ports Type

`object[]` ([Port](pod-defs-port.md))

## volumeMounts

Volumes to mount in the container

`volumeMounts`

*   is optional

*   Type: `object[]` ([VolumeMount](pod-defs-volumemount.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-volumemounts.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/volumeMounts")

### volumeMounts Type

`object[]` ([VolumeMount](pod-defs-volumemount.md))

## tty

Use a pseudo-TTY

`tty`

*   is optional

*   Type: `boolean` ([tty](pod-defs-task-properties-tty.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-tty.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/tty")

### tty Type

`boolean` ([tty](pod-defs-task-properties-tty.md))

## watch

A list of files to watch for changes, and restart the task if they change

`watch`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/watch")

### watch Type

`string[]`

## mutex

A mutex to prevent multiple tasks with the same mutex from running at the same time

`mutex`

*   is optional

*   Type: `string` ([mutex](pod-defs-task-properties-mutex.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-mutex.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/mutex")

### mutex Type

`string` ([mutex](pod-defs-task-properties-mutex.md))

## semaphore

A semaphore to limit the number of tasks with the same semaphore that can run at the same time

`semaphore`

*   is optional

*   Type: `string` ([semaphore](pod-defs-task-properties-semaphore.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-semaphore.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/semaphore")

### semaphore Type

`string` ([semaphore](pod-defs-task-properties-semaphore.md))

## dependencies

A list of tasks to run before this task

`dependencies`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/dependencies")

### dependencies Type

`string[]`

## targets

A list of files this task will create. If these exist, and they're newer than the watched files, the task is skipped.

`targets`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/targets")

### targets Type

`string[]`

## restartPolicy

The restart policy, e.g. Always, Never, OnFailure. Defaults depends on the type of task.

`restartPolicy`

*   is optional

*   Type: `string` ([restartPolicy](pod-defs-task-properties-restartpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-restartpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/restartPolicy")

### restartPolicy Type

`string` ([restartPolicy](pod-defs-task-properties-restartpolicy.md))
