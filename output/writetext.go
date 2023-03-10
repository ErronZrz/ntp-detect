package output

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	configPath         = "../resource/"
	outputPathKey      = "output.dir_path"
	fileTimeFormat     = "/2006-01-02_15-04-05_"
	dividingLineFormat = "------------ 15:04:05.000 ------------\n"
	beforeParsed       = "--- parsed ---\n"
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

func WriteToFile(raw, parsed, cmd string, seq int, rcvTime, now time.Time) {
	dirPath := viper.GetString(outputPathKey)
	cmd = strings.Replace(cmd, "/", "_", 1)

	filePath := dirPath + now.Format(fileTimeFormat) + cmd + ".txt"

	var file *os.File

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("error creating file %s: %s", filePath, err)
			return
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("error opening file %s: %s", filePath, err)
			return
		}
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("error closing file %s: %s", filePath, err)
		}
	}(file)

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString("#" + strconv.Itoa(seq) + "\n")
	if err != nil {
		fmt.Printf("err writing sequence `%d`: %s", seq, err)
		return
	}

	dividingLine := rcvTime.Format(dividingLineFormat)
	_, err = writer.WriteString(dividingLine)
	if err != nil {
		fmt.Printf("error writing string `%s`: %s", dividingLine, err)
		return
	}

	_, err = writer.WriteString(raw)
	if err != nil {
		fmt.Printf("error writing string `%s`: %s", raw, err)
		return
	}

	_, err = writer.WriteString(beforeParsed)
	if err != nil {
		fmt.Printf("error writing string `%s`: %s", beforeParsed, err)
		return
	}

	_, err = writer.WriteString(parsed)
	if err != nil {
		fmt.Printf("error writing string `%s`: %s", parsed, err)
		return
	}

	err = writer.Flush()
	if err != nil {
		fmt.Printf("error flushing writer: %s", err)
		return
	}
}
