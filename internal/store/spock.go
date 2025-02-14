// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package store

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

var (
	spockVersionSelect     sq.SelectBuilder
	spockVersionNumSelect  sq.SelectBuilder
	spockReplicationSelect sq.SelectBuilder
	spockLagTrackerSelect  sq.SelectBuilder
)

const (
	spockVersionFunction     = "spock.spock_version()"
	spockVersionNumFunction  = "spock.spock_version_num()"
	spockReplicationFunction = "spock.sub_show_status()"
	spockLagTrackerTable     = "spock.lag_tracker"
)

type SpockVersion struct {
	Version string `json:"spock_version" db:"spock_version"`
}

type SpockVersionNum struct {
	VersionNum string `json:"spock_version_num" db:"spock_version_num"`
}

type SpockStatus struct {
	SubscriptionName string `json:"subscription_name" db:"subscription_name"`
	Status           string `json:"status" db:"status"`
	ReplicationSets  string `json:"replication_sets" db:"replication_sets"`
}

type SpockLagTracker struct {
	SlotName            string `json:"slot_name" db:"slot_name"`
	CommitTimestamp     string `json:"commit_timestamp" db:"commit_timestamp"`
	ReplicationLag      string `json:"replication_lag" db:"replication_lag"`
	ReplicationLagBytes string `json:"replication_lag_bytes" db:"replication_lag_bytes"`
}

func init() {
	spockVersionSelect = sq.Select("spock_version").From(spockVersionFunction)
	spockVersionNumSelect = sq.Select("spock_version_num").From(spockVersionNumFunction)
	spockReplicationSelect = sq.Select("subscription_name", "status", "replication_sets").From(spockReplicationFunction)
	spockLagTrackerSelect = sq.Select("slot_name", "commit_timestamp", "replication_lag", "replication_lag_bytes").From(spockLagTrackerTable)
}

func (sqlStore *SQLStore) GetSpockVersion() (*SpockVersion, error) {
	builder := spockVersionSelect

	var version []*SpockVersion
	err := sqlStore.selectBuilder(sqlStore.db, &version, builder)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for spock version")
	}

	return version[0], nil
}

func (sqlStore *SQLStore) GetSpockVersionNum() (*SpockVersionNum, error) {
	builder := spockVersionNumSelect

	var versionNum []*SpockVersionNum
	err := sqlStore.selectBuilder(sqlStore.db, &versionNum, builder)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for spock version num")
	}

	return versionNum[0], nil
}

func (sqlStore *SQLStore) GetSpockReplicationStatus() (*SpockStatus, error) {
	builder := spockReplicationSelect

	var status []*SpockStatus
	err := sqlStore.selectBuilder(sqlStore.db, &status, builder)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for spock replication status")
	}

	return status[0], nil
}

func (sqlStore *SQLStore) GetSpockLag() (*SpockLagTracker, error) {
	builder := spockLagTrackerSelect

	var lag []*SpockLagTracker
	err := sqlStore.selectBuilder(sqlStore.db, &lag, builder)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for spock replication lag")
	}

	return lag[0], nil
}
