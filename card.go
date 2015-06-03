package dzpk

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

//花色
type CardFace int

const (
	Spade   CardFace = 0
	Hearts  CardFace = 1
	Clubs   CardFace = 2
	Diamond CardFace = 3
)

//CardValue ranges from 1 to 13
type CardValue int

func (cv CardValue) String() string {
	switch cv {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		return strconv.Itoa(int(cv))
	}
}

type Card struct {
	Face  CardFace
	Value CardValue
}

//n must be in [1, 52]
func newCard(n int) (*Card, error) {
	if n < 0 || n > 52 {
		return nil, errors.New("Invalid n to newCard")
	}

	v := n % 13
	if v == 0 {
		v = 13
	}
	switch math.Ceil(float64(n) / 13.0) {
	case 4:
		return &Card{Diamond, CardValue(v)}, nil
	case 3:
		return &Card{Clubs, CardValue(v)}, nil
	case 2:
		return &Card{Hearts, CardValue(v)}, nil
	case 1:
		return &Card{Spade, CardValue(v)}, nil
	}

	return nil, fmt.Errorf("someting goning wrong if you see this error: n=%d", n)
}

func parseCard(card string) (*Card, error) {
	if len(card) < 2 {
		return nil, fmt.Errorf("invalid card representation:%s", card)
	}

	var face CardFace
	switch card[0] {
	case 'S':
		face = Spade
	case 'H':
		face = Hearts
	case 'C':
		face = Clubs
	case 'D':
		face = Diamond
	}

	v, err := strconv.Atoi(card[1:])
	if err != nil {
		switch card[1] {
		case 'K':
			v = 13
		case 'Q':
			v = 12
		case 'J':
			v = 11
		case 'A':
			v = 1
		default:
			return nil, fmt.Errorf("invalid card representation:%s", card)
		}
	}

	return &Card{face, CardValue(v)}, nil

}

