# EnvVar Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/Task/properties/env/items
```

A environment variable.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## items Type

`object` ([EnvVar](pod-defs-envvar.md))

# items Properties

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                |
| :---------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)           | `string` | Required | cannot be null | [Untitled schema](pod-defs-envvar-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/name")   |
| [value](#value)         | `string` | Required | cannot be null | [Untitled schema](pod-defs-envvar-properties-value.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/value") |
| [valueFrom](#valuefrom) | `object` | Optional | cannot be null | [Untitled schema](pod-defs-envvarsource.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/valueFrom")        |

## name



`name`

*   is required

*   Type: `string` ([name](pod-defs-envvar-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvar-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/name")

### name Type

`string` ([name](pod-defs-envvar-properties-name.md))

## value



`value`

*   is required

*   Type: `string` ([value](pod-defs-envvar-properties-value.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvar-properties-value.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/value")

### value Type

`string` ([value](pod-defs-envvar-properties-value.md))

## valueFrom



`valueFrom`

*   is optional

*   Type: `object` ([EnvVarSource](pod-defs-envvarsource.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-envvarsource.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/EnvVar/properties/valueFrom")

### valueFrom Type

`object` ([EnvVarSource](pod-defs-envvarsource.md))
