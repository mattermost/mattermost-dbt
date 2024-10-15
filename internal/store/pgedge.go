// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package store

import (
	"fmt"

	"github.com/mattermost/mattermost-dbt/model"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type PgedgeNodeStore struct {
	Node  *model.PgedgeNode
	Store *SQLStore
}

func NewStoreForPgedgeNode(db *model.PgedgeDatabase, node *model.PgedgeNode, logger log.FieldLogger) (*SQLStore, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?connect_timeout=10&sslmode=disable",
		db.Username, db.Password, node.IPAddress, node.Port, db.Name,
	)

	store, err := New(dsn, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new store")
	}

	_, err = store.GetConnectionCount()
	if err != nil {
		return nil, errors.Wrap(err, "database connection health check failed")
	}

	return store, nil
}

func NewStoresForAllPgedgeNodes(pgedgeConfig *model.PgedgeClusterConfig, logger log.FieldLogger) ([]*PgedgeNodeStore, error) {
	nodes := pgedgeConfig.NodeGroups.Nodes()
	if len(nodes) == 0 {
		return nil, errors.New("config contains no database nodes")
	}

	var nodesStores []*PgedgeNodeStore
	for _, node := range nodes {
		store, err := NewStoreForPgedgeNode(pgedgeConfig.Database.Databases[0], &node, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create store for pgEdge node %s", node.Name)
		}

		nodesStores = append(nodesStores, &PgedgeNodeStore{
			&node,
			store,
		})
	}

	return nodesStores, nil
}
