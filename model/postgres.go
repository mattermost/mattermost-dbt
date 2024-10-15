package model

type PgIndexes map[string]*PgIndex

type PgIndex struct {
	SchemaName string
	TableName  string
	IndexName  string
	IndexDef   string
}

func BuildPgIndexesFromSlice(input []*PgIndex) PgIndexes {
	indexes := make(PgIndexes, len(input))
	for _, i := range input {
		indexes[i.IndexName] = i
	}
	return indexes
}

func (pgi *PgIndexes) GetNames() []string {
	names := []string{}
	for n := range *pgi {
		names = append(names, n)
	}
	return names
}

type PgIndexFilter struct {
	SchemaName string
	TableName  string
}
