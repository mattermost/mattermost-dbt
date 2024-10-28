// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package store

const (
	ChannelsTable = "channels"
)

func (sqlStore *SQLStore) GetChannelsTableCount() (int64, error) {
	return sqlStore.CountRows(ChannelsTable)
}
