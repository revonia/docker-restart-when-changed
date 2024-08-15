# Docker Restart When Changed

`docker-restart-when-changed` 是一个基于 Go 编写的实用工具，用于监控指定文件的变化，并在文件发生变化时自动重启关联的 Docker 容器。该工具运行于 Docker 环境，通过使用 Docker 的 Unix Socket 直接与 Docker Daemon 通信来执行容器的重启操作。由 ChatGPT 编写。

## 功能

- 监控一个或多个文件的变化
- 自动重启与变化文件关联的 Docker 容器
- 通过环境变量进行文件与容器的动态配置

## 使用方式

### 镜像拉取

你可以通过以下命令从 GitHub Container Registry 拉取镜像：

```bash
docker pull ghcr.io/revonia/docker-restart-when-changed:latest
```

### 启动容器
使用该镜像时，你可以指定环境变量 FILE_CONTAINER_PAIRS 来配置文件-容器的映射对。每个映射对使用 file:container 形式，通过逗号分隔多个映射。

```bash
docker run -d \
  -v /path/on/host/to/your/file1:/app/watched_file1 \
  -v /path/on/host/to/your/file2:/app/watched_file2 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e FILE_CONTAINER_PAIRS="/app/watched_file1:container1,/app/watched_file2:container2" \
  ghcr.io/your-username/docker-restart-when-changed:latest
```

### 参数说明
 * -v /path/on/host/to/your/file:/app/watched_file：将主机上的文件挂载到容器中进行监控。
 * -v /var/run/docker.sock:/var/run/docker.sock：挂载 Docker Unix socket，使容器能够与 Docker Daemon 通信。
 * -e FILE_CONTAINER_PAIRS：环境变量，用于指定文件-容器对，每对格式为 file:container，多个对之间用逗号分隔。

