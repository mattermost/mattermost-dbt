// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package store

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (sqlStore *SQLStore) CountRows(tableName string) (int64, error) {
	type Count struct {
		Count int64
	}
	var counts []Count

	builder := sq.Select("Count (*) as Count").From(tableName)

	err := sqlStore.selectBuilder(sqlStore.db, &counts, builder)
	if err != nil {
		return 0, errors.Wrap(err, "failed to query")
	}

	return counts[0].Count, nil
}
