# RemoteSync Protocol Documentation

## Message Structure
All messages follow this JSON format:
```json
{
    "type": "message_type",
    "timestamp": "2024-04-03T12:00:00Z",
    "payload": {},
    "sender_id": "client_id"
}
```

## Message Types

### Identification (identify)
Used when a client connects to the server.
```json
{
    "type": "identify",
    "payload": {
        "name": "client_name",
        "hostname": "machine_hostname"
    }
}
```

### Heartbeat (ping/pong)
Used to maintain connection status.
```json
{
    "type": "ping",
    "payload": {
        "timestamp": "2024-04-03T12:00:00Z"
    }
}
```

### Control Request (control_request)
Used to request control of another client.
```json
{
    "type": "control_request",
    "payload": {
        "target_id": "target_client_id",
        "mode": "exclusive|simultaneous"
    }
}
```

### Mouse/Keyboard Events
Used to transmit input events.
```json
{
    "type": "mouse_event",
    "payload": {
        "x": 100,
        "y": 200,
        "button": 1,
        "action": "click|move|scroll"
    }
}
```