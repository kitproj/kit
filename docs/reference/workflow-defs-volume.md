# Volume Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Workflow/properties/volumes/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## items Type

`object` ([Volume](workflow-defs-volume.md))

# items Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                        |
| :-------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| [name](#name)         | `string` | Required | cannot be null | [Untitled schema](workflow-defs-volume-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/name") |
| [hostPath](#hostpath) | `object` | Required | cannot be null | [Untitled schema](workflow-defs-hostpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/hostPath")           |

## name

Volume's name.

`name`

*   is required

*   Type: `string` ([name](workflow-defs-volume-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-volume-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/name")

### name Type

`string` ([name](workflow-defs-volume-properties-name.md))

## hostPath

HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container.

`hostPath`

*   is required

*   Type: `object` ([HostPath](workflow-defs-hostpath.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-hostpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/hostPath")

### hostPath Type

`object` ([HostPath](workflow-defs-hostpath.md))
