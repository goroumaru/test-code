package FTX

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func analysis() {
	data := readData()
	fmt.Println(data)
}

func readData() (data [][]string) {
	data = make([][]string, 0)

	csvFile, err := os.Open("./BTC_USD BitMEX 過去データ.csv")
	defer csvFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(csvFile)
	var line []string
	fmt.Println("--- START ---")
	for {
		line, err = reader.Read()
		if err == io.EOF {
			fmt.Println("--- END ---")
			break
		} else if err != nil {
			fmt.Println("--- ERROR END ---")
			fmt.Println(err)
			break
		}
		fmt.Println(line)
		data = append(data, line)
	}
	return data
}
