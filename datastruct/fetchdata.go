package datastruct

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func DataFromCSV(path string) ([]*Statistic, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s error: %v", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	res := make([]*Statistic, 0)
	reader := csv.NewReader(file)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read file %s error: %v", path, err)
		}
		next, err := parseRow(row)
		if err != nil {
			return nil, err
		}
		res = append(res, next)
	}

	return res, nil
}

func parseRow(row []string) (*Statistic, error) {
	stratum, err := strconv.Atoi(row[3])
	if err != nil {
		return nil, fmt.Errorf("parse stratum error: %v", err)
	}
	poll, err := strconv.Atoi(row[4])
	if err != nil {
		return nil, fmt.Errorf("parse poll error: %v", err)
	}
	precision, err := strconv.Atoi(row[5])
	if err != nil {
		return nil, fmt.Errorf("parse precision error: %v", err)
	}
	delay, err := strconv.Atoi(row[6])
	if err != nil {
		return nil, fmt.Errorf("parse delay error: %v", err)
	}
	offset, err := strconv.Atoi(row[7])
	if err != nil {
		return nil, fmt.Errorf("parse offset error: %v", err)
	}
	processingTime, err := strconv.Atoi(row[8])
	if err != nil {
		return nil, fmt.Errorf("parse processing time error: %v", err)
	}
	rootDelay, err := strconv.Atoi(row[10])
	if err != nil {
		return nil, fmt.Errorf("parse root delay error: %v", err)
	}
	rootDispersion, err := strconv.Atoi(row[11])
	if err != nil {
		return nil, fmt.Errorf("parse root dispersion error: %v", err)
	}

	return &Statistic{
		Domain:         row[0],
		IP:             row[1],
		Country:        row[2],
		Stratum:        stratum,
		Poll:           poll,
		Precision:      precision,
		Delay:          delay,
		Offset:         offset,
		ProcessingTime: processingTime,
		RefCountry:     row[9],
		RootDelay:      rootDelay,
		RootDisp:       rootDispersion,
	}, nil
}
