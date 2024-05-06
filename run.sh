PORT=25565
USER=
REMOTE=
ssh -R \*:$PORT:localhost:$PORT -o ServerAliveInterval=60 $USER@$REMOTE
