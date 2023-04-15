package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/RCRalph/CurveApproximator/pkg/approximator"
	"github.com/atotto/clipboard"
)

func getFloat64(value string, separator *string) (result float64) {
	value = strings.Replace(value, *separator, ".", 1)
	result, err := strconv.ParseFloat(value, 64)

	if err != nil {
		panic(err)
	}

	return result
}

func setData(data *approximator.Data, filename, separator *string, delimiter *rune) {
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = *delimiter

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		data.Target = append(data.Target, getFloat64(row[0], separator))
		data.Weights = append(data.Weights, getFloat64(row[1], separator))

		if len((data.Values)) == 0 {
			data.Values = make([][]float64, len(row)-2)
		}

		for i := 0; i < len(row)-2; i++ {
			data.Values[i] = append(data.Values[i], getFloat64(row[i+2], separator))
		}
	}

	// Expected data requires summation of data from the last to first rows
	for i := 0; i < len(data.Values); {
		var sum float64
		for j := len(data.Values[i]) - 1; j >= 0; j-- {
			data.Values[i][j], sum = sum, sum+data.Values[i][j]
		}

		if math.Round(data.Values[i][0]) == 0 {
			data.Values[i] = data.Values[len(data.Values)-1]
			data.Values = data.Values[:len(data.Values)-1]
		} else {
			i++
		}
	}

	fmt.Printf("Successfully gathered %d datasets.\n", len(data.Values))
}

func main() {
	filename := flag.String("file", "examples/cli/Data.csv", "Filename of the file with datasets")
	decimalPlaces := flag.Uint("precision", 1, "Precision of approximation in decimal places, ex. 3 = 0.001")
	separator := flag.String("separator", ".", "Decimal separator")
	d := flag.String("delimiter", ",", "Field delimiter")
	flag.Parse()

	delimiter := rune((*d)[0])
	if *d == "\\t" {
		delimiter = '\t'
	}

	var data approximator.Data
	data.DecimalPlaces = uint8(*decimalPlaces)

	setData(&data, filename, separator, &delimiter)

	bestDist := approximator.Approximate(&data, runtime.GOMAXPROCS(-1))
	bestPartition := bestDist.ToArray(&data)

	printPrecision := int(data.DecimalPlaces) - 2
	if printPrecision < 0 {
		printPrecision = 0
	}

	result := ""
	for i := range bestPartition {
		result += fmt.Sprintf("%.*f%%\t", printPrecision, bestPartition[i]*100)
	}

	if *separator != "." {
		result = strings.ReplaceAll(result, ".", *separator)
	}

	fmt.Println("Result:", result)

	if clipboard.WriteAll(result) == nil {
		fmt.Println("Result copied to clipboard")
	}
}
