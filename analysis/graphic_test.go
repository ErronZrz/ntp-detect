package analysis

import "testing"

func TestGenerateGraphic(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain21_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	err := GenerateGraphic(srcPath, dstDir)
	if err != nil {
		t.Error(err)
	}
}
