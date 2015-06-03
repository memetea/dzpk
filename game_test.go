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

func TestGetStraight(t *testing.T) {
	hands := [][]string{
		[]string{"S2", "C5", "H4", "S3", "S6", "H7", "SA"},
		[]string{"S2", "C5", "H4", "S3", "S7", "H8", "SA"},
		[]string{"SJ", "CQ", "HK", "S10", "S7", "H8", "SA"},
		[]string{"SJ", "CQ", "HK", "S10", "S9", "H8", "SA"},
	}

	expectVal := [][]string{
		[]string{"S3", "H4", "C5", "S6", "H7"},
		[]string{"SA", "S2", "S3", "H4", "C5"},
		[]string{"S10", "SJ", "CQ", "HK", "SA"},
		[]string{"S10", "SJ", "CQ", "HK", "SA"},
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
		sort.Sort(SortByValue(cards))
		straight := getTopStraight(cards)
		if len(straight) != len(expectVal[i]) {
			t.Errorf("Error return:%v != %v", straight, expectVal[i])
			t.FailNow()
		}

		for j, v := range straight {
			if v.String() != expectVal[i][j] {
				t.Errorf("Expect val of %v is %s. got %v", h, expectVal[i], straight)
				t.FailNow()
			}
		}
	}
}

func TestSelectCards(t *testing.T) {
	hands := [][]string{
		[]string{"SA", "SK", "SQ", "S10", "SJ", "H8", "D2"},
		[]string{"S7", "S9", "SJ", "S8", "S10", "H2", "D5"},
		[]string{"SA", "SJ", "HA", "C3", "CA", "D2", "DA"},
		[]string{"S3", "C3", "HQ", "CQ", "C9", "SQ", "D8"},
		[]string{"HA", "H8", "C3", "H6", "H5", "S4", "H9"},
		[]string{"S2", "C5", "H4", "S3", "S6", "H7", "SA"},
		[]string{"SA", "HA", "CA", "D2", "C3", "H8", "D9"},
		[]string{"SA", "HA", "C3", "D2", "D3", "H8", "DJ"},
		[]string{"S2", "C3", "H2", "D9", "CK", "S8", "SA"},
		[]string{"S2", "S4", "S5", "H7", "C9", "DJ", "SA"},
	}

	expectFace := []CardFaceType{
		RoyalFlush,
		StraightFlush,
		FourOfAKind,
		FullHouse,
		Flush,
		Straight,
		ThreeOfAKind,
		TwoPairs,
		Pair,
		HighCard,
	}

	expectVal := [][]string{
		[]string{"S10", "SJ", "SQ", "SK", "SA"},
		[]string{"S7", "S8", "S9", "S10", "SJ"},
		[]string{"SA", "HA", "CA", "DA", "SJ"},
		[]string{"SQ", "HQ", "CQ", "S3", "C3"},
		[]string{"H5", "H6", "H8", "H9", "HA"},
		[]string{"S3", "H4", "C5", "S6", "H7"},
		[]string{"SA", "HA", "CA", "D9", "H8"},
		[]string{"SA", "HA", "C3", "D3", "DJ"},
		[]string{"S2", "H2", "SA", "CK", "D9"},
		[]string{"SA", "DJ", "C9", "H7", "S5"},
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
			t.Errorf("Expect face of %v is %v, got: %v", h, expectFace[i], collection.FaceType)
			t.FailNow()
		}

		for j, v := range collection.TopCards {
			if v.String() != expectVal[i][j] {
				t.Errorf("Expect top of %v is %s. got: %v", h, expectVal[i], collection.TopCards)
				t.FailNow()
			}
		}
	}
}

func TestTopCardsCmp(t *testing.T) {
	top1 := [][]string{
		[]string{"S10", "SJ", "SQ", "SK", "SA"},
		[]string{"S10", "SJ", "SQ", "SK", "SA"},
		[]string{"S9", "S10", "SJ", "SQ", "SK"},
		[]string{"S9", "S10", "SJ", "SQ", "SK"},
		[]string{"SA", "HA", "CA", "DA", "SJ"},
		[]string{"SA", "HA", "CA", "DA", "SJ"},
	}

	top2 := [][]string{
		[]string{"S9", "S10", "SJ", "SQ", "SK"},
		[]string{"S10", "SJ", "SQ", "SK", "SA"},
		[]string{"S10", "SJ", "SQ", "SK", "SA"},
		[]string{"SA", "HA", "CA", "DA", "SJ"},
		[]string{"SQ", "HQ", "CQ", "S3", "C3"},
		[]string{"SA", "HA", "CA", "DA", "SQ"},
	}

	expect := []int{1, 0, -1, 1, 1, -1}

	for i, h := range top1 {
		var cards []*Card
		for _, c := range h {
			card, err := parseCard(c)
			if err != nil {
				t.Errorf("parseCard err:%v", err)
				t.FailNow()
			}
			cards = append(cards, card)
		}
		collection1 := SelectTop5(cards)

		cards = []*Card(nil)
		for _, c := range top2[i] {
			card, err := parseCard(c)
			if err != nil {
				t.Errorf("parseCard err:%v", err)
				t.FailNow()
			}
			cards = append(cards, card)
		}
		collection2 := SelectTop5(cards)

		cmp := collection1.CmpTo(collection2)
		if cmp != expect[i] {
			t.Errorf("%d: expect %d, got %d.", i, expect[i], cmp)
			t.FailNow()
		}
	}
}
