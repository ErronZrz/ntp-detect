package analysis

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func HistogramChart(srcPath, dstDir, prefix string, params histParams) error {
	sharedParams = params

	histMap, err := generateHistMap(srcPath)
	if err != nil {
		return err
	}

	for _, histData := range histMap {
		err := generateHistogramChart(histData, dstDir, prefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateHistogramChart(data histData, dstDir, prefix string) error {
	p := plot.New()
	p.Title.Text = sharedParams.subject + data.name
	p.X.Label.Text = sharedParams.xText

	hist, err := plotter.NewHist(plotter.Values(data.values), 20)
	if err != nil {
		return fmt.Errorf("create histogram error: %v", err)
	}

	hist.LineStyle.Width = vg.Length(0)
	hist.Color = plotutil.Color(0)
	// hist.Normalize(1)

	p.Add(hist)

	chartWidth := 6 * vg.Inch
	chartHeight := 4 * vg.Inch

	err = p.Save(chartWidth, chartHeight, fmt.Sprintf("%s/%s%s.png", dstDir, prefix, p.Title.Text))
	if err != nil {
		return fmt.Errorf("save histogram error: %v", err)
	}

	return nil
}
