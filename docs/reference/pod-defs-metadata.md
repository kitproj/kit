# Metadata Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/Pod/properties/metadata
```

Metadata is the metadata for the pod.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## metadata Type

`object` ([Metadata](pod-defs-metadata.md))

# metadata Properties

| Property                    | Type     | Required | Nullable       | Defined by                                                                                                                                                |
| :-------------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)               | `string` | Optional | cannot be null | [Untitled schema](pod-defs-metadata-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/name")               |
| [annotations](#annotations) | `object` | Optional | cannot be null | [Untitled schema](pod-defs-metadata-properties-annotations.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/annotations") |

## name

Name is the name of the resource.

`name`

*   is optional

*   Type: `string` ([name](pod-defs-metadata-properties-name.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-metadata-properties-name.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/name")

### name Type

`string` ([name](pod-defs-metadata-properties-name.md))

## annotations

Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata.

`annotations`

*   is optional

*   Type: `object` ([annotations](pod-defs-metadata-properties-annotations.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-metadata-properties-annotations.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Metadata/properties/annotations")

### annotations Type

`object` ([annotations](pod-defs-metadata-properties-annotations.md))
