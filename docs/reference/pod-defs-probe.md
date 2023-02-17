# Probe Schema

```txt
https://github.com/alexec/kit/internal/types/pod#/$defs/Probe
```

A probe to check if the task is alive, it will be restarted if not.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [pod.schema.json\*](../../out/pod.schema.json "open original schema") |

## Probe Type

`object` ([Probe](pod-defs-probe.md))

# Probe Properties

| Property                                    | Type      | Required | Nullable       | Defined by                                                                                                                                                         |
| :------------------------------------------ | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [tcpSocket](#tcpsocket)                     | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-tcpsocketaction.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/tcpSocket")                                |
| [httpGet](#httpget)                         | `object`  | Optional | cannot be null | [Untitled schema](pod-defs-httpgetaction.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/httpGet")                                    |
| [initialDelaySeconds](#initialdelayseconds) | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-initialdelayseconds.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/initialDelaySeconds") |
| [periodSeconds](#periodseconds)             | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-periodseconds.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/periodSeconds")             |
| [successThreshold](#successthreshold)       | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-successthreshold.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/successThreshold")       |
| [failureThreshold](#failurethreshold)       | `integer` | Optional | cannot be null | [Untitled schema](pod-defs-probe-properties-failurethreshold.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/failureThreshold")       |

## tcpSocket

TCPSocketAction describes an action based on opening a socket

`tcpSocket`

*   is optional

*   Type: `object` ([TCPSocketAction](pod-defs-tcpsocketaction.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-tcpsocketaction.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/tcpSocket")

### tcpSocket Type

`object` ([TCPSocketAction](pod-defs-tcpsocketaction.md))

## httpGet

HTTPGetAction describes an action based on HTTP Get requests.

`httpGet`

*   is optional

*   Type: `object` ([HTTPGetAction](pod-defs-httpgetaction.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-httpgetaction.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/httpGet")

### httpGet Type

`object` ([HTTPGetAction](pod-defs-httpgetaction.md))

## initialDelaySeconds

Number of seconds after the process has started before the probe is initiated.

`initialDelaySeconds`

*   is optional

*   Type: `integer` ([initialDelaySeconds](pod-defs-probe-properties-initialdelayseconds.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-initialdelayseconds.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/initialDelaySeconds")

### initialDelaySeconds Type

`integer` ([initialDelaySeconds](pod-defs-probe-properties-initialdelayseconds.md))

## periodSeconds

How often (in seconds) to perform the probe.

`periodSeconds`

*   is optional

*   Type: `integer` ([periodSeconds](pod-defs-probe-properties-periodseconds.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-periodseconds.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/periodSeconds")

### periodSeconds Type

`integer` ([periodSeconds](pod-defs-probe-properties-periodseconds.md))

## successThreshold

Minimum consecutive successes for the probe to be considered successful after having failed.

`successThreshold`

*   is optional

*   Type: `integer` ([successThreshold](pod-defs-probe-properties-successthreshold.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-successthreshold.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/successThreshold")

### successThreshold Type

`integer` ([successThreshold](pod-defs-probe-properties-successthreshold.md))

## failureThreshold

Minimum consecutive failures for the probe to be considered failed after having succeeded.

`failureThreshold`

*   is optional

*   Type: `integer` ([failureThreshold](pod-defs-probe-properties-failurethreshold.md))

*   cannot be null

*   defined in: [Untitled schema](pod-defs-probe-properties-failurethreshold.md "https://github.com/alexec/kit/internal/types/pod#/$defs/Probe/properties/failureThreshold")

### failureThreshold Type

`integer` ([failureThreshold](pod-defs-probe-properties-failurethreshold.md))
