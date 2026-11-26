package history

import "github.com/Camilo-845/type-cli/internal/game"

const RecentRunsWindow = 10

type Summary struct {
	TotalRuns             int
	PersonalBestWPM       float64
	RecentRuns            int
	RecentAverageWPM      float64
	RecentAverageAccuracy float64
}

func Summarize(results []game.Result) Summary {
	summary := Summary{TotalRuns: len(results)}
	if len(results) == 0 {
		return summary
	}

	recent := results
	if len(recent) > RecentRunsWindow {
		recent = recent[len(recent)-RecentRunsWindow:]
	}
	summary.RecentRuns = len(recent)

	for _, result := range results {
		if result.WPM > summary.PersonalBestWPM {
			summary.PersonalBestWPM = result.WPM
		}
	}

	for _, result := range recent {
		summary.RecentAverageWPM += result.WPM
		summary.RecentAverageAccuracy += result.Accuracy
	}

	summary.RecentAverageWPM /= float64(summary.RecentRuns)
	summary.RecentAverageAccuracy /= float64(summary.RecentRuns)
	return summary
}
