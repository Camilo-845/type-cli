package history

import (
	"testing"

	"github.com/Camilo-845/type-cli/internal/game"
)

func TestSummarizeEmpty(t *testing.T) {
	got := Summarize(nil)
	if got.TotalRuns != 0 || got.RecentRuns != 0 {
		t.Fatalf("expected empty summary, got %+v", got)
	}
	if got.PersonalBestWPM != 0 || got.RecentAverageWPM != 0 || got.RecentAverageAccuracy != 0 {
		t.Fatalf("expected zero values for empty summary, got %+v", got)
	}
}

func TestSummarizeRecentRunsAndPersonalBest(t *testing.T) {
	results := make([]game.Result, 12)
	for i := range results {
		results[i] = game.Result{
			WPM:      float64(i + 1),
			Accuracy: float64(80 + i),
		}
	}

	got := Summarize(results)
	if got.TotalRuns != 12 {
		t.Errorf("TotalRuns = %d, want 12", got.TotalRuns)
	}
	if got.PersonalBestWPM != 12 {
		t.Errorf("PersonalBestWPM = %v, want 12", got.PersonalBestWPM)
	}
	if got.RecentRuns != RecentRunsWindow {
		t.Errorf("RecentRuns = %d, want %d", got.RecentRuns, RecentRunsWindow)
	}
	if got.RecentAverageWPM != 7.5 {
		t.Errorf("RecentAverageWPM = %v, want 7.5", got.RecentAverageWPM)
	}
	if got.RecentAverageAccuracy != 86.5 {
		t.Errorf("RecentAverageAccuracy = %v, want 86.5", got.RecentAverageAccuracy)
	}
}

func TestSummarizeFewerThanRecentWindow(t *testing.T) {
	results := []game.Result{
		{WPM: 40, Accuracy: 95},
		{WPM: 60, Accuracy: 99},
	}

	got := Summarize(results)
	if got.RecentRuns != 2 {
		t.Errorf("RecentRuns = %d, want 2", got.RecentRuns)
	}
	if got.RecentAverageWPM != 50 {
		t.Errorf("RecentAverageWPM = %v, want 50", got.RecentAverageWPM)
	}
	if got.RecentAverageAccuracy != 97 {
		t.Errorf("RecentAverageAccuracy = %v, want 97", got.RecentAverageAccuracy)
	}
}
