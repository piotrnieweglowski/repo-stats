package service

// RepoStatistics helds statistics of language usage for repository
type RepoStatistics struct {
	Language   string
	Size       int
	Percentage float64
}

func getTotalSumOfBytes(statistics []RepoStatistics) int {
	sum := 0
	for _, v := range statistics {
		sum += v.Size
	}

	return sum
}

func groupStatisticsByLanguage(statistics []RepoStatistics) map[string]int {
	grouped := make(map[string]int)
	for _, val := range statistics {
		if elem, ok := grouped[val.Language]; ok {
			grouped[val.Language] = elem + val.Size
		} else {
			grouped[val.Language] = val.Size
		}
	}

	return grouped
}
