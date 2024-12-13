package controllers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/markany-inc/docker-log-server/models"
	log "github.com/markany-inc/docker-log-server/utils"
)

func getLog(containerId string, filters []string, excludes []bool, opts container.LogsOptions) ([]models.LogEntry, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Sprintf("RECOVERED: Error getting logs: %v", r))
		}
	}()

	var logs []models.LogEntry
	ctx := context.Background()
	reader, err := cli.ContainerLogs(ctx, containerId, opts)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting logs: %v", err))
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	log_string := buf.String()

	if len(filters) > 0 {
		logs, err = filterLogs(log_string, filters, excludes)
		if err != nil {
			log.Error(fmt.Sprintf("Error filtering logs: %v", err))
		}
	} else {
		scanner := bufio.NewScanner(strings.NewReader(log_string))
		for scanner.Scan() {
			line := scanner.Text()
			logs = append(logs, models.LogEntry{Line: line})
		}
	}

	return logs, nil
}

func filterLogs(log_string string, filters []string, excludes []bool) ([]models.LogEntry, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Sprintf("RECOVERED: Error getting logs: %v", r))
		}
	}()

	var logs []models.LogEntry
	scanner := bufio.NewScanner(strings.NewReader(log_string))
	log.Info(fmt.Sprintf("Log filtering start FILTER:[%v] EXCLUDE:[%v]", filters, excludes))

	for scanner.Scan() { //새로운 줄을 line에 저장하며 반복
		line := scanner.Text()
		for index, filter := range filters {
			if excludes[index] {
				if !strings.Contains(line, filter) { //filter word가 포함되지 않은 줄만 추가 (grep -v)
					logs = append(logs, models.LogEntry{Line: line})
				}
			} else {
				if strings.Contains(line, filter) { //filter word가 포함된 줄만 추가 (grep)
					logs = append(logs, models.LogEntry{Line: line})
				}
			}
		}
	}
	log.Info(fmt.Sprintf("Log filtering end. Total logs: %v", len(logs)))
	return logs, nil
}

func streamLogs(c echo.Context, containerId string, filters []string, excludes []bool, opts container.LogsOptions) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to set websocket upgrade: %v", err))
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	reader, err := cli.ContainerLogs(ctx, containerId, opts)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting logs: %v", err))
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		for index, filter := range filters {
			if excludes[index] {
				if !strings.Contains(line, filter) {
					if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
						log.Error(fmt.Sprintf("Error writing message: %v", err))
						return err
					}
				}
			} else {
				if strings.Contains(line, filter) {
					if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
						log.Error(fmt.Sprintf("Error writing message: %v", err))
						return err
					}
				}
			}
		}
	}
	return nil
}

func GetLog(c echo.Context) error {
	params := c.QueryParams()
	containerId := params.Get("id")
	if containerId == "" {
		return c.JSON(http.StatusBadRequest, "id is required")
	}
	log.Info(fmt.Sprintf("GetLog: getting logs for container: %s", containerId))

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
		return c.JSON(http.StatusBadRequest, "filters and excludes must have the same length")
	}
	for i, exclude := range excludes {
		excludesBool[i] = strings.ToLower(exclude) == "true"
	}

	if logOpts.Follow {
		log.Info("GetLog: sending log through websocket")
		return streamLogs(c, containerId, filters, excludesBool, logOpts)
	}
	logs, err := getLog(containerId, filters, excludesBool, logOpts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error getting logs")
	}
	log.Info("GetLog: sending log successfuly end")
	return c.JSON(http.StatusOK, logs)
}
