# Protohackers Solutions

## Notes / Helpers

Forwarding port example (from local server to remote machine)

```bash
ssh -R \*:$PORT:localhost:$PORT -o ServerAliveInterval=60 $USER@$REMOTE
```

## JSON Parser
