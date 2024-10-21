package model

import (
	"encoding/json"
	"os"
)

type PgedgeClusterConfig struct {
	JSONVersion string                `json:"json_version"`
	ClusterName string                `json:"cluster_name"`
	LogLevel    string                `json:"log_level"`
	UpdateDate  string                `json:"update_date"`
	Pgedge      *PgedgeDatabaseConfig `json:"pgedge"`
	NodeGroups  []*PgedgeNode         `json:"node_groups"`
}

type PgedgeDatabaseConfig struct {
	PgVersion int               `json:"pg_version"`
	AutoStart string            `json:"auto_start"`
	Spock     *SpockConfig      `json:"spock"`
	Databases []*PgedgeDatabase `json:"databases"`
}

type SpockConfig struct {
	SpockVersion string `json:"spock_version"`
	AutoDDL      string `json:"auto_ddl"`
}

type PgedgeDatabase struct {
	Name     string `json:"db_name"`
	Username string `json:"db_user"`
	Password string `json:"db_password"`
}

type PgedgeNode struct {
	SSH       *PgedgeHostConfig `json:"ssh"`
	Name      string            `json:"name"`
	IsActive  string            `json:"is_active"`
	PublicIP  string            `json:"public_ip"`
	PrivateIP string            `json:"private_ip"`
	Port      string            `json:"port"`
	Path      string            `json:"path"`
	Backrest  *PgedgeBackrest   `json:"backrest"`
}

type PgedgeHostConfig struct {
	OSUser     string `json:"os_user"`
	PrivateKey string `json:"private_key"`
}

type PgedgeBackrest struct {
	Stanza            string `json:"stanza"`
	RepoPath          string `json:"repo1-path"`
	RepoRetentionFull string `json:"repo1-retention-full"`
	LogLevelConsole   string `json:"log-level-console"`
	RepoCypherType    string `json:"repo1-cipher-type"`
}

func (dbc *PgedgeDatabaseConfig) Sanitize() {
	for _, database := range dbc.Databases {
		database.Sanitize()
	}
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