func newCards(cc ...int) ([]*Card, error) {
	var cards []*Card
	for _, v := range cc {
		c, err := newCard(v)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (c *Card) String() string {
	switch c.Face {
	case Spade:
		return "S" + c.Value.String()
	case Hearts:
		return "H" + c.Value.String()
	case Clubs:
		return "C" + c.Value.String()
	case Diamond:
		return "D" + c.Value.String()
	}

	return ""
}

//gen 52 cards that permutated randomly
func genRandCards() []*Card {
	var cards []*Card
	rand.Seed(time.Now().UnixNano())
	for _, n := range rand.Perm(13 * 4) {
		c, err := newCard(n + 1)
		if err != nil {
			panic(err)
		}
		cards = append(cards, c)
	}
	return cards
}

//牌型
type CardFaceType int

const (
	RoyalFlush    CardFaceType = 9 //皇家同花顺
	StraightFlush CardFaceType = 8 //同花顺
	FourOfAKind   CardFaceType = 7 //四条
	FullHouse     CardFaceType = 6 //葫芦
	Flush         CardFaceType = 5 //同花
	Straight      CardFaceType = 4 //顺子
	ThreeOfAKind  CardFaceType = 3 //三条
	TwoPairs      CardFaceType = 2 //两队
	Pair          CardFaceType = 1 //一对
	HighCard      CardFaceType = 0 //高牌
	CardFaceNone  CardFaceType = -1
)

func (ct CardFaceType) String() string {
	switch ct {
	case RoyalFlush:
		return "RoyalFlush"
	case StraightFlush:
		return "StraightFlush"
	case FourOfAKind:
		return "FourOfAKind"
	case FullHouse:
		return "FullHouse"
	case Flush:
		return "Flush"
	case Straight:
		return "Straight"
	case ThreeOfAKind:
		return "ThreeOfAKind"
	case TwoPairs:
		return "TwoPairs"
	case Pair:
		return "Pair"
	case HighCard:
		return "HighCard"
	default:
		return "-"
	}
}

//sort cards by face and then by their value secondarily
//notice: Ace is the largest one among the same face cards
type SortByFaceAndValue []*Card

func (a SortByFaceAndValue) Len() int      { return len(a) }
func (a SortByFaceAndValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByFaceAndValue) Less(i, j int) bool {
	v1, v2 := a[i].Value, a[j].Value
	if v1 == 1 {
		v1 = 14
	}
	if v2 == 1 {
		v2 = 14
	}
	return int(a[i].Face*13)+int(v1) < int(a[j].Face*13)+int(v2)
}

//sort cards by there value and then by their face secondarily
//notice: Ace is the largest card
type SortByValue []*Card

func (a SortByValue) Len() int      { return len(a) }
func (a SortByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByValue) Less(i, j int) bool {
	if a[i].Value == a[j].Value {
		return a[i].Face < a[j].Face
	} else {
		v1, v2 := a[i].Value, a[j].Value
		if v1 == 1 {
			v1 = 14
		}
		if v2 == 1 {
			v2 = 14
		}
		return v1 < v2
	}
}

type CardCollection struct {
	FaceType CardFaceType
	TopCards [5]*Card
}

func (cc *CardCollection) CmpTo(bb *CardCollection) int {
	if cc.FaceType != bb.FaceType {
		if cc.FaceType > bb.FaceType {
			return 1
		} else {
			return -1
		}
	}
	for i, j := 0, 0; i < len(cc.TopCards) && j < len(bb.TopCards); i, j = i+1, j+1 {
		if cc.TopCards[i].Value != bb.TopCards[j].Value {
			if cc.TopCards[i].Value > bb.TopCards[j].Value {
				return 1
			} else {
				return -1
			}
		}
	}

	return 0
}

//insert cards at index i.
func (cc *CardCollection) insertTopCardsAt(i int, cards ...*Card) {
	for k, j := len(cc.TopCards)-i-len(cards), len(cc.TopCards)-1; k > 0 && j >= 0; k, j = k-1, j-1 {
		cc.TopCards[j] = cc.TopCards[i+k-1]
	}
	for k := 0; k < len(cards); k++ {
		cc.TopCards[i+k] = cards[k]
	}
}

//pick straight from a sorted cards collection
//returns a straight or empty slice
//notice: Ace 2 3 4 5 is also a straight. it's the smallest one.
func getTopStraight(cards []*Card) []*Card {
	var straight []*Card
	for i := 0; i < len(cards)-4; i++ {
		if cards[i+3].Value-cards[i].Value == 3 && cards[i+2].Value-cards[i].Value == 2 && cards[i+1].Value-cards[i].Value == 1 {
			if cards[i+4].Value-cards[i].Value == 4 || cards[i+4].Value-cards[i].Value == (1-10) {
				straight = append([]*Card(nil), cards[i:i+5]...)
			} else if cards[len(cards)-1].Value == 1 && cards[i].Value == 2 {
				straight = append([]*Card(nil), cards[len(cards)-1])
				straight = append(straight, cards[i:i+4]...)
			}
		}
	}

	return straight
}

func SelectTop5(cards []*Card) *CardCollection {
	cc := &CardCollection{
		FaceType: HighCard,
	}

	//check if straight exists
	sort.Sort(SortByValue(cards))
	straight := getTopStraight(cards)
	if len(straight) == 5 {
		cc.FaceType = Straight
		cc.insertTopCardsAt(0, straight...)
	}

	valueCounter := make(map[CardValue]int)
	faceCounter := make(map[CardFace]int)
	for i, v := range cards {
		valueCounter[v.Value] += 1
		faceCounter[v.Face] += 1

		if valueCounter[v.Value] == 4 && cc.FaceType < FourOfAKind {
			cc.FaceType = FourOfAKind
			cc.insertTopCardsAt(3, cards[i])
			break
		} else if valueCounter[v.Value] == 3 {
			if cc.FaceType == TwoPairs {
				cc.FaceType = FullHouse
				cc.insertTopCardsAt(2, cards[i])
			} else if cc.FaceType == Pair {
				cc.FaceType = ThreeOfAKind
				cc.insertTopCardsAt(2, cards[i])
			}
		} else if valueCounter[v.Value] == 2 {
			if cc.FaceType == ThreeOfAKind {
				cc.FaceType = FullHouse
				cc.insertTopCardsAt(3, cards[i-1:i+1]...)
			} else if cc.FaceType == Pair {
				cc.FaceType = TwoPairs
				cc.insertTopCardsAt(0, cards[i-1:i+1]...)
			}
			if cc.FaceType < Pair {
				cc.FaceType = Pair
				cc.insertTopCardsAt(0, cards[i-1:i+1]...)
			}
		}
	}

	//beacuse we sort cards by their value previously.
	//so we need check if there'are flush straight or flush exists next
FLUSHCHECK:
	for k, v := range faceCounter {
		if v >= 5 {
			sort.Sort(SortByFaceAndValue(cards))
			for i := 0; i < len(cards); i++ {
				if cards[i].Face == k {
					straight = getTopStraight(cards[i : i+v])
					if len(straight) == 5 {
						//同花顺
						cc.FaceType = StraightFlush
						cc.insertTopCardsAt(0, straight...)
						if straight[4].Value == 1 {
							cc.FaceType = RoyalFlush
						}
					} else {
						if cc.FaceType < Flush {
							cc.FaceType = Flush
							//同花
							cc.insertTopCardsAt(0, cards[i+v-5:i+v]...)
						}
					}
					break FLUSHCHECK
				}
			}
		}
	}

	//pad nil pos
	for i := 0; i < len(cc.TopCards); i++ {
		if cc.TopCards[i] != nil {
			continue
		}
		for j := len(cards) - 1; j >= 0; j-- {
			included := false
			for k := len(cc.TopCards) - 1; k >= 0; k-- {
				if cc.TopCards[k] == cards[j] {
					included = true
					break
				}
			}
			if !included {
				cc.TopCards[i] = cards[j]
				break
			}
		}
	}

	return cc
}
