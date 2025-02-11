# Untitled schema Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [workflow.schema.json](../../out/workflow.schema.json "open original schema") |

## Untitled schema Type

unknown

# Untitled schema Definitions

## Definitions group Duration

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Duration"}
```

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                |
| :-------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
| [Duration](#duration) | `object` | Required | cannot be null | [Untitled schema](workflow-defs-duration.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Duration/properties/Duration") |

### Duration



`Duration`

*   is required

*   Type: `object` ([Duration](workflow-defs-duration.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-duration.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Duration/properties/Duration")

#### Duration Type

`object` ([Duration](workflow-defs-duration.md))

## Definitions group EnvVars

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/EnvVars"}
```

| Property | Type     | Required | Nullable       | Defined by                                                                                                                                                  |
| :------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `.*`     | `string` | Optional | cannot be null | [Untitled schema](workflow-defs-envvars-patternproperties-.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/EnvVars/patternProperties/.*") |

### Pattern: `.*`



`.*`

*   is optional

*   Type: `string`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-envvars-patternproperties-.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/EnvVars/patternProperties/.*")

#### .\* Type

`string`

## Definitions group Envfile

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Envfile"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group HTTPGetAction

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction"}
```

| Property          | Type      | Required | Nullable       | Defined by                                                                                                                                                          |
| :---------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [scheme](#scheme) | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/scheme") |
| [port](#port)     | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/port")     |
| [path](#path)     | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/path")     |

### scheme

Scheme to use for connecting to the host. Defaults to HTTP.

`scheme`

*   is optional

*   Type: `string` ([scheme](workflow-defs-httpgetaction-properties-scheme.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/scheme")

#### scheme Type

`string` ([scheme](workflow-defs-httpgetaction-properties-scheme.md))

### port

Number of the port

`port`

*   is optional

*   Type: `integer` ([port](workflow-defs-httpgetaction-properties-port.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/port")

#### port Type

`integer` ([port](workflow-defs-httpgetaction-properties-port.md))

### path

Path to access on the HTTP server.

`path`

*   is optional

*   Type: `string` ([path](workflow-defs-httpgetaction-properties-path.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/path")

#### path Type

`string` ([path](workflow-defs-httpgetaction-properties-path.md))

## Definitions group HostPath

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/HostPath"}
```

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                            |
| :-------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
| [path](#path-1) | `string` | Required | cannot be null | [Untitled schema](workflow-defs-hostpath-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HostPath/properties/path") |

### path

Path of the directory on the host.

`path`

*   is required

*   Type: `string` ([path](workflow-defs-hostpath-properties-path.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-hostpath-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HostPath/properties/path")

#### path Type

`string` ([path](workflow-defs-hostpath-properties-path.md))

## Definitions group Port

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port"}
```

| Property                        | Type      | Required | Nullable       | Defined by                                                                                                                                                      |
| :------------------------------ | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [containerPort](#containerport) | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/containerPort") |
| [hostPort](#hostport)           | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/hostPort")           |

### containerPort

The container port to expose

`containerPort`

*   is optional

*   Type: `integer` ([containerPort](workflow-defs-port-properties-containerport.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/containerPort")

#### containerPort Type

`integer` ([containerPort](workflow-defs-port-properties-containerport.md))

### hostPort

The host port to route to the container port

`hostPort`

*   is optional

*   Type: `integer` ([hostPort](workflow-defs-port-properties-hostport.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/hostPort")

#### hostPort Type

`integer` ([hostPort](workflow-defs-port-properties-hostport.md))

## Definitions group Ports

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Ports"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group Probe

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe"}
```

| Property                                    | Type      | Required | Nullable       | Defined by                                                                                                                                                                    |
| :------------------------------------------ | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [tcpSocket](#tcpsocket)                     | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-tcpsocketaction.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/tcpSocket")                                |
| [httpGet](#httpget)                         | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-httpgetaction.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/httpGet")                                    |
| [initialDelaySeconds](#initialdelayseconds) | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-probe-properties-initialdelayseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/initialDelaySeconds") |
| [periodSeconds](#periodseconds)             | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-probe-properties-periodseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/periodSeconds")             |
| [successThreshold](#successthreshold)       | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-probe-properties-successthreshold.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/successThreshold")       |
| [failureThreshold](#failurethreshold)       | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-probe-properties-failurethreshold.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/failureThreshold")       |

### tcpSocket

TCPSocketAction describes an action based on opening a socket

`tcpSocket`

*   is optional

*   Type: `object` ([TCPSocketAction](workflow-defs-tcpsocketaction.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-tcpsocketaction.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/tcpSocket")

#### tcpSocket Type

`object` ([TCPSocketAction](workflow-defs-tcpsocketaction.md))

### httpGet

HTTPGetAction describes an action based on HTTP Locks requests.

`httpGet`

*   is optional

*   Type: `object` ([HTTPGetAction](workflow-defs-httpgetaction.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-httpgetaction.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/httpGet")

#### httpGet Type

`object` ([HTTPGetAction](workflow-defs-httpgetaction.md))

### initialDelaySeconds

Number of seconds after the process has started before the probe is initiated.

`initialDelaySeconds`

*   is optional

*   Type: `integer` ([initialDelaySeconds](workflow-defs-probe-properties-initialdelayseconds.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-probe-properties-initialdelayseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/initialDelaySeconds")

#### initialDelaySeconds Type

`integer` ([initialDelaySeconds](workflow-defs-probe-properties-initialdelayseconds.md))

### periodSeconds

How often (in seconds) to perform the probe.

`periodSeconds`

*   is optional

*   Type: `integer` ([periodSeconds](workflow-defs-probe-properties-periodseconds.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-probe-properties-periodseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/periodSeconds")

#### periodSeconds Type

`integer` ([periodSeconds](workflow-defs-probe-properties-periodseconds.md))

### successThreshold

Minimum consecutive successes for the probe to be considered successful after having failed.

`successThreshold`

*   is optional

*   Type: `integer` ([successThreshold](workflow-defs-probe-properties-successthreshold.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-probe-properties-successthreshold.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/successThreshold")

#### successThreshold Type

`integer` ([successThreshold](workflow-defs-probe-properties-successthreshold.md))

### failureThreshold

Minimum consecutive failures for the probe to be considered failed after having succeeded.

`failureThreshold`

*   is optional

*   Type: `integer` ([failureThreshold](workflow-defs-probe-properties-failurethreshold.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-probe-properties-failurethreshold.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/failureThreshold")

#### failureThreshold Type

`integer` ([failureThreshold](workflow-defs-probe-properties-failurethreshold.md))

## Definitions group Strings

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Strings"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group TCPSocketAction

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/TCPSocketAction"}
```

| Property        | Type      | Required | Nullable       | Defined by                                                                                                                                                          |
| :-------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [port](#port-1) | `integer` | Required | cannot be null | [Untitled schema](workflow-defs-tcpsocketaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/TCPSocketAction/properties/port") |

### port

Port number of the port to probe.

`port`

*   is required

*   Type: `integer` ([port](workflow-defs-tcpsocketaction-properties-port.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-tcpsocketaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/TCPSocketAction/properties/port")

#### port Type

`integer` ([port](workflow-defs-tcpsocketaction-properties-port.md))

## Definitions group Task

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task"}
```

| Property                            | Type      | Required | Nullable       | Defined by                                                                                                                                                          |
| :---------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [type](#type)                       | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-type.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/type")                       |
| [log](#log)                         | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-log.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/log")                         |
| [image](#image)                     | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-image.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/image")                     |
| [imagePullPolicy](#imagepullpolicy) | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-imagepullpolicy.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/imagePullPolicy") |
| [livenessProbe](#livenessprobe)     | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-probe.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/livenessProbe")                             |
| [readinessProbe](#readinessprobe)   | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-probe.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/readinessProbe")                            |
| [command](#command)                 | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/command")                                 |
| [args](#args)                       | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/args")                                    |
| [sh](#sh)                           | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-sh.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/sh")                           |
| [manifests](#manifests)             | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/manifests")                               |
| [namespace](#namespace)             | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-namespace.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/namespace")             |
| [workingDir](#workingdir)           | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-workingdir.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/workingDir")           |
| [user](#user)                       | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-user.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/user")                       |
| [env](#env)                         | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-envvars.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/env")                                     |
| [envfile](#envfile)                 | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-envfile.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/envfile")                                 |
| [ports](#ports)                     | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-ports.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/ports")                                     |
| [volumeMounts](#volumemounts)       | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-volumemounts.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/volumeMounts")       |
| [tty](#tty)                         | `boolean` | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-tty.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/tty")                         |
| [watch](#watch)                     | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/watch")                                   |
| [mutex](#mutex)                     | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-mutex.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/mutex")                     |
| [semaphore](#semaphore)             | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-semaphore.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/semaphore")             |
| [dependencies](#dependencies)       | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/dependencies")                            |
| [targets](#targets)                 | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/targets")                                 |
| [restartPolicy](#restartpolicy)     | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-task-properties-restartpolicy.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/restartPolicy")     |
| [stalledTimeout](#stalledtimeout)   | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-duration.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/stalledTimeout")                         |

### type

Type is the type of the task: "service" or "job". If omitted, if there are ports, it's a service, otherwise it's a job.
This is only needed when you have service that does not listen on ports.
Services are running in the background.

`type`

*   is optional

*   Type: `string` ([type](workflow-defs-task-properties-type.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-type.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/type")

#### type Type

`string` ([type](workflow-defs-task-properties-type.md))

### log

Where to log the output of the task. E.g. if the task is verbose. Defaults to /dev/stdout. Maybe a file, or /dev/null.

`log`

*   is optional

*   Type: `string` ([log](workflow-defs-task-properties-log.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-log.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/log")

#### log Type

`string` ([log](workflow-defs-task-properties-log.md))

### image

Either the container image to run, or a directory containing a Dockerfile. If omitted, the process runs on the host.

`image`

*   is optional

*   Type: `string` ([image](workflow-defs-task-properties-image.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-image.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/image")

#### image Type

`string` ([image](workflow-defs-task-properties-image.md))

### imagePullPolicy

Pull policy, e.g. Always, Never, IfNotPresent

`imagePullPolicy`

*   is optional

*   Type: `string` ([imagePullPolicy](workflow-defs-task-properties-imagepullpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-imagepullpolicy.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/imagePullPolicy")

#### imagePullPolicy Type

`string` ([imagePullPolicy](workflow-defs-task-properties-imagepullpolicy.md))

### livenessProbe

A probe to check if the task is alive, it will be restarted if not.

`livenessProbe`

*   is optional

*   Type: `object` ([Probe](workflow-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-probe.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/livenessProbe")

#### livenessProbe Type

`object` ([Probe](workflow-defs-probe.md))

### readinessProbe

A probe to check if the task is alive, it will be restarted if not.

`readinessProbe`

*   is optional

*   Type: `object` ([Probe](workflow-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-probe.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/readinessProbe")

#### readinessProbe Type

`object` ([Probe](workflow-defs-probe.md))

### command

The command to run in the container or on the host. If both the image and the command are omitted, this is a noop.

`command`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/command")

#### command Type

`string[]`

### args

The arguments to pass to the command

`args`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/args")

#### args Type

`string[]`

### sh

The shell script to run, instead of the command

`sh`

*   is optional

*   Type: `string` ([sh](workflow-defs-task-properties-sh.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-sh.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/sh")

#### sh Type

`string` ([sh](workflow-defs-task-properties-sh.md))

### manifests

A directories or files of Kubernetes manifests to apply. Once running the task will wait for the resources to be ready.

`manifests`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/manifests")

#### manifests Type

`string[]`

### namespace

The namespace to run the Kubernetes resource in. Defaults to the namespace of the current Kubernetes context.

`namespace`

*   is optional

*   Type: `string` ([namespace](workflow-defs-task-properties-namespace.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-namespace.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/namespace")

#### namespace Type

`string` ([namespace](workflow-defs-task-properties-namespace.md))

### workingDir

The working directory in the container or on the host

`workingDir`

*   is optional

*   Type: `string` ([workingDir](workflow-defs-task-properties-workingdir.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-workingdir.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/workingDir")

#### workingDir Type

`string` ([workingDir](workflow-defs-task-properties-workingdir.md))

### user

The user to run the task as.

`user`

*   is optional

*   Type: `string` ([user](workflow-defs-task-properties-user.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-user.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/user")

#### user Type

`string` ([user](workflow-defs-task-properties-user.md))

### env

A list of environment variables.

`env`

*   is optional

*   Type: `object` ([EnvVars](workflow-defs-envvars.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-envvars.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/env")

#### env Type

`object` ([EnvVars](workflow-defs-envvars.md))

### envfile

Environment file (e.g. .env) to use

`envfile`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-envfile.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/envfile")

#### envfile Type

`string[]`

### ports

A list of ports to expose.

`ports`

*   is optional

*   Type: `object[]` ([Port](workflow-defs-port.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-ports.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/ports")

#### ports Type

`object[]` ([Port](workflow-defs-port.md))

### volumeMounts

Volumes to mount in the container

`volumeMounts`

*   is optional

*   Type: `object[]` ([VolumeMount](workflow-defs-volumemount.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-volumemounts.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/volumeMounts")

#### volumeMounts Type

`object[]` ([VolumeMount](workflow-defs-volumemount.md))

### tty

Use a pseudo-TTY

`tty`

*   is optional

*   Type: `boolean` ([tty](workflow-defs-task-properties-tty.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-tty.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/tty")

#### tty Type

`boolean` ([tty](workflow-defs-task-properties-tty.md))

### watch

A list of files to watch for changes, and restart the task if they change

`watch`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/watch")

#### watch Type

`string[]`

### mutex

A mutex to prevent multiple tasks with the same mutex from running at the same time

`mutex`

*   is optional

*   Type: `string` ([mutex](workflow-defs-task-properties-mutex.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-mutex.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/mutex")

#### mutex Type

`string` ([mutex](workflow-defs-task-properties-mutex.md))

### semaphore

A semaphore to limit the number of tasks with the same semaphore that can run at the same time

`semaphore`

*   is optional

*   Type: `string` ([semaphore](workflow-defs-task-properties-semaphore.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-semaphore.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/semaphore")

#### semaphore Type

`string` ([semaphore](workflow-defs-task-properties-semaphore.md))

### dependencies

A list of tasks to run before this task

`dependencies`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/dependencies")

#### dependencies Type

`string[]`

### targets

A list of files this task will create. If these exist, and they're newer than the watched files, the task is skipped.

`targets`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-strings.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/targets")

#### targets Type

`string[]`

### restartPolicy

The restart policy, e.g. Always, Never, OnFailure. Defaults depends on the type of task.

`restartPolicy`

*   is optional

*   Type: `string` ([restartPolicy](workflow-defs-task-properties-restartpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task-properties-restartpolicy.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/restartPolicy")

#### restartPolicy Type

`string` ([restartPolicy](workflow-defs-task-properties-restartpolicy.md))

### stalledTimeout

The timeout for the task to be considered stalled. If omitted, the task will be considered stalled after 30 seconds of no activity.

`stalledTimeout`

*   is optional

*   Type: `object` ([Duration](workflow-defs-duration.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-duration.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/stalledTimeout")

#### stalledTimeout Type

`object` ([Duration](workflow-defs-duration.md))

## Definitions group Tasks

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Tasks"}
```

| Property | Type     | Required | Nullable       | Defined by                                                                                                                          |
| :------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------- |
| `.*`     | `object` | Optional | cannot be null | [Untitled schema](workflow-defs-task.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Tasks/patternProperties/.*") |

### Pattern: `.*`

A task is a container or a command to run.

`.*`

*   is optional

*   Type: `object` ([Task](workflow-defs-task.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Tasks/patternProperties/.*")

#### .\* Type

`object` ([Task](workflow-defs-task.md))

## Definitions group Volume

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume"}
```

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                        |
| :-------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| [name](#name)         | `string` | Required | cannot be null | [Untitled schema](workflow-defs-volume-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/name") |
| [hostPath](#hostpath) | `object` | Required | cannot be null | [Untitled schema](workflow-defs-hostpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/hostPath")           |

### name

Volume's name.

`name`

*   is required

*   Type: `string` ([name](workflow-defs-volume-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-volume-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/name")

#### name Type

`string` ([name](workflow-defs-volume-properties-name.md))

### hostPath

HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container.

`hostPath`

*   is required

*   Type: `object` ([HostPath](workflow-defs-hostpath.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-hostpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/hostPath")

#### hostPath Type

`object` ([HostPath](workflow-defs-hostpath.md))

## Definitions group VolumeMount

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount"}
```

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                            |
| :---------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name-1)         | `string` | Required | cannot be null | [Untitled schema](workflow-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/name")           |
| [mountPath](#mountpath) | `string` | Required | cannot be null | [Untitled schema](workflow-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/mountPath") |

### name

This must match the name of a volume.

`name`

*   is required

*   Type: `string` ([name](workflow-defs-volumemount-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/name")

#### name Type

`string` ([name](workflow-defs-volumemount-properties-name.md))

### mountPath

Path within the container at which the volume should be mounted.

`mountPath`

*   is required

*   Type: `string` ([mountPath](workflow-defs-volumemount-properties-mountpath.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/mountPath")

#### mountPath Type

`string` ([mountPath](workflow-defs-volumemount-properties-mountpath.md))

## Definitions group Workflow

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow"}
```

| Property                                                        | Type      | Required | Nullable       | Defined by                                                                                                                                                                                              |
| :-------------------------------------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [terminationGracePeriodSeconds](#terminationgraceperiodseconds) | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-workflow-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/terminationGracePeriodSeconds") |
| [tasks](#tasks)                                                 | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-tasks.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/tasks")                                                                     |
| [volumes](#volumes)                                             | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-workflow-properties-volumes.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/volumes")                                             |
| [semaphores](#semaphores)                                       | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-workflow-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/semaphores")                                       |
| [env](#env-1)                                                   | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-envvars.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/env")                                                                     |
| [envfile](#envfile-1)                                           | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-envfile.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/envfile")                                                                 |

### terminationGracePeriodSeconds



`terminationGracePeriodSeconds`

*   is optional

*   Type: `integer` ([terminationGracePeriodSeconds](workflow-defs-workflow-properties-terminationgraceperiodseconds.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-workflow-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/terminationGracePeriodSeconds")

#### terminationGracePeriodSeconds Type

`integer` ([terminationGracePeriodSeconds](workflow-defs-workflow-properties-terminationgraceperiodseconds.md))

### tasks



`tasks`

*   is optional

*   Type: `object` ([Tasks](workflow-defs-tasks.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-tasks.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/tasks")

#### tasks Type

`object` ([Tasks](workflow-defs-tasks.md))

### volumes



`volumes`

*   is optional

*   Type: `object[]` ([Volume](workflow-defs-volume.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-workflow-properties-volumes.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/volumes")

#### volumes Type

`object[]` ([Volume](workflow-defs-volume.md))

### semaphores



`semaphores`

*   is optional

*   Type: `object` ([semaphores](workflow-defs-workflow-properties-semaphores.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-workflow-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/semaphores")

#### semaphores Type

`object` ([semaphores](workflow-defs-workflow-properties-semaphores.md))

### env

A list of environment variables.

`env`

*   is optional

*   Type: `object` ([EnvVars](workflow-defs-envvars.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-envvars.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/env")

#### env Type

`object` ([EnvVars](workflow-defs-envvars.md))

### envfile



`envfile`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-envfile.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/envfile")

#### envfile Type

`string[]`
