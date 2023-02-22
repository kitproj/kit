# HTTPGetAction Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/Probe/properties/httpGet
```

HTTPGetAction describes an action based on HTTP Get requests.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## httpGet Type

`object` ([HTTPGetAction](pod-defs-httpgetaction.md))

# httpGet Properties

| Property          | Type      | Required | Nullable       | Defined by                                                                                                                                                |
| :---------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [scheme](#scheme) | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/scheme") |
| [port](#port)     | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/port")     |
| [path](#path)     | `string`  | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/path")     |

## scheme

Scheme to use for connecting to the host. Defaults to HTTP.

`scheme`

*   is optional

*   Type: `string` ([scheme](pod-defs-httpgetaction-properties-scheme.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/scheme")

### scheme Type

`string` ([scheme](pod-defs-httpgetaction-properties-scheme.md))

## port

Number of the port

`port`

*   is optional

*   Type: `integer` ([port](pod-defs-httpgetaction-properties-port.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/port")

### port Type

`integer` ([port](pod-defs-httpgetaction-properties-port.md))

## path

Path to access on the HTTP server.

`path`

*   is optional

*   Type: `string` ([path](pod-defs-httpgetaction-properties-path.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/HTTPGetAction/properties/path")

### path Type

`string` ([path](pod-defs-httpgetaction-properties-path.md))
