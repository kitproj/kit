# Untitled schema Schema

```txt
https://github.com/kitproj/kit/internal/types/pod
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [pod.schema.json](../../out/pod.schema.json "open original schema") |

## Untitled schema Type

unknown

# Untitled schema Definitions

## Definitions group EnvVar

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar"}
```

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                |
| :-------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)   | `string` | Required | cannot be null | [Untitled schema](pod-defs-envvar-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/name")   |
| [value](#value) | `string` | Required | cannot be null | [Untitled schema](pod-defs-envvar-properties-value.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/value") |

### name



`name`

*   is required

*   Type: `string` ([name](pod-defs-envvar-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvar-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/name")

#### name Type

`string` ([name](pod-defs-envvar-properties-name.md))

### value



`value`

*   is required

*   Type: `string` ([value](pod-defs-envvar-properties-value.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvar-properties-value.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/value")

#### value Type

`string` ([value](pod-defs-envvar-properties-value.md))

## Definitions group EnvVars

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVars"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group HTTPGetAction

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction"}
```

| Property          | Type      | Required | Nullable       | Defined by                                                                                                                                                |
| :---------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [scheme](#scheme) | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/scheme") |
| [port](#port)     | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/port")     |
| [path](#path)     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/path")     |

### scheme

Scheme to use for connecting to the host. Defaults to HTTP.

`scheme`

*   is optional

*   Type: `string` ([scheme](pod-defs-httpgetaction-properties-scheme.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/scheme")

#### scheme Type

`string` ([scheme](pod-defs-httpgetaction-properties-scheme.md))

### port

Number of the port

`port`

*   is optional

*   Type: `integer` ([port](pod-defs-httpgetaction-properties-port.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/port")

#### port Type

`integer` ([port](pod-defs-httpgetaction-properties-port.md))

### path

Path to access on the HTTP server.

`path`

*   is optional

*   Type: `string` ([path](pod-defs-httpgetaction-properties-path.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/path")

#### path Type

`string` ([path](pod-defs-httpgetaction-properties-path.md))

## Definitions group HostPath

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/HostPath"}
```

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                  |
| :-------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------ |
| [path](#path-1) | `string` | Required | cannot be null | [Untitled schema](pod-defs-hostpath-properties-path.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HostPath/properties/path") |

### path

Path of the directory on the host.

`path`

*   is required

*   Type: `string` ([path](pod-defs-hostpath-properties-path.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-hostpath-properties-path.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HostPath/properties/path")

#### path Type

`string` ([path](pod-defs-hostpath-properties-path.md))

## Definitions group Metadata

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata"}
```

| Property                    | Type     | Required | Nullable       | Defined by                                                                                                                                                |
| :-------------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name-1)             | `string` | Required | cannot be null | [Untitled schema](pod-defs-metadata-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/name")               |
| [annotations](#annotations) | `object` | Optional | cannot be null | [Untitled schema](pod-defs-metadata-properties-annotations.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/annotations") |

### name

Name is the name of the resource.

`name`

*   is required

*   Type: `string` ([name](pod-defs-metadata-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-metadata-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/name")

#### name Type

`string` ([name](pod-defs-metadata-properties-name.md))

### annotations

Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata.

`annotations`

*   is optional

*   Type: `object` ([annotations](pod-defs-metadata-properties-annotations.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-metadata-properties-annotations.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/annotations")

#### annotations Type

`object` ([annotations](pod-defs-metadata-properties-annotations.md))

## Definitions group Pod

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod"}
```

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                    |
| :------------------------ | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------- |
| [spec](#spec)             | `object` | Required | cannot be null | [Untitled schema](pod-defs-podspec.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/spec")                         |
| [apiVersion](#apiversion) | `string` | Optional | cannot be null | [Untitled schema](pod-defs-pod-properties-apiversion.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/apiVersion") |
| [kind](#kind)             | `string` | Optional | cannot be null | [Untitled schema](pod-defs-pod-properties-kind.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/kind")             |
| [metadata](#metadata)     | `object` | Required | cannot be null | [Untitled schema](pod-defs-metadata.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/metadata")                    |

### spec

Task is a unit of work that should be run.

`spec`

*   is required

*   Type: `object` ([PodSpec](pod-defs-podspec.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-podspec.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/spec")

#### spec Type

`object` ([PodSpec](pod-defs-podspec.md))

### apiVersion

APIVersion must be `kit/v1`.

`apiVersion`

*   is optional

*   Type: `string` ([apiVersion](pod-defs-pod-properties-apiversion.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-pod-properties-apiversion.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/apiVersion")

#### apiVersion Type

`string` ([apiVersion](pod-defs-pod-properties-apiversion.md))

### kind

Kind must be `Tasks`.

`kind`

*   is optional

*   Type: `string` ([kind](pod-defs-pod-properties-kind.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-pod-properties-kind.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/kind")

#### kind Type

`string` ([kind](pod-defs-pod-properties-kind.md))

### metadata

Metadata is the metadata for the pod.

`metadata`

*   is required

*   Type: `object` ([Metadata](pod-defs-metadata.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-metadata.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/metadata")

#### metadata Type

`object` ([Metadata](pod-defs-metadata.md))

## Definitions group PodSpec

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec"}
```

| Property                                                        | Type      | Required | Nullable       | Defined by                                                                                                                                                                                  |
| :-------------------------------------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [terminationGracePeriodSeconds](#terminationgraceperiodseconds) | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-podspec-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/terminationGracePeriodSeconds") |
| [tasks](#tasks)                                                 | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-tasks.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/tasks")                                                                    |
| [volumes](#volumes)                                             | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-podspec-properties-volumes.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/volumes")                                             |
| [semaphores](#semaphores)                                       | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-podspec-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/semaphores")                                       |

### terminationGracePeriodSeconds

TerminationGracePeriodSeconds is the grace period for terminating the pod.

`terminationGracePeriodSeconds`

*   is optional

*   Type: `integer` ([terminationGracePeriodSeconds](pod-defs-podspec-properties-terminationgraceperiodseconds.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-podspec-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/terminationGracePeriodSeconds")

#### terminationGracePeriodSeconds Type

`integer` ([terminationGracePeriodSeconds](pod-defs-podspec-properties-terminationgraceperiodseconds.md))

### tasks

Tasks is a list of tasks that should be run.

`tasks`

*   is optional

*   Type: `object[]` ([Task](pod-defs-task.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-tasks.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/tasks")

#### tasks Type

`object[]` ([Task](pod-defs-task.md))

### volumes

Volumes is a list of volumes that can be mounted by containers belonging to the pod.

`volumes`

*   is optional

*   Type: `object[]` ([Volume](pod-defs-volume.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-podspec-properties-volumes.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/volumes")

#### volumes Type

`object[]` ([Volume](pod-defs-volume.md))

### semaphores

Semaphores is a list of semaphores that can be acquired by tasks.

`semaphores`

*   is optional

*   Type: `object` ([semaphores](pod-defs-podspec-properties-semaphores.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-podspec-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/PodSpec/properties/semaphores")

#### semaphores Type

`object` ([semaphores](pod-defs-podspec-properties-semaphores.md))

## Definitions group Port

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Port"}
```

| Property                        | Type      | Required | Nullable       | Defined by                                                                                                                                            |
| :------------------------------ | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
| [containerPort](#containerport) | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/containerPort") |
| [hostPort](#hostport)           | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/hostPort")           |

### containerPort

The container port to expose

`containerPort`

*   is optional

*   Type: `integer` ([containerPort](pod-defs-port-properties-containerport.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/containerPort")

#### containerPort Type

`integer` ([containerPort](pod-defs-port-properties-containerport.md))

### hostPort

The host port to route to the container port

`hostPort`

*   is optional

*   Type: `integer` ([hostPort](pod-defs-port-properties-hostport.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/hostPort")

#### hostPort Type

`integer` ([hostPort](pod-defs-port-properties-hostport.md))

## Definitions group Ports

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Ports"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group Probe

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe"}
```

| Property                                    | Type      | Required | Nullable       | Defined by                                                                                                                                                          |
| :------------------------------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [tcpSocket](#tcpsocket)                     | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-tcpsocketaction.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/tcpSocket")                                |
| [httpGet](#httpget)                         | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/httpGet")                                    |
| [initialDelaySeconds](#initialdelayseconds) | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-initialdelayseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/initialDelaySeconds") |
| [periodSeconds](#periodseconds)             | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-periodseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/periodSeconds")             |
| [successThreshold](#successthreshold)       | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-successthreshold.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/successThreshold")       |
| [failureThreshold](#failurethreshold)       | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-failurethreshold.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/failureThreshold")       |

### tcpSocket

TCPSocketAction describes an action based on opening a socket

`tcpSocket`

*   is optional

*   Type: `object` ([TCPSocketAction](pod-defs-tcpsocketaction.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-tcpsocketaction.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/tcpSocket")

#### tcpSocket Type

`object` ([TCPSocketAction](pod-defs-tcpsocketaction.md))

### httpGet

HTTPGetAction describes an action based on HTTP Locks requests.

`httpGet`

*   is optional

*   Type: `object` ([HTTPGetAction](pod-defs-httpgetaction.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/httpGet")

#### httpGet Type

`object` ([HTTPGetAction](pod-defs-httpgetaction.md))

### initialDelaySeconds

Number of seconds after the process has started before the probe is initiated.

`initialDelaySeconds`

*   is optional

*   Type: `integer` ([initialDelaySeconds](pod-defs-probe-properties-initialdelayseconds.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-initialdelayseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/initialDelaySeconds")

#### initialDelaySeconds Type

`integer` ([initialDelaySeconds](pod-defs-probe-properties-initialdelayseconds.md))

### periodSeconds

How often (in seconds) to perform the probe.

`periodSeconds`

*   is optional

*   Type: `integer` ([periodSeconds](pod-defs-probe-properties-periodseconds.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-periodseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/periodSeconds")

#### periodSeconds Type

`integer` ([periodSeconds](pod-defs-probe-properties-periodseconds.md))

### successThreshold

Minimum consecutive successes for the probe to be considered successful after having failed.

`successThreshold`

*   is optional

*   Type: `integer` ([successThreshold](pod-defs-probe-properties-successthreshold.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-successthreshold.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/successThreshold")

#### successThreshold Type

`integer` ([successThreshold](pod-defs-probe-properties-successthreshold.md))

### failureThreshold

Minimum consecutive failures for the probe to be considered failed after having succeeded.

`failureThreshold`

*   is optional

*   Type: `integer` ([failureThreshold](pod-defs-probe-properties-failurethreshold.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-failurethreshold.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/failureThreshold")

#### failureThreshold Type

`integer` ([failureThreshold](pod-defs-probe-properties-failurethreshold.md))

## Definitions group Strings

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Strings"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group TCPSocketAction

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/TCPSocketAction"}
```

| Property        | Type      | Required | Nullable       | Defined by                                                                                                                                                |
| :-------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [port](#port-1) | `integer` | Required | cannot be null | [Untitled schema](pod-defs-tcpsocketaction-properties-port.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/TCPSocketAction/properties/port") |

### port

Port number of the port to probe.

`port`

*   is required

*   Type: `integer` ([port](pod-defs-tcpsocketaction-properties-port.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-tcpsocketaction-properties-port.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/TCPSocketAction/properties/port")

#### port Type

`integer` ([port](pod-defs-tcpsocketaction-properties-port.md))

## Definitions group Task

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Task"}
```

| Property                            | Type      | Required | Nullable       | Defined by                                                                                                                                                |
| :---------------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name-2)                     | `string`  | Required | cannot be null | [Untitled schema](pod-defs-task-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/name")                       |
| [image](#image)                     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-image.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/image")                     |
| [imagePullPolicy](#imagepullpolicy) | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-imagepullpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/imagePullPolicy") |
| [livenessProbe](#livenessprobe)     | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/livenessProbe")                             |
| [readinessProbe](#readinessprobe)   | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/readinessProbe")                            |
| [command](#command)                 | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/command")                                 |
| [args](#args)                       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/args")                                    |
| [workingDir](#workingdir)           | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-workingdir.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/workingDir")           |
| [user](#user)                       | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-user.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/user")                       |
| [env](#env)                         | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-envvars.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/env")                                     |
| [ports](#ports)                     | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-ports.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/ports")                                     |
| [volumeMounts](#volumemounts)       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-volumemounts.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/volumeMounts")       |
| [tty](#tty)                         | `boolean` | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-tty.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/tty")                         |
| [watch](#watch)                     | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/watch")                                   |
| [mutex](#mutex)                     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-mutex.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/mutex")                     |
| [semaphore](#semaphore)             | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-semaphore.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/semaphore")             |
| [dependencies](#dependencies)       | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/dependencies")                            |
| [restartPolicy](#restartpolicy)     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-task-properties-restartpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/restartPolicy")     |

### name

The name of the task, must be unique

`name`

*   is required

*   Type: `string` ([name](pod-defs-task-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/name")

#### name Type

`string` ([name](pod-defs-task-properties-name.md))

### image

Either the container image to run, or a directory containing a Dockerfile. If omitted, the process runs on the host.

`image`

*   is optional

*   Type: `string` ([image](pod-defs-task-properties-image.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-image.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/image")

#### image Type

`string` ([image](pod-defs-task-properties-image.md))

### imagePullPolicy

Pull policy, e.g. Always, Never, IfNotPresent

`imagePullPolicy`

*   is optional

*   Type: `string` ([imagePullPolicy](pod-defs-task-properties-imagepullpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-imagepullpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/imagePullPolicy")

#### imagePullPolicy Type

`string` ([imagePullPolicy](pod-defs-task-properties-imagepullpolicy.md))

### livenessProbe

A probe to check if the task is alive, it will be restarted if not.

`livenessProbe`

*   is optional

*   Type: `object` ([Probe](pod-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/livenessProbe")

#### livenessProbe Type

`object` ([Probe](pod-defs-probe.md))

### readinessProbe

A probe to check if the task is alive, it will be restarted if not.

`readinessProbe`

*   is optional

*   Type: `object` ([Probe](pod-defs-probe.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/readinessProbe")

#### readinessProbe Type

`object` ([Probe](pod-defs-probe.md))

### command

The command to run in the container or on the host. If both the image and the command are omitted, this is a noop.

`command`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/command")

#### command Type

`string[]`

### args

The arguments to pass to the command

`args`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/args")

#### args Type

`string[]`

### workingDir

The working directory in the container or on the host

`workingDir`

*   is optional

*   Type: `string` ([workingDir](pod-defs-task-properties-workingdir.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-workingdir.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/workingDir")

#### workingDir Type

`string` ([workingDir](pod-defs-task-properties-workingdir.md))

### user

The user to run the task as.

`user`

*   is optional

*   Type: `string` ([user](pod-defs-task-properties-user.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-user.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/user")

#### user Type

`string` ([user](pod-defs-task-properties-user.md))

### env

A list of environment variables.

`env`

*   is optional

*   Type: `object[]` ([EnvVar](pod-defs-envvar.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvars.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/env")

#### env Type

`object[]` ([EnvVar](pod-defs-envvar.md))

### ports

A list of ports to expose.

`ports`

*   is optional

*   Type: `object[]` ([Port](pod-defs-port.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-ports.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/ports")

#### ports Type

`object[]` ([Port](pod-defs-port.md))

### volumeMounts

Volumes to mount in the container

`volumeMounts`

*   is optional

*   Type: `object[]` ([VolumeMount](pod-defs-volumemount.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-volumemounts.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/volumeMounts")

#### volumeMounts Type

`object[]` ([VolumeMount](pod-defs-volumemount.md))

### tty

Use a pseudo-TTY

`tty`

*   is optional

*   Type: `boolean` ([tty](pod-defs-task-properties-tty.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-tty.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/tty")

#### tty Type

`boolean` ([tty](pod-defs-task-properties-tty.md))

### watch

A list of files to watch for changes, and restart the task if they change

`watch`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/watch")

#### watch Type

`string[]`

### mutex

A mutex to prevent multiple tasks with the same mutex from running at the same time

`mutex`

*   is optional

*   Type: `string` ([mutex](pod-defs-task-properties-mutex.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-mutex.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/mutex")

#### mutex Type

`string` ([mutex](pod-defs-task-properties-mutex.md))

### semaphore

A semaphore to limit the number of tasks with the same semaphore that can run at the same time

`semaphore`

*   is optional

*   Type: `string` ([semaphore](pod-defs-task-properties-semaphore.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-semaphore.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/semaphore")

#### semaphore Type

`string` ([semaphore](pod-defs-task-properties-semaphore.md))

### dependencies

A list of tasks to run before this task

`dependencies`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-strings.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/dependencies")

#### dependencies Type

`string[]`

### restartPolicy

The restart policy, e.g. Always, Never, OnFailure

`restartPolicy`

*   is optional

*   Type: `string` ([restartPolicy](pod-defs-task-properties-restartpolicy.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task-properties-restartpolicy.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/restartPolicy")

#### restartPolicy Type

`string` ([restartPolicy](pod-defs-task-properties-restartpolicy.md))

## Definitions group Tasks

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Tasks"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group Volume

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/Volume"}
```

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                              |
| :-------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name-3)       | `string` | Required | cannot be null | [Untitled schema](pod-defs-volume-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Volume/properties/name") |
| [hostPath](#hostpath) | `object` | Required | cannot be null | [Untitled schema](pod-defs-hostpath.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Volume/properties/hostPath")           |

### name

Volume's name.

`name`

*   is required

*   Type: `string` ([name](pod-defs-volume-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-volume-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Volume/properties/name")

#### name Type

`string` ([name](pod-defs-volume-properties-name.md))

### hostPath

HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container.

`hostPath`

*   is required

*   Type: `object` ([HostPath](pod-defs-hostpath.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-hostpath.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Volume/properties/hostPath")

#### hostPath Type

`object` ([HostPath](pod-defs-hostpath.md))

## Definitions group VolumeMount

Reference this group by using

```json
{"$ref":"https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount"}
```

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                  |
| :---------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name-4)         | `string` | Required | cannot be null | [Untitled schema](pod-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/name")           |
| [mountPath](#mountpath) | `string` | Required | cannot be null | [Untitled schema](pod-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/mountPath") |

### name

This must match the name of a volume.

`name`

*   is required

*   Type: `string` ([name](pod-defs-volumemount-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/name")

#### name Type

`string` ([name](pod-defs-volumemount-properties-name.md))

### mountPath

Path within the container at which the volume should be mounted.

`mountPath`

*   is required

*   Type: `string` ([mountPath](pod-defs-volumemount-properties-mountpath.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/mountPath")

#### mountPath Type

`string` ([mountPath](pod-defs-volumemount-properties-mountpath.md))
