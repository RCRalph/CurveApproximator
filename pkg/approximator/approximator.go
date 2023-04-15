package approximator

import (
	"math"
	"sync"
)

func calculateDeviation(data *Data, props *properties, dist *Distribution) {
	var value float64
	var x uint64
	for i := range data.Target {
		value, x = 0, dist.Partition
		for j := 0; j < len(data.Values); j++ {
			value += float64(x%props.base) * data.Values[j][i]
			x /= props.base
		}

		dist.Deviation += math.Abs(data.Target[i]-(value*props.precision)) * data.Weights[i]
	}
}

func calculateBestDistribution(
	start, end uint64,
	data *Data,
	props *properties,
	result *Distribution,
	wg *sync.WaitGroup,
	m *sync.Mutex,
) {
	defer wg.Done()

	var dist, bestDist Distribution
	bestDist.Deviation = math.MaxFloat64

	var temp, sum uint64
	for i := start; i < end; i++ {
		temp, sum = i, 0
		for j := 0; j < len(data.Values); j++ {
			sum += temp % props.base
			temp /= props.base

			if sum > props.base {
				break
			}
		}

		if sum == props.base {
			dist.Partition, dist.Deviation = i, 0
			calculateDeviation(data, props, &dist)

			if bestDist.Deviation > dist.Deviation {
				bestDist.Deviation, bestDist.Partition = dist.Deviation, dist.Partition
			}
		}
	}

	m.Lock()
	if result.Deviation > bestDist.Deviation {
		*result = bestDist
	}
	m.Unlock()
}

func Approximate(data *Data, threads int) (result Distribution) {
	// Calculate values for properties
	props := properties{
		base:       1,
		totalDists: 1,
		precision:  1,
	}

	for i := uint8(0); i < data.DecimalPlaces; i++ {
		props.base *= 10
		props.precision /= 10
	}

	for i := 0; i < len(data.Values); i++ {
		props.totalDists *= props.base
	}

	// Invoke goroutines
	var start, stop uint64
	step, waitGroup, mutex := props.totalDists/uint64(threads), sync.WaitGroup{}, sync.Mutex{}
	result.Deviation = math.MaxFloat64

	waitGroup.Add(threads)
	for i := 0; i < threads; i++ {
		start = uint64(i) * step
		if i == threads-1 {
			stop = props.totalDists
		} else {
			stop = start + step
		}

		go calculateBestDistribution(
			start, stop, data, &props,
			&result, &waitGroup, &mutex,
		)
	}

	waitGroup.Wait()

	return result
}
