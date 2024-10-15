// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package store

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-dbt/model"
	"github.com/pkg/errors"
)

const (
	pgStatTable  = "pg_stat_activity"
	pgIndexTable = "pg_indexes"
)

var pgIndexSelect sq.SelectBuilder

func init() {
	pgIndexSelect = sq.Select("schemaname", "tablename", "indexname", "indexdef").From(pgIndexTable)
}

func (sqlStore *SQLStore) GetIndexes(filter *model.PgIndexFilter) (model.PgIndexes, error) {
	builder := pgIndexSelect
	if filter.SchemaName != "" {
		builder = builder.Where("schemaname = ?", filter.SchemaName)
	}
	if filter.TableName != "" {
		builder = builder.Where("tablename = ?", filter.TableName)
	}

	var indexes []*model.PgIndex
	err := sqlStore.selectBuilder(sqlStore.db, &indexes, builder)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for pg indexes")
	}

	return model.BuildPgIndexesFromSlice(indexes), nil
}

func (sqlStore *SQLStore) GetConnectionCount() (int64, error) {
	return sqlStore.CountRows(pgStatTable)
}

func (sqlStore *SQLStore) CreateIndex(index *model.PgIndex) error {
	_, err := sqlStore.exec(sqlStore.db, index.IndexDef)
	if err != nil {
		return errors.Wrap(err, "failed to create pg index")
	}

	return nil
}
