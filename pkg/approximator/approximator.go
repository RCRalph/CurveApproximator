package approximator

import (
	"math"
	"sync"
)

func calculateDeviation(data *Data, props *properties, dist *Distribution) {
	var value float64
	var x uint64
	for i := range data.target {
		value, x = 0, dist.partition
		for j := 0; j < len(data.values); j++ {
			value += float64(x%props.base) * data.values[j][i]
			x /= props.base
		}

		dist.deviation += math.Abs(data.target[i]-(value*props.precision)) * data.weights[i]
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
	bestDist.deviation = math.MaxFloat64

	var temp, sum uint64
	for i := start; i < end; i++ {
		temp, sum = i, 0
		for j := 0; j < len(data.values); j++ {
			sum += temp % props.base
			temp /= props.base

			if sum > props.base {
				break
			}
		}

		if sum == props.base {
			dist.partition, dist.deviation = i, 0
			calculateDeviation(data, props, &dist)

			if bestDist.deviation > dist.deviation {
				bestDist.deviation, bestDist.partition = dist.deviation, dist.partition
			}
		}
	}

	m.Lock()
	if result.deviation > bestDist.deviation {
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

	for i := uint8(0); i < data.decimalPlaces; i++ {
		props.base *= 10
		props.precision /= 10
	}

	for i := 0; i < len(data.values); i++ {
		props.totalDists *= props.base
	}

	// Invoke goroutines
	var start, stop uint64
	step, waitGroup, mutex := props.totalDists/uint64(threads), sync.WaitGroup{}, sync.Mutex{}
	result.deviation = math.MaxFloat64

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
