package store

const (
	TeamsTable = "teams"
)

func (sqlStore *SQLStore) GetTeamsTableCount() (int64, error) {
	return sqlStore.CountRows(TeamsTable)
}
