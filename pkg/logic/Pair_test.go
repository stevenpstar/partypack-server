package logic

import (
	"fmt"
	"strconv"
	"testing"
)

func TestPair(t *testing.T) {
	//setup data
	var all AllData
	for n := 1; n <= 5; n++ {

		var i []int
		var p = PlayerData{n, i}
		all.Players = append(all.Players, p)

		var pl []int
		var img = ImageData{n, pl}
		all.Images = append(all.Images, img)

	}

	for i := 0; i < 100000; i++ {
		t.Logf("-------------------------\n")
		Pair(all)
		if !allPaired(all.Images) {
			t.Errorf("Not all images / players were matched")
		} else {
			for img := range all.Images {
				//t.Logf("%s : %s, %s", strconv.Itoa(all.Images[img].PlayerId),
				//strconv.Itoa(all.Images[img].Paired[0]),
				//strconv.Itoa(all.Images[img].Paired[1]))
				fmt.Printf(strconv.Itoa(all.Images[img].Paired[0]))

			}
		}
	}

}
