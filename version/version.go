package version

// todo use git tagging after release
const version = "v0.2"

var commit string

func GetVersion() string {
	return version
}

func GetCommit() string {
	return commit
}
