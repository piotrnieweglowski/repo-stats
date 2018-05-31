package service

import "testing"

func TestGetTotalSumOfBytesArrayIsEmptyReturnsZero(t *testing.T) {
	var statistics = make([]RepoStatistics, 0)
	sum := getTotalSumOfBytes(statistics)

	if sum != 0 {
		t.Error("Sum should be equal 0")
	}
}

func TestGetTotalSumOfBytesReturnsSumOfSize(t *testing.T) {
	var statistics = make([]RepoStatistics, 0)
	statistics = append(statistics, RepoStatistics{Size: 1})
	statistics = append(statistics, RepoStatistics{Size: 2})
	statistics = append(statistics, RepoStatistics{Size: 3})

	expectedSum := 1 + 2 + 3
	sum := getTotalSumOfBytes(statistics)

	if sum != expectedSum {
		t.Errorf("Sum should be equal %d", expectedSum)
	}
}

func TestGropStatisticsByLanguageMergesTheSameLanguages(t *testing.T) {
	var statistics = make([]RepoStatistics, 0)
	statistics = append(statistics, RepoStatistics{Language: "Go", Size: 100})
	statistics = append(statistics, RepoStatistics{Language: "Java", Size: 200})
	statistics = append(statistics, RepoStatistics{Language: "Go", Size: 1000})

	mapOfStatistics := groupStatisticsByLanguage(statistics)

	if len(mapOfStatistics) != 2 {
		t.Error("map should contain 2 items")
	}

	if mapOfStatistics["Go"] != 1100 {
		t.Error("Total size for Go should equal 1100")
	}

	if mapOfStatistics["Java"] != 200 {
		t.Error("Total size for Java should equal 200")
	}
}
