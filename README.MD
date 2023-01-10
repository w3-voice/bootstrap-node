## Build
docker build -t hoodboot .
## Run
docker run -d --name hoodboot --restart always --net=host -v $(pwd)/data:/hbnode/data --env GOLOG_LOG_LEVEL="warn,boothood=debug" hoodboot:latest