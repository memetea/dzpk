package dzpk

import (
	"fmt"
	"sort"
	"testing"
)

func TestGenCards(t *testing.T) {
	cards := genRandCards()
	if len(cards) != 52 {
		t.Errorf("len of cards unexpected: %s != 52", len(cards))
		t.FailNow()
	}

	var cardsFaceNum [4]int
	var cardsValNum [13]int

	for i, c := range cards {
		cardsFaceNum[int(c.Face)] += 1
		cardsValNum[int(c.Value)-1] += 1

		fmt.Printf("%s\t", c)
		if (i+1)%13 == 0 && i != 0 {
			fmt.Println()
		}
	}

	//每个花色有13张牌
	for i, v := range cardsFaceNum {
		if v != 13 {
			t.Errorf("Card face %d expect 13, get %d", i, v)
			t.FailNow()
		}
	}
	//每个相同大小的牌， 有4个花色
	for i, v := range cardsValNum {
		if v != 4 {
			t.Errorf("Card num %d expect 4, got %d", i, v)
			t.FailNow()
		}
	}
}

func TestSortCards(t *testing.T) {
	arr, err := newCards(1, 14, 3, 5, 15, 2, 4)
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}
	t.Logf("%v", arr)
	sort.Sort(SortByFaceAndValue(arr))
	t.Logf("%v", arr)
}
