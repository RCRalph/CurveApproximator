package approximator

type Data struct {
	decimalPlaces   uint8
	target, weights []float64
	values          [][]float64
}

type Distribution struct {
	partition uint64
	deviation float64
}

type properties struct {
	base, totalDists uint64
	precision        float64
}

func (dist Distribution) toArray(data *Data) (result []float64) {
	base, precision := uint64(1), float64(1)
	for i := uint8(0); i < data.decimalPlaces; i++ {
		base *= 10
		precision /= 10
	}

	partition := dist.partition
	for i := 0; i < len(data.values); i++ {
		result[i] = float64(partition%base) * precision
		partition /= base
	}

	return result
}
