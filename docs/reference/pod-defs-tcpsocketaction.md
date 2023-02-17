# TCPSocketAction Schema

```txt
https://github.com/alexec/kit/internal/types/pod#/$defs/TCPSocketAction
```

TCPSocketAction describes an action based on opening a socket

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## TCPSocketAction Type

`object` ([TCPSocketAction](pod-defs-tcpsocketaction.md))

# TCPSocketAction Properties

| Property      | Type      | Required | Nullable       | Defined by                                                                                                                                               |
| :------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [port](#port) | `integer` | Required | cannot be null | [Untitled schema](pod-defs-tcpsocketaction-properties-port.md "https://github.com/alexec/kit/internal/types/pod#/$defs/TCPSocketAction/properties/port") |

## port

Port number of the port to probe.

`port`

*   is required

*   Type: `integer` ([port](pod-defs-tcpsocketaction-properties-port.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-tcpsocketaction-properties-port.md "https://github.com/alexec/kit/internal/types/pod#/$defs/TCPSocketAction/properties/port")

### port Type

`integer` ([port](pod-defs-tcpsocketaction-properties-port.md))
