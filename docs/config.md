# RemoteSync Configuration Guide

## Server Configuration

### Network Settings
```json
{
    "port": "8080",
    "max_clients": 2,
    "tls_enabled": true,
    "cert_file": "server.crt",
    "key_file": "server.key"
}
```

### Security Settings
```json
{
    "auth_required": true,
    "token_expiry": "24h",
    "allowed_ips": ["192.168.1.*"]
}
```

## Client Configuration

### Connection Settings
```json
{
    "server_host": "localhost",
    "server_port": "8080",
    "retry_delay": "5s",
    "max_retries": 5
}
```

### Performance Settings
```json
{
    "frame_rate": 30,
    "quality": 75,
    "max_bandwidth": 5242880
}
```