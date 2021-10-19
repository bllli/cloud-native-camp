# CloudNativeCamp

# demo-server
### Docker
```shell
docker build -t demo-server .
docker run -d -p 8888:8888 --name demo-server demo-server
# 创建dockerhub仓库
docker image tag demo-server blllicn/demo-server
docker push blllicn/demo-server
```
see: https://hub.docker.com/r/blllicn/demo-server
