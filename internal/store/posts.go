package store

const (
	PostsTable = "posts"
)

func (sqlStore *SQLStore) GetPostsTableCount() (int64, error) {
	return sqlStore.CountRows(PostsTable)
}
