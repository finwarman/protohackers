
Forwarding port example (from local server to remote machine)

```
ssh -R \*:$PORT:localhost:$PORT -o ServerAliveInterval=60 $USER@$REMOTE
```
