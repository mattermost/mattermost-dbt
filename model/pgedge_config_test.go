package model_test

import (
	"testing"

	"github.com/mattermost/mattermost-dbt/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClusterConfigFromBytes(t *testing.T) {
	configString := `
{
  "json_version": "1.0",
  "cluster_name": "test1",
  "log_level": "debug",
  "update_date": "2024-1-1 01:01:01 GMT",
  "pgedge": {
    "pg_version": 14,
    "auto_start": "off",
    "spock": {
      "spock_version": "",
      "auto_ddl": "on"
    },
    "databases": [
      {
        "db_name": "mattermost",
        "db_user": "testusername",
        "db_password": "testpassword"
      }
    ]
  },
  "node_groups": [
    {
      "ssh": {
        "os_user": "ubuntu",
        "private_key": ""
      },
      "name": "n1",
      "is_active": "on",
      "public_ip": "127.0.0.1",
      "private_ip": "127.0.0.1",
      "port": "7432",
      "path": "/home/ubuntu/pgedge/cluster/test1/n1",
      "backrest": {
        "stanza": "test1_stanza",
        "repo1-path": "/var/lib/pgbackrest",
        "repo1-retention-full": "7",
        "log-level-console": "info",
        "repo1-cipher-type": "aes-256-cbc"
      }
    },
    {
      "ssh": {
        "os_user": "ubuntu",
        "private_key": ""
      },
      "name": "n2",
      "is_active": "on",
      "public_ip": "127.0.0.1",
      "private_ip": "127.0.0.1",
      "port": "7433",
      "path": "/home/ubuntu/pgedge/cluster/test1/n2",
      "backrest": {
        "stanza": "test1_stanza",
        "repo1-path": "/var/lib/pgbackrest",
        "repo1-retention-full": "7",
        "log-level-console": "info",
        "repo1-cipher-type": "aes-256-cbc"
      }
    }
  ]
}`

	config, err := model.NewClusterConfigFromBytes([]byte(configString))
	require.NoError(t, err)
	require.NotNil(t, config)

	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "test1", config.ClusterName)
		assert.Equal(t, "1.0", config.JSONVersion)
		assert.Equal(t, "2024-1-1 01:01:01 GMT", config.UpdateDate)
		assert.Equal(t, "debug", config.LogLevel)
	})

	t.Run("databases", func(t *testing.T) {
		require.Len(t, config.Pgedge.Databases, 1)
		assert.Equal(t, "mattermost", config.Pgedge.Databases[0].Name)
		assert.Equal(t, "testusername", config.Pgedge.Databases[0].Username)
		assert.Equal(t, "testpassword", config.Pgedge.Databases[0].Password)

		config.Pgedge.Sanitize()
		assert.NotEqual(t, "testpassword", config.Pgedge.Databases[0].Password)
	})

	t.Run("nodes", func(t *testing.T) {
		nodes := config.NodeGroups
		require.Len(t, nodes, 2)
		assert.Equal(t, "n1", nodes[0].Name)
		assert.Equal(t, "n2", nodes[1].Name)
	})
}
