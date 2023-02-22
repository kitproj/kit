# Port Schema

```txt
https://github.com/kitproj/kit/internal/types/pod#/$defs/Ports/items
```

A port to expose.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## items Type

`object` ([Port](pod-defs-port.md))

# items Properties

| Property                        | Type      | Required | Nullable       | Defined by                                                                                                                                            |
| :------------------------------ | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
| [containerPort](#containerport) | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/containerPort") |
| [hostPort](#hostport)           | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/hostPort")           |

## containerPort

The container port to expose

`containerPort`

*   is optional

*   Type: `integer` ([containerPort](pod-defs-port-properties-containerport.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-port-properties-containerport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/containerPort")

### containerPort Type

`integer` ([containerPort](pod-defs-port-properties-containerport.md))

## hostPort

The host port to route to the container port

`hostPort`

*   is optional

*   Type: `integer` ([hostPort](pod-defs-port-properties-hostport.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-port-properties-hostport.md "https://github.com/kitproj/kit/internal/types/pod#/$defs/Port/properties/hostPort")

### hostPort Type

`integer` ([hostPort](pod-defs-port-properties-hostport.md))
