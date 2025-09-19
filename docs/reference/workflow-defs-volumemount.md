# VolumeMount Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount
```

VolumeMount describes a mounting of a Volume within a container.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## VolumeMount Type

`object` ([VolumeMount](workflow-defs-volumemount.md))

# VolumeMount Properties

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                            |
| :---------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)           | `string` | Required | cannot be null | [Untitled schema](workflow-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/name")           |
| [mountPath](#mountpath) | `string` | Required | cannot be null | [Untitled schema](workflow-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/mountPath") |

## name

This must match the name of a volume.

`name`

* is required

* Type: `string` ([name](workflow-defs-volumemount-properties-name.md))

* cannot be null

* defined in: [Untitled schema](workflow-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/name")

### name Type

`string` ([name](workflow-defs-volumemount-properties-name.md))

## mountPath

Path within the container at which the volume should be mounted.

`mountPath`

* is required

* Type: `string` ([mountPath](workflow-defs-volumemount-properties-mountpath.md))

* cannot be null

* defined in: [Untitled schema](workflow-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/VolumeMount/properties/mountPath")

### mountPath Type

`string` ([mountPath](workflow-defs-volumemount-properties-mountpath.md))
