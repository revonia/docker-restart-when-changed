package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Docker Unix socket path
const dockerSocket = "/var/run/docker.sock"

// 重启容器的函数，通过 Unix socket 发送 HTTP POST 请求
func restartContainer(containerName string) error {
	// 创建一个通过 Unix socket 连接的 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", dockerSocket)
			},
		},
		Timeout: 10 * time.Second,
	}

	// 发送POST请求重启容器
	url := fmt.Sprintf("http://localhost/containers/%s/restart", containerName)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to restart container %s, status code: %d", containerName, resp.StatusCode)
	}

	fmt.Printf("Container %s restarted successfully.\n", containerName)
	return nil
}

func main() {
	// 从环境变量中读取文件-容器对的映射，格式为 "file1:container1,file2:container2"
	fileContainerPairs := os.Getenv("FILE_CONTAINER_PAIRS")
	if fileContainerPairs == "" {
		fmt.Println("Error: Missing environment variable FILE_CONTAINER_PAIRS.")
		os.Exit(1)
	}

	// 解析文件-容器对
	pairs := strings.Split(fileContainerPairs, ",")
	fileToContainer := make(map[string]string)
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			fmt.Printf("Invalid pair: %s\n", pair)
			os.Exit(1)
		}
		fileToContainer[parts[0]] = parts[1]
	}

	// 创建 fsnotify 监控器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		os.Exit(1)
	}
	defer watcher.Close()

	// 监听每个文件
	for file := range fileToContainer {
		err = watcher.Add(file)
		if err != nil {
			fmt.Printf("Error adding file %s to watcher: %v\n", file, err)
			os.Exit(1)
		}
		fmt.Printf("Watching file: %s for container: %s\n", file, fileToContainer[file])
	}

	// 事件循环
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// 处理文件修改事件
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf("File %s changed.\n", event.Name)
				containerName := fileToContainer[event.Name]
				err := restartContainer(containerName)
				if err != nil {
					fmt.Printf("Error restarting container: %v\n", err)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
	}
}
