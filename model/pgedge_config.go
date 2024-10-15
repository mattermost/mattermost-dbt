package model

import (
	"encoding/json"
	"os"
)

type PgedgeClusterConfig struct {
	Name       string                 `json:"name"`
	Style      string                 `json:"style"`
	CreateDate string                 `json:"create_date"`
	Remote     PgedgeHostConfig       `json:"remote,omitempty"`
	LocalHost  PgedgeHostConfig       `json:"localhost,omitempty"`
	Database   PgedgeDatabaseConfig   `json:"database"`
	NodeGroups PgedgeNodeGroupsConfig `json:"node_groups"`
}

type PgedgeHostConfig struct {
	OSUser string `json:"os_user"`
	SSHKey string `json:"ssh_key"`
}

type PgedgeDatabaseConfig struct {
	Databases []*PgedgeDatabase `json:"databases"`
	PgVersion int               `json:"pg_version"`
	AutoDDL   string            `json:"auto_ddl"`
	AutoStart string            `json:"auto_start"`
}

type PgedgeDatabase struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type PgedgeNodeGroupsConfig struct {
	Remote    []PgedgeLocalhostNodes `json:"remote,omitempty"`
	Localhost []PgedgeLocalhostNodes `json:"localhost,omitempty"`
}

type PgedgeLocalhostNodes struct {
	Nodes []PgedgeNode `json:"nodes"`
}

type PgedgeNode struct {
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Path      string `json:"path"`
}

func (dbc *PgedgeDatabaseConfig) Sanitize() {
	for _, database := range dbc.Databases {
		database.Sanitize()
	}
}

func (ng *PgedgeNodeGroupsConfig) Nodes() []PgedgeNode {
	var nodes []PgedgeNode

	if len(ng.Remote) != 0 {
		for _, l := range ng.Remote {
			nodes = append(nodes, l.Nodes...)
		}
		return nodes
	}

	for _, l := range ng.Localhost {
		nodes = append(nodes, l.Nodes...)
	}

	return nodes
}

func (db *PgedgeDatabase) Sanitize() {
	db.Password = "****************"
}

func NewClusterConfigFromBytes(bytes []byte) (*PgedgeClusterConfig, error) {
	var config PgedgeClusterConfig
	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewClusterConfigFromFile(file string) (*PgedgeClusterConfig, error) {
	rawConfig, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return NewClusterConfigFromBytes(rawConfig)
}
