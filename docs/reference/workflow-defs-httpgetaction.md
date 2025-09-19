# HTTPGetAction Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Probe/properties/httpGet
```

HTTPGetAction describes an action based on HTTP Locks requests.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## httpGet Type

`object` ([HTTPGetAction](workflow-defs-httpgetaction.md))

# httpGet Properties

| Property          | Type      | Required | Nullable       | Defined by                                                                                                                                                          |
| :---------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [scheme](#scheme) | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/scheme") |
| [port](#port)     | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/port")     |
| [path](#path)     | `string`  | Optional | cannot be null | [Untitled schema](workflow-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/path")     |

## scheme

Scheme to use for connecting to the host. Defaults to HTTP.

`scheme`

* is optional

* Type: `string` ([scheme](workflow-defs-httpgetaction-properties-scheme.md))

* cannot be null

* defined in: [Untitled schema](workflow-defs-httpgetaction-properties-scheme.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/scheme")

### scheme Type

`string` ([scheme](workflow-defs-httpgetaction-properties-scheme.md))

## port

Number of the port

`port`

* is optional

* Type: `integer` ([port](workflow-defs-httpgetaction-properties-port.md))

* cannot be null

* defined in: [Untitled schema](workflow-defs-httpgetaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/port")

### port Type

`integer` ([port](workflow-defs-httpgetaction-properties-port.md))

## path

Path to access on the HTTP server.

`path`

* is optional

* Type: `string` ([path](workflow-defs-httpgetaction-properties-path.md))

* cannot be null

* defined in: [Untitled schema](workflow-defs-httpgetaction-properties-path.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/HTTPGetAction/properties/path")

### path Type

`string` ([path](workflow-defs-httpgetaction-properties-path.md))
