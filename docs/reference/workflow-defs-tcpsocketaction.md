# TCPSocketAction Schema

```txt
https://github.com/kitproj/kit/internal/types/workflow#/$defs/TCPSocketAction
```

TCPSocketAction describes an action based on opening a socket

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [workflow.schema.json\*](../../out/workflow.schema.json "open original schema") |

## TCPSocketAction Type

`object` ([TCPSocketAction](workflow-defs-tcpsocketaction.md))

# TCPSocketAction Properties

| Property      | Type      | Required | Nullable       | Defined by                                                                                                                                                          |
| :------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [port](#port) | `integer` | Required | cannot be null | [Untitled schema](workflow-defs-tcpsocketaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/TCPSocketAction/properties/port") |

## port

Port number of the port to probe.

`port`

*   is required

*   Type: `integer` ([port](workflow-defs-tcpsocketaction-properties-port.md))

*   cannot be null

*   defined in: [Untitled schema](workflow-defs-tcpsocketaction-properties-port.md "https://github.com/kitproj/kit/internal/types/workflow#/$defs/TCPSocketAction/properties/port")

### port Type

`integer` ([port](workflow-defs-tcpsocketaction-properties-port.md))
