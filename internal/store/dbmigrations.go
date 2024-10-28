// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package store

import (
	sq "github.com/Masterminds/squirrel"
	mattermostmodel "github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"
)

const (
	dbMigrationsTable = "db_migrations"
)

var dbMigrationsSelect sq.SelectBuilder

func init() {
	dbMigrationsSelect = sq.Select("version", "name").From(dbMigrationsTable)
}

func (sqlStore *SQLStore) GetDBMigrations() ([]*mattermostmodel.AppliedMigration, error) {
	builder := dbMigrationsSelect

	var migrations []*mattermostmodel.AppliedMigration
	err := sqlStore.selectBuilder(sqlStore.db, &migrations, builder)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for db migrations")
	}

	return migrations, nil
}

func (sqlStore *SQLStore) GetDBMigrationsTableCount() (int64, error) {
	return sqlStore.CountRows(dbMigrationsTable)
}
