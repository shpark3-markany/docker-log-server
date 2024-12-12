package controllers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/labstack/echo/v4"
	log "github.com/markany-inc/docker-log-server/utils"
)

func getLog(containerId string, filters []string, excludes []bool, opts container.LogsOptions) string {
	ctx := context.Background()
	reader, err := cli.ContainerLogs(ctx, containerId, opts)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting logs: %v", err))
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	logs := buf.String()
	if len(filters) > 0 {
		filtered_logs := filterLogs(logs, filters, excludes)
		return filtered_logs
	}
	return logs
}

func filterLogs(logs string, filters []string, excludes []bool) string {
	var filteredLogs string
	scanner := bufio.NewScanner(strings.NewReader(logs))
	log.Info(fmt.Sprintf("Log filtering start FILTER:[%v] EXCLUDE:[%v]", filters, excludes))

	for scanner.Scan() { //새로운 줄을 line에 저장하며 반복
		line := scanner.Text()
		for index, filter := range filters {
			if excludes[index] {
				if !strings.Contains(line, filter) { //filter word가 포함되지 않은 줄만 추가 (grep -v)
					filteredLogs += line + "\n"
				}
			} else {
				if strings.Contains(line, filter) { //filter word가 포함된 줄만 추가 (grep)
					filteredLogs += line + "\n"
				}
			}
		}
	}
	log.Info("Log filtering end")
	return filteredLogs
}

func GetLog(c echo.Context) error {
	params := c.QueryParams()
	containerId := c.Param("id")

	var logOpts container.LogsOptions
	logOpts.ShowStdout = strings.ToLower(params.Get("stdout")) == "true"
	logOpts.ShowStderr = strings.ToLower(params.Get("stderr")) == "true"
	logOpts.Timestamps = strings.ToLower(params.Get("timestamps")) == "true"
	logOpts.Follow = strings.ToLower(params.Get("follow")) == "true"
	logOpts.Tail = params.Get("tail")
	logOpts.Since = params.Get("since")
	logOpts.Until = params.Get("until")

	filters := params["filters"]
	excludes := params["excludes"]
	excludesBool := make([]bool, len(excludes))
	if len(filters) != len(excludes) {
		return c.String(http.StatusBadRequest, "filters and excludes must have the same length")
	}
	for i, exclude := range excludes {
		excludesBool[i] = strings.ToLower(exclude) == "true"
	}

	logs := getLog(containerId, filters, excludesBool, logOpts)
	return c.String(http.StatusOK, logs)
}
