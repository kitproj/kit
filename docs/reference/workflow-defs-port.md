# Port Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/Task/properties/ports/items
```

A port to expose.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## items Type

`object` ([Port](workflow-defs-port.md))

# items Properties

| Property                        | Type      | Required | Nullable       | Defined by                                                                                                                                                      |
| :------------------------------ | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [containerPort](#containerport) | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/containerPort") |
| [hostPort](#hostport)           | `integer` | Optional | cannot be null | [Untitled schema](workflow-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/hostPort")           |

## containerPort

The container port to expose

`containerPort`

*   is optional

*   Type: `integer` ([containerPort](workflow-defs-port-properties-containerport.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/containerPort")

### containerPort Type

`integer` ([containerPort](workflow-defs-port-properties-containerport.md))

## hostPort

The host port to route to the container port

`hostPort`

*   is optional

*   Type: `integer` ([hostPort](workflow-defs-port-properties-hostport.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/Port/properties/hostPort")

### hostPort Type

`integer` ([hostPort](workflow-defs-port-properties-hostport.md))
