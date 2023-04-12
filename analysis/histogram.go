package analysis

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
	"os"
	"sort"
	"strconv"
)

type histParams struct {
	nameCol    int
	valCol     int
	title      string
	xText      string
	divisor    float64
	partitions []int
}

type histData struct {
	name   string
	values []float64
}

var (
	sharedParams histParams
)

func HistogramBarChart(srcPath, dstDir, prefix string, params histParams) error {
	sharedParams = params

	histMap, err := generateHistMap(srcPath)
	if err != nil {
		return err
	}

	for _, histData := range histMap {
		err := generateHistogramBarChart(histData, dstDir, prefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateHistMap(srcPath string) (map[string]histData, error) {
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

	res := make(map[string]histData)
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv error: %v", err)
		}
		val, err := strconv.ParseInt(row[sharedParams.valCol], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse int error: %v", err)
		}
		name := row[sharedParams.nameCol]
		data, ok := res[name]
		if !ok {
			data = histData{name: name}
		}
		data.values = append(data.values, float64(val)/sharedParams.divisor)
		res[name] = data
	}
	return res, nil
}

func generateHistogramBarChart(data histData, dstDir, prefix string) error {
	ps := sharedParams.partitions
	n := len(ps)

	values := make(plotter.Values, n+1)
	labels := make([]string, n+1)

	sort.Float64s(data.values)
	idx := 0
	for _, val := range data.values {
		for idx < n && val >= float64(ps[idx]) {
			idx++
		}
		values[idx]++
	}

	for i, p := range ps {
		labels[i] = fmt.Sprintf("<%d", p)
	}
	labels[n] = fmt.Sprintf(">=%d", ps[n-1])

	p := plot.New()
	p.Title.Text = sharedParams.title + data.name
	p.X.Label.Text = sharedParams.xText

	bars, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		return fmt.Errorf("create bar chart error: %v", err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)
	bars.ShowValue = true

	p.Add(bars)
	p.NominalX(labels...)

	chartWidth := (1 + vg.Length(n+1)*0.4) * vg.Inch
	chartHeight := 4 * vg.Inch

	err = p.Save(chartWidth, chartHeight, fmt.Sprintf("%s/%s%s.png", dstDir, prefix, p.Title.Text))
	if err != nil {
		return fmt.Errorf("save bar chart error: %v", err)
	}

	return nil
}
