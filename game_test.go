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

func TestSelectCards(t *testing.T) {
	hands := [][]string{
		//[]string{"SA", "SK", "SQ", "S10", "SJ", "H8", "D2"},
		//[]string{"S7", "S9", "SJ", "S8", "S10", "H2", "D5"},
		[]string{"SA", "SJ", "HA", "C3", "CA", "D2", "DA"},
		//[]string{"S3", "C3", "HQ", "CQ", "C9", "SQ", "D8"},
	}

	expectFace := []CardFaceType{
		//RoyalFlush,
		//StraightFlush,
		FourOfAKind,
		//FullHouse,
	}

	expectVal := [][]string{
		//[]string{"S10", "SJ", "SQ", "SK", "SA"},
		//[]string{"S7", "S8", "S9", "S10", "SJ"},
		[]string{"SA", "HA", "CA", "DA", "SJ"},
		//[]string{"SQ", "HQ", "CQ", "S3", "C3"},
	}

	for i, h := range hands {
		var cards []*Card
		for _, c := range h {
			card, err := parseCard(c)
			if err != nil {
				t.Errorf("parseCard err:%v", err)
				t.FailNow()
			}
			cards = append(cards, card)
		}
		collection := SelectTop5(cards)
		if collection.FaceType != expectFace[i] {
			t.Errorf("Expect face of %v is %v, got %v", h, expectFace[i], collection.FaceType)
			t.FailNow()
		}

		for j, v := range collection.TopCards {
			if v.String() != expectVal[i][j] {
				t.Errorf("Expect val of %v is %s. got %v", h, expectVal[i], collection.TopCards)
				t.FailNow()
			}
		}
	}
}
