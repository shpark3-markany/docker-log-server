package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	log "github.com/markany-inc/docker-log-server/utils"
)

var (
	cli      *client.Client
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func MakeClient() {
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Sprintf("Error creating docker client: %v", r))
		}
	}()
	if _, err := os.Stat("/var/run/docker.sock"); os.IsNotExist(err) {
		log.Error("/var/run/docker.sock does not exist")
		return
	}

	cli_, err := client.NewClientWithOpts(client.FromEnv, client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(fmt.Sprintf("Error creating docker client: %v", err))
		return
	}
	log.Info("Docker client created")
	cli = cli_
}
