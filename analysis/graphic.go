package analysis

import (
	"active/datastruct"
	"encoding/csv"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"io"
	"os"
	"sort"
	"strconv"

	"gonum.org/v1/plot/plotter"
)

const (
	binCount = 10
)

func GenerateGraphic(srcPath, dstDir string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open file %s error: %v", srcPath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	dataList := make([]*datastruct.Statistic, 0)
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read csv error: %v", err)
		}
		stratum, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			return fmt.Errorf("parse stratum error: %v", err)
		}
		dataList = append(dataList, &datastruct.Statistic{
			Domain:  row[0],
			IP:      row[1],
			Country: row[2],
			Stratum: int(stratum),
		})
	}

	countryMap := make(map[string][]float64)
	for _, data := range dataList {
		countryMap[data.Country] = append(countryMap[data.Country], float64(data.Stratum))
	}

	for country, stratumList := range countryMap {
		sort.Float64s(stratumList)
		min, max := stratumList[0], stratumList[len(stratumList)-1]
		dividers := make([]float64, binCount+1)
		dividers[0], dividers[binCount] = min-1, max+1
		diff := (max - min) / float64(binCount)
		for i := 1; i < binCount; i++ {
			dividers[i] = min + diff*float64(i)
		}
		hist := stat.Histogram(nil, dividers, stratumList, nil)
		values := make(plotter.Values, len(hist))
		for i, v := range hist {
			values[i] = v
		}

		p := plot.New()
		p.Title.Text = fmt.Sprintf("Histogram of %s", country)
		p.X.Label.Text = "Stratum"
		p.Y.Label.Text = "Count"
		bars, err := plotter.NewBarChart(values, vg.Points(50))
		if err != nil {
			return fmt.Errorf("create bar chart error: %v", err)
		}
		p.Add(bars)

		err = p.Save(8*vg.Inch, 4*vg.Inch, fmt.Sprintf("%s/%s.png", dstDir, country))
		if err != nil {
			return fmt.Errorf("save plot error: %v", err)
		}
	}

	return nil
}
