package lib

import (
	"strconv"
)

var (
	fmtf = strconv.FormatFloat
)

func Find(goal float64, list []float64, prec int) [][]float64 {
	return addLR(goal, nil, list, prec)
}

func addLR(goal float64, L, R []float64, prec int) [][]float64 {
	Goal := fmtf(goal, 'f', prec, 64)
	l := Sum(L)
	var lists [][]float64
	for ri, r := range R {
		list := append([]float64{r}, L...)
		if fmtf(l+r, 'f', prec, 64) != Goal {
			lists = append(lists, addLR(goal, list, R[ri+1:], prec)...)
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
