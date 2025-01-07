# Tasks Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/tasks
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## tasks Type

`object` ([Tasks](workflow-defs-tasks.md))

# tasks Properties

| Property | Type     | Required | Nullable       | Defined by                                                                                                                          |
| :------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------- |
| `.*`     | `object` | Optional | cannot be null | [Untitled schema](workflow-defs-task.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Tasks/patternProperties/.*") |

## Pattern: `.*`

A task is a container or a command to run.

`.*`

*   is optional

*   Type: `object` ([Task](workflow-defs-task.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-task.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Tasks/patternProperties/.*")

### .\* Type

`object` ([Task](workflow-defs-task.md))
