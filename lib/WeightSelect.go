package lib

//import (
//	"fmt"
//)

type WeightItem struct {
	Index       int
	StartNumber int
	EndNumber   int
	Weight      int
}

func ChooseItem(items []*WeightItem) (selectedIndex int) {
	if 1 >= len(items) {
		return
	}

	totalWeight := 0

	for _, item := range items {
		totalWeight += item.Weight
	}

	randNumber := GetRandomNumber(1, totalWeight)

	for _, item := range items {
		if randNumber >= item.StartNumber && item.EndNumber >= randNumber {
			selectedIndex = item.Index
			break
		}
	}

	return
}
