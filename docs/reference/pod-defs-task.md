# Task Schema

```txt
https://github.com/alexec/kit/internal/types/pod#/$defs/Tasks/items
```

A task is a container or a command to run.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## items Type

`object` ([Task](pod-defs-task.md))

# items Properties

| Property                            | Type      | Required | Nullable       | Defined by                                                                                                                                               |
| :---------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)                       | `string`  | Required | cannot be null | [Untitled schema](pod-defs-task-properties-name.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/name")                       |
| [image](#image)                     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-image.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/image")                     |
| [imagePullPolicy](#imagepullpolicy) | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-imagepullpolicy.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/imagePullPolicy") |
| [livenessProbe](#livenessprobe)     | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-probe.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/livenessProbe")                             |
| [readinessProbe](#readinessprobe)   | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-probe.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/readinessProbe")                            |
| [command](#command)                 | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/command")                                 |
| [args](#args)                       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/args")                                    |
| [workingDir](#workingdir)           | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-workingdir.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/workingDir")           |
| [user](#user)                       | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-user.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/user")                       |
| [env](#env)                         | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-envvars.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/env")                                     |
| [ports](#ports)                     | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-ports.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/ports")                                     |
| [volumeMounts](#volumemounts)       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-volumemounts.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/volumeMounts")       |
| [tty](#tty)                         | `boolean` | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-tty.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/tty")                         |
| [watch](#watch)                     | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/watch")                                   |
| [mutex](#mutex)                     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-mutex.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/mutex")                     |
| [dependencies](#dependencies)       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/dependencies")                            |
| [restartPolicy](#restartpolicy)     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-restartpolicy.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/restartPolicy")     |

## name

The name of the task, must be unique

`name`

*   is required

*   Type: `string` ([name](pod-defs-task-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-name.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/name")

### name Type

`string` ([name](pod-defs-task-properties-name.md))

## image

Either the container image to run, or a directory containing a Dockerfile. If omitted, the process runs on the host.

`image`

*   is optional

*   Type: `string` ([image](pod-defs-task-properties-image.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-image.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/image")

### image Type

`string` ([image](pod-defs-task-properties-image.md))

## imagePullPolicy

Pull policy, e.g. Always, Never, IfNotPresent

`imagePullPolicy`

*   is optional

*   Type: `string` ([imagePullPolicy](pod-defs-task-properties-imagepullpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-imagepullpolicy.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/imagePullPolicy")

### imagePullPolicy Type

`string` ([imagePullPolicy](pod-defs-task-properties-imagepullpolicy.md))

## livenessProbe

A probe to check if the task is alive, it will be restarted if not.

`livenessProbe`

*   is optional

*   Type: `object` ([Probe](pod-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/livenessProbe")

### livenessProbe Type

`object` ([Probe](pod-defs-probe.md))

## readinessProbe

A probe to check if the task is alive, it will be restarted if not.

`readinessProbe`

*   is optional

*   Type: `object` ([Probe](pod-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/readinessProbe")

### readinessProbe Type

`object` ([Probe](pod-defs-probe.md))

## command

The command to run in the container or on the host. If both the image and the command are omitted, this is a noop.

`command`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/command")

### command Type

`string[]`

## args

The arguments to pass to the command

`args`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/args")

### args Type

`string[]`

## workingDir

The working directory in the container or on the host

`workingDir`

*   is optional

*   Type: `string` ([workingDir](pod-defs-task-properties-workingdir.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-workingdir.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/workingDir")

### workingDir Type

`string` ([workingDir](pod-defs-task-properties-workingdir.md))

## user

The user to run the task as.

`user`

*   is optional

*   Type: `string` ([user](pod-defs-task-properties-user.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-user.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/user")

### user Type

`string` ([user](pod-defs-task-properties-user.md))

## env

A list of environment variables.

`env`

*   is optional

*   Type: `object[]` ([EnvVar](pod-defs-envvar.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvars.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/env")

### env Type

`object[]` ([EnvVar](pod-defs-envvar.md))

## ports

A list of ports to expose.

`ports`

*   is optional

*   Type: `object[]` ([Port](pod-defs-port.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-ports.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/ports")

### ports Type

`object[]` ([Port](pod-defs-port.md))

## volumeMounts

Volumes to mount in the container

`volumeMounts`

*   is optional

*   Type: `object[]` ([VolumeMount](pod-defs-volumemount.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-volumemounts.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/volumeMounts")

### volumeMounts Type

`object[]` ([VolumeMount](pod-defs-volumemount.md))

## tty

Use a pseudo-TTY

`tty`

*   is optional

*   Type: `boolean` ([tty](pod-defs-task-properties-tty.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-tty.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/tty")

### tty Type

`boolean` ([tty](pod-defs-task-properties-tty.md))

## watch

A list of files to watch for changes, and restart the task if they change

`watch`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/watch")

### watch Type

`string[]`

## mutex

A mutex to prevent multiple tasks with the same mutex from running at the same time

`mutex`

*   is optional

*   Type: `string` ([mutex](pod-defs-task-properties-mutex.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-mutex.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/mutex")

### mutex Type

`string` ([mutex](pod-defs-task-properties-mutex.md))

## dependencies

A list of tasks to run before this task

`dependencies`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/dependencies")

### dependencies Type

`string[]`

## restartPolicy

The restart policy, e.g. Always, Never, OnFailure

`restartPolicy`

*   is optional

*   Type: `string` ([restartPolicy](pod-defs-task-properties-restartpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-restartpolicy.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Task/properties/restartPolicy")

### restartPolicy Type

`string` ([restartPolicy](pod-defs-task-properties-restartpolicy.md))
