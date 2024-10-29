package store

const (
	UsersTable = "users"
)

func (sqlStore *SQLStore) GetUsersTableCount() (int64, error) {
	return sqlStore.CountRows(UsersTable)
}
