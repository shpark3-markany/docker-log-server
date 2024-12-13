package controllers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/labstack/echo/v4"
)

func ListContainers(c echo.Context) {
	ctx := context.Background()
	cli.ContainerList(ctx, container.ListOptions{})
}
