package version

// todo use git tagging after release
var version = "v0.2"

var commit string

func GetVersion() string {
	return version
}

func GetCommit() string {
	return commit
}
