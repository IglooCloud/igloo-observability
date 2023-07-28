ssh -t igloo-observability 'echo "== Connecting to remote server ==" \
    && cd ~/igloo-observability \
    && echo "== Fetching latest version from git ==" \
    && git pull \
    && echo "== Building the application ==" \
    && go build -o build/server ./cmd \
    && echo "== Restarting systemctl service ==" \
    && sudo systemctl restart observability-server.service \
    && echo "== Deployed successfully =="'