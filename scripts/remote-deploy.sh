echo "== Building binaries =="
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 \
    go build -o build/server -ldflags "-linkmode external -extldflags -static" ./cmd/server
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 \
    go build -o build/os-instrumentation -ldflags "-linkmode external -extldflags -static" ./cmd/instrumentation/os

echo "== Copying binaries to remote server =="
scp build/* igloo-observability:~/igloo-observability
scp init/* igloo-observability:~/igloo-observability

ssh -t igloo-observability 'echo "== Connecting to remote server ==" \
    && echo "== Restarting systemctl service ==" \
    && sudo systemctl restart observability-server.service \
    && sudo systemctl restart observability-os.service \
    && echo "== Deployed successfully =="'