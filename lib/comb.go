package lib

func Find(goal float64, list []float64) [][]float64 {
	return addLR(goal, nil, list)
}

func addLR(goal float64, L, R []float64) [][]float64 {
	l := Sum(L)
	var lists [][]float64
	for ri, r := range R {
		list := append([]float64{r}, L...)
		if l+r != goal {
			lists = append(lists, addLR(goal, list, R[ri+1:])...)
		} else if len(list) > 0 {
			lists = append(lists, list)
		}
	}
	return lists
}

func Sum(all []float64) float64 {
	sum := 0.
	for _, x := range all {
		sum += x
	}
	return sum
}
