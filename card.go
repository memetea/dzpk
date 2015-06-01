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

func newCard(n int) (*Card, error) {
	//n must be in [1, 52]
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

func genRandCards() []*Card {
	var cards []*Card
	//gen 52 cards randomly. cards span [0, 52]
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

type SortByFaceAndValue []*Card

func (a SortByFaceAndValue) Len() int      { return len(a) }
func (a SortByFaceAndValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByFaceAndValue) Less(i, j int) bool {
	return int(a[i].Face*13)+int(a[i].Value) < int(a[j].Face*13)+int(a[j].Value)
}

type SortByValue []*Card

func (a SortByValue) Len() int      { return len(a) }
func (a SortByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByValue) Less(i, j int) bool {
	return a[i].Value < a[j].Value
}

type CardCollection struct {
	FaceType  CardFaceType
	Top5Cards []*Card
}

func getTopStraight(cards []*Card) []*Card {
	var straight []*Card
	var faceType = CardFaceNone
	var straightTop CardValue
	for i := 4; i < len(cards); i++ {
		if cards[i].Value-cards[i-4].Value == 4 {
			faceType = Straight
			straightTop = cards[i].Value
		} else if cards[i].Value == 13 && cards[i-3].Value == (13-10) && cards[0].Value == 1 {
			//10, J, Q, K, A
			faceType = Straight
			straightTop = 1
		}

		if faceType == Straight {
			if straightTop == 1 || straight[4].Value < straightTop {
				if straightTop == 1 {
					copy(straight, cards[i-3:i+1])
					straight = append(straight, cards[0])
				} else {
					copy(straight, cards[i-4:i+1])
				}
			}
		}
	}
	return straight
}

func insertBefore(i []*Card, cc *[]*Card) {
	*cc = append(append([]*Card{}, i...), (*cc)[:5-len(*cc)]...)
}

func SelectTop5(cards []*Card) *CardCollection {
	cc := &CardCollection{FaceType: HighCard}

	//检查是否有顺子
	sort.Sort(SortByValue(cards))
	straight := getTopStraight(cards)
	if len(straight) == 5 {
		cc.FaceType = Straight
		copy(cc.Top5Cards, straight)
	}

	valueCounter := make(map[CardValue]int)
	faceCounter := make(map[CardFace]int)
	for i, v := range cards {
		valueCounter[v.Value] += 1
		faceCounter[v.Face] += 1

		if valueCounter[v.Value] == 4 && cc.FaceType < FourOfAKind {
			cc.FaceType = FourOfAKind
			insertBefore(cards[i-3:i+1], &cc.Top5Cards)
			break
		} else if valueCounter[v.Value] == 3 {
			if cc.FaceType == Pair && cc.FaceType < FullHouse {
				cc.FaceType = FullHouse
				insertBefore(cards[i-2:i+1], &cc.Top5Cards)
			}
			if cc.FaceType < ThreeOfAKind {
				cc.FaceType = ThreeOfAKind
				insertBefore(cards[i-2:i+1], &cc.Top5Cards)
			}
		} else if valueCounter[v.Value] == 2 {
			if cc.FaceType == ThreeOfAKind && cc.FaceType < FullHouse {
				cc.FaceType = FullHouse
				cc.Top5Cards = append(cc.Top5Cards[0:3], cards[i-1:i+1]...)
			}
			if cc.FaceType == Pair && cc.FaceType < TwoPairs {
				cc.FaceType = TwoPairs
				insertBefore(cards[i-1:i+1], &cc.Top5Cards)
			}
			if cc.FaceType < Pair {
				insertBefore(cards[i-1:i+1], &cc.Top5Cards)
			}
		} else if cc.FaceType == HighCard {
			insertBefore(cards[i:i+1], &cc.Top5Cards)
		}
	}

	//检查是否有同花
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
						copy(cc.Top5Cards, straight)
						if straight[4].Value == 1 {
							cc.FaceType = RoyalFlush
						}
					} else {
						if cc.FaceType < Flush {
							cc.FaceType = Flush
							//同花
							copy(cc.Top5Cards, cards[i+v-5:i+v])
						}
					}
					break FLUSHCHECK
				}
			}
		}
	}

	return cc
}

func (cc *CardCollection) CmpTo(bb *CardCollection) int {

	return 0
}
