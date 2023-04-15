package approximator

type Data struct {
	DecimalPlaces   uint8
	Target, Weights []float64
	Values          [][]float64
}

type Distribution struct {
	Partition uint64
	Deviation float64
}

type properties struct {
	base, totalDists uint64
	precision        float64
}

func (dist Distribution) ToArray(data *Data) (result []float64) {
	base, precision := uint64(1), float64(1)
	for i := uint8(0); i < data.DecimalPlaces; i++ {
		base *= 10
		precision /= 10
	}

	partition := dist.Partition
	for i := 0; i < len(data.Values); i++ {
		result = append(result, float64(partition%base)*precision)
		partition /= base
	}

	return result
}
