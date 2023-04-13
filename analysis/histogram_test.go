package analysis

import "testing"

func TestHistogramBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "domain28_17_"
	partitions := []int{
		1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 60, 70, 80, 90, 100, 120, 150, 200, 300, 500, 1000,
	}

	params := histParams{
		nameCol:    3,
		valCol:     10,
		subject:    "Stratum ",
		xText:      "Root Delay",
		unit:       "ms",
		divisor:    float64(2<<16) / 1000,
		partitions: partitions,
	}
	err := HistogramBarChart(srcPath, dstDir, prefix, params)
	if err != nil {
		t.Error(err)
	}
}
