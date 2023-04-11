package analysis

import (
	"active/utils"
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
	"os"
	"strconv"
)

const (
	stratumLimit = 16
)

var (
	nominalX = []string{
		"unsynchronized",
		"1", "2", "3", "4", "5", "6", "7", "8",
		"9", "10", "11", "12", "13", "14", "15",
	}
)

func StratumBarChart(srcPath, dstDir, prefix string) error {
	countryMap, err := generateMap(srcPath)
	if err != nil {
		return err
	}

	countryList := make([]string, 0)
	for country := range countryMap {
		countryList = append(countryList, country)
	}
	engList := utils.TranslateCountry(countryList)

	for i, country := range countryList {
		err := generateBarChart(country, engList[i], countryMap[country], dstDir, prefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateMap(srcPath string) (map[string][]int, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("open file %s error: %v", srcPath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	countryMap := make(map[string][]int)
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv error: %v", err)
		}
		stratum, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse stratum error: %v", err)
		}
		if stratum >= stratumLimit {
			stratum = 0
		}
		bins, ok := countryMap[row[2]]
		if !ok {
			bins = make([]int, stratumLimit)
			bins[stratum] = 1
			countryMap[row[2]] = bins
		} else {
			bins[stratum]++
		}
	}

	return countryMap, nil
}

func generateBarChart(country, eng string, list []int, dir, prefix string) error {
	i, j := 0, stratumLimit
	for list[i] == 0 {
		i++
	}
	for list[j-1] == 0 {
		j--
	}
	diff := j - i
	floatList := make([]float64, diff)
	for k := 0; k < diff; k++ {
		floatList[k] = float64(list[i+k])
	}

	group := plotter.Values(floatList)

	p := plot.New()
	p.Title.Text = "Distribution of Stratum for " + eng
	p.X.Label.Text = "Stratum"
	p.Y.Label.Text = "Count"

	width := vg.Points(20)

	bars, err := plotter.NewBarChart(group, width)
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)

	p.Add(bars)
	p.NominalX(nominalX[i:j]...)

	// TODO: 修正宽度
	chartWidth := vg.Length(diff+3) * vg.Inch

	err = p.Save(chartWidth, 4*vg.Inch, fmt.Sprintf("%s/%s%s.png", dir, prefix, country))
	if err != nil {
		return fmt.Errorf("save bar chart error: %v", err)
	}

	return nil
}
