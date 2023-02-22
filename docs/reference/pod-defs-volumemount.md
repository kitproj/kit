# VolumeMount Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount
```

VolumeMount describes a mounting of a Volume within a container.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## VolumeMount Type

`object` ([VolumeMount](pod-defs-volumemount.md))

# VolumeMount Properties

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                  |
| :---------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)           | `string` | Required | cannot be null | [Untitled schema](pod-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/name")           |
| [mountPath](#mountpath) | `string` | Required | cannot be null | [Untitled schema](pod-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/mountPath") |

## name

This must match the name of a volume.

`name`

*   is required

*   Type: `string` ([name](pod-defs-volumemount-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-volumemount-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/name")

### name Type

`string` ([name](pod-defs-volumemount-properties-name.md))

## mountPath

Path within the container at which the volume should be mounted.

`mountPath`

*   is required

*   Type: `string` ([mountPath](pod-defs-volumemount-properties-mountpath.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-volumemount-properties-mountpath.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/VolumeMount/properties/mountPath")

### mountPath Type

`string` ([mountPath](pod-defs-volumemount-properties-mountpath.md))
