package dzpk

import (
	"strconv"
	"testing"
)

func TestGenCards(t *testing.T) {
	cards := genNCards(7)
	t.Logf("len:%d, cards:%v", len(cards), cards)
	var cc CardCollection
	for _, c := range cards {
		cc.SetCard(c)
	}
	t.Logf("%s", strconv.FormatUint(uint64(cc), 2))
	for _, c := range cards {
		if !cc.CardExists(c) {
			t.Logf("card %d expect to exist in cardcollection", c)
			t.FailNow()
		}
	}
}

func TestSelectTop5(t *testing.T) {
	cards := [][]string{
		[]string{"S8", "D7", "SQ", "S9", "SJ", "S10", "S3"},
		[]string{"S2", "S4", "H4", "C4", "D4", "S9", "C3"},
		[]string{"H5", "D5", "C2", "C9", "D2", "S5", "SA"},
		[]string{"S1", "S5", "C9", "S7", "D2", "S6", "S3"},
		[]string{"S2", "H3", "C4", "D5", "S6", "H7", "C8"},
		[]string{"S4", "H4", "C4", "D8", "S9", "C10", "SJ"},
		[]string{"S8", "H8", "S9", "C9", "D2", "D6", "C3"},
		[]string{"S3", "D3", "S5", "D6", "C8", "H9", "SJ"},
		[]string{"SQ", "D10", "C8", "H7", "C5", "S3", "D2"},
	}
	expectFace := []CardFaceType{
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
	expectCards := [][]string{
		[]string{"S8", "S9", "S10", "SJ", "SQ"},
		[]string{"S4", "H4", "C4", "D4", "S9"},
		[]string{"S5", "H5", "D5", "C2", "D2"},
		[]string{"SA", "S7", "S6", "S5", "S3"},
		[]string{"C4", "D5", "S6", "H7", "C8"},
		[]string{"S4", "H4", "C4", "SJ", "C10", "S9"},
		[]string{"S9", "C9", "S8", "H8", "D6"},
		[]string{"S3", "D3", "SJ", "H9", "C8"},
		[]string{"SQ", "D10", "C8", "H7", "C5"},
	}

	for i := 0; i < len(cards); i++ {
		var cc []Card
		for j := 0; j < len(cards[i]); j++ {
			c, err := parseCard(cards[i][j])
			if err != nil {
				t.Logf("%v", err)
				t.FailNow()
			}
			cc = append(cc, c)
		}
		ck := NewCardsCheck(cc)
		if ck.CardFace() != expectFace[i] {
			t.Logf("Expect face:%v. got:%v", expectFace[i], ck.CardFace())
			t.FailNow()
		}
		for j, v := range ck.Top5Cards() {
			if expectCards[i][j] != v.String() {
				t.Logf("top5 expect:%v, got:%v", expectCards[i], ck.Top5Cards())
				t.FailNow()
			}
		}
	}
}
