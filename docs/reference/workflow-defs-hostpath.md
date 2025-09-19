# HostPath Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Volume/properties/hostPath
```

HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## hostPath Type

`object` ([HostPath](workflow-defs-hostpath.md))

# hostPath Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                            |
| :------------ | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
| [path](#path) | `string` | Required | cannot be null | [Untitled schema](workflow-defs-hostpath-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HostPath/properties/path") |

## path

Path of the directory on the host.

`path`

* is required

* Type: `string` ([path](workflow-defs-hostpath-properties-path.md))

* cannot be null

* defined in: [Untitled schema](workflow-defs-hostpath-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HostPath/properties/path")

### path Type

`string` ([path](workflow-defs-hostpath-properties-path.md))
