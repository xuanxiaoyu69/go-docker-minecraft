package main

import (
	"context"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type ContainerInfo struct {
	ID     string
	Image  string
	Status string
}

func main() {
	r := gin.Default()

	r.Static("/static", "./public")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", welcomeHandler)
	r.GET("/containers", containerHandler)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func welcomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", nil)
}

func containerHandler(c *gin.Context) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	containerInfo := make([]ContainerInfo, len(containers))
	for i, container := range containers {
		containerInfo[i] = ContainerInfo{
			ID:     container.ID[:12],
			Image:  container.Image,
			Status: container.Status,
		}
	}

	c.HTML(http.StatusOK, "containers.html", gin.H{
		"Containers": containerInfo,
	})
}
