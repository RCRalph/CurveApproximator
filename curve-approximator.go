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
	"sync"
)

type CurveData struct {
	targetCurve, weights []float64
	values               [][]float64
}

type CurveProperties struct {
	precision              float64
	valueRange, distNumber uint64
}

type Distribution struct {
	deviation float64
	partition []uint64
}

func getFloat64(value string) (result float64) {
	value = strings.Replace(value, ",", ".", 1)
	result, err := strconv.ParseFloat(value, 64)

	if err != nil {
		panic(err)
	}

	return result
}

func ipow(base uint64, exponent uint64) (result uint64) {
	if exponent == 0 {
		return 1
	}

	result = base
	for i := uint64(2); i <= exponent; i++ {
		result *= base
	}

	return result
}

func setData(filename *string, data *CurveData) {
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		data.targetCurve = append(data.targetCurve, getFloat64(row[0]))
		data.weights = append(data.weights, getFloat64(row[1]))

		if len((data.values)) == 0 {
			data.values = make([][]float64, len(row)-2)
		}

		for i := 0; i < len(row)-2; i++ {
			data.values[i] = append(data.values[i], getFloat64(row[i+2]))
		}
	}

	// Expected data requires summation of data from the last to first rows
	for i := 0; i < len(data.values); {
		var sum float64
		for j := len(data.values[i]) - 1; j >= 0; j-- {
			data.values[i][j], sum = sum, sum+data.values[i][j]
		}

		if math.Round(data.values[i][0]) == 0 {
			data.values[i] = data.values[len(data.values)-1]
			data.values = data.values[:len(data.values)-1]
		} else {
			i++
		}
	}

	fmt.Printf("Successfully gathered %d datasets.\n", len(data.values))
}

func calculateDeviation(dist *Distribution, properties *CurveProperties, data *CurveData) {
	dist.deviation = 0
	var value float64
	for i := range data.targetCurve {
		value = 0
		for j := range dist.partition {
			value += float64(dist.partition[j]) * data.values[j][i]
		}

		dist.deviation += math.Abs(data.targetCurve[i]-(value*properties.precision)) / data.weights[i]
	}
}

func calculateBestDistribution(
	start, end uint64,
	data *CurveData,
	properties *CurveProperties,
	result *Distribution,
	wg *sync.WaitGroup,
	m *sync.Mutex,
) {
	defer wg.Done()

	dist := Distribution{
		deviation: 0,
		partition: make([]uint64, len(data.values)),
	}

	bestDist := Distribution{
		deviation: math.MaxFloat64,
		partition: make([]uint64, len(data.values)),
	}

	var temp, sum uint64
	for i := start; i < end; i++ {
		temp, sum = i, 0
		for j := 0; j < len(data.values); j++ {
			dist.partition[j], temp = temp%properties.valueRange, temp/properties.valueRange

			if sum += dist.partition[j]; sum > properties.valueRange {
				break
			}
		}

		if sum == properties.valueRange {
			calculateDeviation(&dist, properties, data)

			if bestDist.deviation > dist.deviation {
				bestDist.deviation = dist.deviation
				copy(bestDist.partition, dist.partition)
			}
		}
	}

	m.Lock()
	if result.deviation > bestDist.deviation {
		*result = bestDist
	}
	m.Unlock()
}

func main() {
	filename := flag.String("file", "examples.csv", "Filename of the file with datasets")
	pd := flag.Float64("precision", 0.1, "Approximation precision in decimal: 1% = 0.01")
	flag.Parse()

	var data CurveData
	setData(filename, &data)

	properties := CurveProperties{
		precision:  *pd,
		valueRange: uint64(1 / *pd),
		distNumber: ipow(uint64(1 / *pd), uint64(len(data.values))),
	}

	threadNumber := uint64(runtime.GOMAXPROCS(-1))

	step := properties.distNumber / threadNumber
	waitGroup, mutex := sync.WaitGroup{}, sync.Mutex{}
	bestDist := Distribution{
		deviation: math.MaxFloat64,
		partition: make([]uint64, len(data.values)),
	}

	var start, stop uint64
	waitGroup.Add(runtime.GOMAXPROCS(-1))
	for i := uint64(0); i < threadNumber; i++ {
		start = i * step
		if i == threadNumber-1 {
			stop = properties.distNumber
		} else {
			stop = start + step
		}

		go calculateBestDistribution(start, stop, &data, &properties, &bestDist, &waitGroup, &mutex)
	}

	waitGroup.Wait()
	fmt.Print("Result: ")
	for _, item := range bestDist.partition {
		fmt.Printf("%.2f%% ", float64(item)*properties.precision*100)
	}
	fmt.Print("\n")
}