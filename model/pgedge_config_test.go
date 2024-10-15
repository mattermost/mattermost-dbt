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
    "name": "test1",
    "style": "localhost",
    "create_date": "2024-01-01",
    "localhost": {
      "os_user": "ubuntu",
      "ssh_key": ""
    },
    "database": {
      "databases": [
        {
          "username": "testusername",
          "password": "testpassword",
          "name": "mattermost"
        }
      ],
      "pg_version": 14,
      "auto_ddl": "on"
    },
    "node_groups": {
      "localhost": [
        {
          "nodes": [
            {
              "name": "n1",
              "is_active": true,
              "ip_address": "127.0.0.1",
              "port": 7432,
              "path": "/home/ubuntu/pgedge/cluster/test1/n1"
            }
          ]
        },
        {
          "nodes": [
            {
              "name": "n2",
              "is_active": true,
              "ip_address": "127.0.0.1",
              "port": 7433,
              "path": "/home/ubuntu/pgedge/cluster/test1/n2"
            }
          ]
        }
      ]
    }
}`

	config, err := model.NewClusterConfigFromBytes([]byte(configString))
	require.NoError(t, err)
	require.NotNil(t, config)

	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "test1", config.Name)
		assert.Equal(t, "localhost", config.Style)
		assert.Equal(t, "2024-01-01", config.CreateDate)
	})

	t.Run("databases", func(t *testing.T) {
		require.Len(t, config.Database.Databases, 1)
		assert.Equal(t, "mattermost", config.Database.Databases[0].Name)
		assert.Equal(t, "testusername", config.Database.Databases[0].Username)
		assert.Equal(t, "testpassword", config.Database.Databases[0].Password)

		config.Database.Redact()
		assert.NotEqual(t, "testpassword", config.Database.Databases[0].Password)
	})

	t.Run("nodes", func(t *testing.T) {
		nodes := config.NodeGroups.Nodes()
		require.Len(t, nodes, 2)
		assert.Equal(t, "n1", nodes[0].Name)
		assert.Equal(t, "n2", nodes[1].Name)
	})
}
