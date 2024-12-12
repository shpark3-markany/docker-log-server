package models

import "github.com/docker/docker/api/types/container"

type Configs struct {
	LogOptions  LogOptions `json:"logOptions"`  //LogOptions 구조체
	ContainerId string     `json:"containerId"` //ContainerId 문자열
}
type (
	LogOptions struct {
		Option  container.LogsOptions `json:"option"`  //ContainerLogs 옵션
		Exclude bool                  `json:"exclude"` //Exclude 사용 여부
		Filter  string                `json:"filter"`  //Filter 단어
	}
)
