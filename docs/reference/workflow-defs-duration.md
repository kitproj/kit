# Duration Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/stalledTimeout
```

The timeout for the task to be considered stalled. If omitted, the task will be considered stalled after 30 seconds of no activity.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## stalledTimeout Type

`object` ([Duration](workflow-defs-duration.md))

# stalledTimeout Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                |
| :-------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
| [Duration](#duration) | `object` | Required | cannot be null | [Untitled schema](workflow-defs-duration.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Duration/properties/Duration") |

## Duration



`Duration`

* is required

* Type: `object` ([Duration](workflow-defs-duration.md))

* cannot be null

* defined in: [Untitled schema](workflow-defs-duration.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Duration/properties/Duration")

### Duration Type

`object` ([Duration](workflow-defs-duration.md))
