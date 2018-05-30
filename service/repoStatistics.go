package service

// RepoStatistics helds statistics of language usage for repository
type RepoStatistics struct {
	Language   string
	Size       int
	Percentage float64
}
