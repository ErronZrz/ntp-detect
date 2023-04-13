package analysis

import "testing"

func TestHistogramBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "new_syn_"
	partitions := []int{
		1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 60, 70, 80, 90, 100, 120, 150, 200, 300, 500, 1000,
	}

	params := histParams{
		nameCol:    3,
		valCol:     8,
		subject:    "Stratum",
		xText:      "Processing Time",
		unit:       "μs",
		divisor:    float64(2<<32) / 1000000,
		partitions: partitions,
	}
	err := HistogramBarChart(srcPath, dstDir, prefix, params)
	if err != nil {
		t.Error(err)
	}
}
