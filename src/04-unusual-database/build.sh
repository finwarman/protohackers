# Build for centos8
GOOS=linux GOARCH=amd64 go build -o 04-udp-db
# Copy to remote
scp 04-udp-db maughan:~
