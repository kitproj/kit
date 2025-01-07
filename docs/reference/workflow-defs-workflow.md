# Workflow Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## Workflow Type

`object` ([Workflow](workflow-defs-workflow.md))

# Workflow Properties

| Property                                                        | Type      | Required | Nullable       | Defined by                                                                                                                                                                                              |
| :-------------------------------------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [terminationGracePeriodSeconds](#terminationgraceperiodseconds) | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-workflow-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/terminationGracePeriodSeconds") |
| [tasks](#tasks)                                                 | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-tasks.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/tasks")                                                                     |
| [volumes](#volumes)                                             | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-workflow-properties-volumes.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/volumes")                                             |
| [semaphores](#semaphores)                                       | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-workflow-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/semaphores")                                       |
| [env](#env)                                                     | `object`  | Optional | cannot be null | [Untitled schema](workflow-defs-envvars.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/env")                                                                     |
| [envfile](#envfile)                                             | `array`   | Optional | cannot be null | [Untitled schema](workflow-defs-envfile.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/envfile")                                                                 |

## terminationGracePeriodSeconds



`terminationGracePeriodSeconds`

*   is optional

*   Type: `integer` ([terminationGracePeriodSeconds](workflow-defs-workflow-properties-terminationgraceperiodseconds.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-workflow-properties-terminationgraceperiodseconds.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/terminationGracePeriodSeconds")

### terminationGracePeriodSeconds Type

`integer` ([terminationGracePeriodSeconds](workflow-defs-workflow-properties-terminationgraceperiodseconds.md))

## tasks



`tasks`

*   is optional

*   Type: `object` ([Tasks](workflow-defs-tasks.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-tasks.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/tasks")

### tasks Type

`object` ([Tasks](workflow-defs-tasks.md))

## volumes



`volumes`

*   is optional

*   Type: `object[]` ([Volume](workflow-defs-volume.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-workflow-properties-volumes.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/volumes")

### volumes Type

`object[]` ([Volume](workflow-defs-volume.md))

## semaphores



`semaphores`

*   is optional

*   Type: `object` ([semaphores](workflow-defs-workflow-properties-semaphores.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-workflow-properties-semaphores.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/semaphores")

### semaphores Type

`object` ([semaphores](workflow-defs-workflow-properties-semaphores.md))

## env

A list of environment variables.

`env`

*   is optional

*   Type: `object` ([EnvVars](workflow-defs-envvars.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-envvars.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/env")

### env Type

`object` ([EnvVars](workflow-defs-envvars.md))

## envfile



`envfile`

*   is optional

*   Type: `string[]`

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-envfile.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/envfile")

### envfile Type

`string[]`
