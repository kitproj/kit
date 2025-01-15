# type Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/type
```

Type is the type of the task: "service" or "job". If omitted, if there are ports, it's a service, otherwise it's a job.
This is only needed when you have service that does not listen on ports.
Services are running in the background.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## type Type

`string` ([type](workflow-defs-task-properties-type.md))
