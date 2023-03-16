package output

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	ntsOutputPathKey      = "output.nts_path"
	ntsFileTimeFormat     = "/2006-01-02_"
	ntsDividingLineFormat = "------------ 15:04:05 ------------\n"
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error reading resource file: %s", err)
	}
}

func WriteNTSToFile(raw, parsed, host string) {
	now := time.Now()
	dirPath := viper.GetString(ntsOutputPathKey)
	filePath := dirPath + now.Format(ntsFileTimeFormat) + host + ".txt"

	var file *os.File

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("error creating file %s: %v", filePath, err)
			return
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("error opening file %s: %v", filePath, err)
			return
		}
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("err closing file %s: %v", filePath, err)
		}
	}(file)

	writer := bufio.NewWriter(file)

	dividingLine := now.Format(ntsDividingLineFormat)
	_, err = writer.WriteString(dividingLine)
	if err != nil {
		fmt.Printf("err writing string `%s`: %v", dividingLine, err)
		return
	}

	_, err = writer.WriteString(raw)
	if err != nil {
		fmt.Printf("err writing string `%s`: %v", raw, err)
		return
	}

	_, err = writer.WriteString(beforeParsed)
	if err != nil {
		fmt.Printf("err writing string `%s`: %v", beforeParsed, err)
		return
	}

	_, err = writer.WriteString(parsed)
	if err != nil {
		fmt.Printf("err writing string `%s`: %v", parsed, err)
		return
	}

	err = writer.Flush()
	if err != nil {
		fmt.Printf("error flushing write: %v", err)
		return
	}
}
