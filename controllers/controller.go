package controllers

import (
	"log"

	"github.com/docker/docker/client"
)

var (
	cli *client.Client
)

func makeClient() {
	cli_, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	cli = cli_
}

func init() {
	makeClient()
}
