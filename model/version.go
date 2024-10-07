package model

// BuildHash holds the git commit hash when we build the server.
var BuildHash string

// ShortBuildHash returns a shortened build hash.
func ShortBuildHash() string {
	if len(BuildHash) >= 7 {
		return BuildHash[0:7]
	}

	return BuildHash
}
