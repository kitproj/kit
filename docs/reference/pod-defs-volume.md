# Volume Schema

```txt
https://github.com/alexec/kit/internal/types/pod#/$defs/Volume
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## Volume Type

`object` ([Volume](pod-defs-volume.md))

# Volume Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                             |
| :-------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)         | `string` | Required | cannot be null | [Untitled schema](pod-defs-volume-properties-name.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Volume/properties/name") |
| [hostPath](#hostpath) | `object` | Required | cannot be null | [Untitled schema](pod-defs-hostpath.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Volume/properties/hostPath")           |

## name

Volume's name.

`name`

*   is required

*   Type: `string` ([name](pod-defs-volume-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-volume-properties-name.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Volume/properties/name")

### name Type

`string` ([name](pod-defs-volume-properties-name.md))

## hostPath

HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container.

`hostPath`

*   is required

*   Type: `object` ([HostPath](pod-defs-hostpath.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-hostpath.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Volume/properties/hostPath")

### hostPath Type

`object` ([HostPath](pod-defs-hostpath.md))
