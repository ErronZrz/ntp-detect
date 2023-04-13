package analysis

import "testing"

func TestVarianceBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "new_syn_"
	params := varParams{
		valCol:       7,
		yText:        "Offset",
		unit:         "ms",
		divisor:      1000,
		useGlobalAvg: true,
		syncOnly:     true,
	}
	err := VarianceBarChart(srcPath, dstDir, prefix, params)
	if err != nil {
		t.Error(err)
	}
}
