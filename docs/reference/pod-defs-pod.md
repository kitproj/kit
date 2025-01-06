# Pod Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## Pod Type

`object` ([Pod](pod-defs-pod.md))

# Pod Properties

| Property                                                        | Type      | Required | Nullable       | Defined by                                                                                                                                                                          |
| :-------------------------------------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [terminationGracePeriodSeconds](#terminationgraceperiodseconds) | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-pod-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/terminationGracePeriodSeconds") |
| [tasks](#tasks)                                                 | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-tasks.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/tasks")                                                                |
| [volumes](#volumes)                                             | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-pod-properties-volumes.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/volumes")                                             |
| [semaphores](#semaphores)                                       | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-pod-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/semaphores")                                       |
| [env](#env)                                                     | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-envvars.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/env")                                                                |
| [envfile](#envfile)                                             | `array`   | Optional | cannot be null | [Untitled schema](pod-defs-envfile.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/envfile")                                                            |

## terminationGracePeriodSeconds



`terminationGracePeriodSeconds`

*   is optional

*   Type: `integer` ([terminationGracePeriodSeconds](pod-defs-pod-properties-terminationgraceperiodseconds.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-pod-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/terminationGracePeriodSeconds")

### terminationGracePeriodSeconds Type

`integer` ([terminationGracePeriodSeconds](pod-defs-pod-properties-terminationgraceperiodseconds.md))

## tasks



`tasks`

*   is optional

*   Type: `object` ([Tasks](pod-defs-tasks.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-tasks.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/tasks")

### tasks Type

`object` ([Tasks](pod-defs-tasks.md))

## volumes



`volumes`

*   is optional

*   Type: `object[]` ([Volume](pod-defs-volume.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-pod-properties-volumes.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/volumes")

### volumes Type

`object[]` ([Volume](pod-defs-volume.md))

## semaphores



`semaphores`

*   is optional

*   Type: `object` ([semaphores](pod-defs-pod-properties-semaphores.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-pod-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/semaphores")

### semaphores Type

`object` ([semaphores](pod-defs-pod-properties-semaphores.md))

## env

A list of environment variables.

`env`

*   is optional

*   Type: `object` ([EnvVars](pod-defs-envvars.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvars.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/env")

### env Type

`object` ([EnvVars](pod-defs-envvars.md))

## envfile



`envfile`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envfile.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/envfile")

### envfile Type

`string[]`
