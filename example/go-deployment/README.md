# go-deployment-test

```bash
# start buildkitd
docker run -d --name buildkitdtcp --privileged moby/buildkit:latest --addr tcp://0.0.0.0:1234

# build binary
GOOS=linux GOARCH=arm64 go build .

# start builder
docker run -it --rm -e BUILDKIT_HOST="tcp://buildkitdtcp.orb.local:1234" -v $(pwd):/runner zeabur/alpine-base bash

# (in container) run /runner
$ /runner/go-deployment

# copy docker.tar to /runner and exit

docker import docker.tar
```
