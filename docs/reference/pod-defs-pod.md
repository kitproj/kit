# Pod Schema

```txt
https://github.com/alexec/kit/internal/types/pod#/$defs/Pod
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## Pod Type

`object` ([Pod](pod-defs-pod.md))

# Pod Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                   |
| :------------------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------- |
| [spec](#spec)             | `object` | Required | cannot be null | [Untitled schema](pod-defs-podspec.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/spec")                         |
| [apiVersion](#apiversion) | `string` | Optional | cannot be null | [Untitled schema](pod-defs-pod-properties-apiversion.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/apiVersion") |
| [kind](#kind)             | `string` | Optional | cannot be null | [Untitled schema](pod-defs-pod-properties-kind.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/kind")             |
| [metadata](#metadata)     | `object` | Required | cannot be null | [Untitled schema](pod-defs-metadata.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/metadata")                    |

## spec

Task is a unit of work that should be run.

`spec`

*   is required

*   Type: `object` ([PodSpec](pod-defs-podspec.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-podspec.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/spec")

### spec Type

`object` ([PodSpec](pod-defs-podspec.md))

## apiVersion

APIVersion must be `kit/v1`.

`apiVersion`

*   is optional

*   Type: `string` ([apiVersion](pod-defs-pod-properties-apiversion.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-pod-properties-apiversion.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/apiVersion")

### apiVersion Type

`string` ([apiVersion](pod-defs-pod-properties-apiversion.md))

## kind

Kind must be `Tasks`.

`kind`

*   is optional

*   Type: `string` ([kind](pod-defs-pod-properties-kind.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-pod-properties-kind.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/kind")

### kind Type

`string` ([kind](pod-defs-pod-properties-kind.md))

## metadata

Metadata is the metadata for the pod.

`metadata`

*   is required

*   Type: `object` ([Metadata](pod-defs-metadata.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-metadata.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Pod/properties/metadata")

### metadata Type

`object` ([Metadata](pod-defs-metadata.md))
