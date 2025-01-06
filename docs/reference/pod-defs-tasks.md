# Tasks Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/Tasks
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## Tasks Type

`object` ([Tasks](pod-defs-tasks.md))

# Tasks Properties

| Property | Type     | Required | Nullable       | Defined by                                                                                                                |
| :------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------ |
| `.*`     | `object` | Optional | cannot be null | [Untitled schema](pod-defs-task.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Tasks/patternProperties/.*") |

## Pattern: `.*`

A task is a container or a command to run.

`.*`

*   is optional

*   Type: `object` ([Task](pod-defs-task.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-task.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Tasks/patternProperties/.*")

### .\* Type

`object` ([Task](pod-defs-task.md))
