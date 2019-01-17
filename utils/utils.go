package utils

func AvgChangePercnt(initialValue, finalValue float64)float64{
	if initialValue == 0{
		return 0
	}
	return  (finalValue-initialValue)/initialValue * 100
}
