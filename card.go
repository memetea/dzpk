package dzpk

//cards operation.
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type CardFace int

const (
	Spade   CardFace = 0
	Hearts  CardFace = 1
	Clubs   CardFace = 2
	Diamond CardFace = 3
)

//牌型
type CardFaceType int

const (
	//RoyalFlush    CardFaceType = 9 //皇家同花顺
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

// card face type percent(based on emulated calc):
// 0 percent: 0.174064
// 1 percent: 0.403919
// 2 percent: 0.269111
// 3 percent: 0.044540
// 4 percent: 0.046323
// 5 percent: 0.030246
// 6 percent: 0.029808
// 7 percent: 0.001676
// 8 percent: 0.000313

// [1, 52].
type Card uint32

func (c Card) Face() CardFace {
	return CardFace((c - 1) / 13)
}
func (c Card) Value() uint32 {
	return uint32(c) - uint32(c.Face()*13)
}
func (c Card) String() string {
	var v string
	switch c.Value() {
	case 1:
		v = "A"
	case 11:
		v = "J"
	case 12:
		v = "Q"
	case 13:
		v = "K"
	default:
		v = strconv.Itoa(int(c.Value()))
	}

	switch c.Face() {
	case Spade:
		return "S" + v
	case Hearts:
		return "H" + v
	case Clubs:
		return "C" + v
	case Diamond:
		return "D" + v
	}

	return ""
}

//represents a collection of cards in 64 bits
//Spade     [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K] at bit [63, 62, 61, 60, 59, 58, 57, 56, 55, 54, 53, 52, 51]
//Hearts    [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K] at bit [47, 46, 45, 44, 43, 42, 41, 40, 39, 38, 37, 36, 35]
//Clubs     [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K] at bit [31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19]
//Diamond   [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K] at bit [15, 14, 13, 12, 11, 10,  9,  8,  7,  6,  5,  4, 3]
type CardCollection uint64

// var bitsMap [52]uint = [52]uint{
// 	63, 62, 61, 60, 59, 58, 57, 56, 55, 54, 53, 52, 51,
// 	47, 46, 45, 44, 43, 42, 41, 40, 39, 38, 37, 36, 35,
// 	31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19,
// 	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3}

var bitsMap [52]uint = [52]uint{
	61, 60, 59, 58, 57, 56, 55, 54, 53, 52, 51, 50, 49,
	45, 44, 43, 42, 41, 40, 39, 38, 37, 36, 35, 34, 33,
	29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17,
	13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
}

func (cc *CardCollection) SetCard(c Card) {
	*cc = (*cc) | (0x01 << bitsMap[int(c-1)])
	if c.Value() == 1 {
		*cc = *cc | (0x01 << uint(3-c.Face()) * 16)
	}
}

func (cc *CardCollection) CardExists(c Card) bool {
	return ((*cc) & (0x01 << bitsMap[int(c-1)])) != 0
}

func NewCardCollection(cc []Card) CardCollection {
	var cardcc CardCollection
	for _, p := range cc {
		cardcc.SetCard(p)
	}
	return cardcc
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func genNCards(n int) []Card {
	//rand.Seed(time.Now().UnixNano())
	var cards []Card
	for len(cards) < n {
		// random card in [0, 52)
		c := Card(rand.Int31n(52) + 1)
		conflict := false
		for _, v := range cards {
			if v == c {
				conflict = true
				break
			}
		}
		if !conflict {
			cards = append(cards, c)
		}
	}
	return cards
}

func parseCard(card string) (Card, error) {
	if len(card) < 2 {
		return Card(0), fmt.Errorf("invalid card representation:%s", card)
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
			return Card(0), fmt.Errorf("invalid card representation:%s", card)
		}
	}

	return Card(int(face)*13 + v), nil
}

type CardsCheck struct {
	face     CardFaceType
	s        uint16    // 用于判断顺子. 如果存在同花顺， 保存同花顺最小的牌.
	s4       [4]uint16 // 用于判断同花顺
	kinds    [14]byte  // 相同大小牌的计数
	flushs   [4]byte   // 相同花色的计数
	topCards []Card    // 最大的5张牌
}

func NewCardsCheck(cc []Card) *CardsCheck {
	ck := &CardsCheck{face: CardFaceNone}

	for _, c := range cc {
		f, v := c.Face(), c.Value()
		ck.flushs[f]++
		ck.kinds[v-1]++
		ck.s4[f] = ck.s4[f] | (0x01 << uint(14-v))
	}
	ck.kinds[13] = ck.kinds[0]
	ck.s = ck.s4[0] | ck.s4[1] | ck.s4[2] | ck.s4[3]

	return ck
}

func NewCardsCheckFromCC(cc CardCollection) *CardsCheck {
	ck := &CardsCheck{face: CardFaceNone}

	ck.s4 = [4]uint16{
		uint16(cc >> 48),
		uint16(cc >> 32),
		uint16(cc >> 16),
		uint16(cc),
	}
	ck.s = ck.s4[0] | ck.s4[1] | ck.s4[2] | ck.s4[3]

	for i, p := range ck.s4 {
		for j := uint(1); j < 14; j++ {
			if (p>>j)&0x01 == 0x01 {
				ck.kinds[13-j] += 1
				ck.flushs[i] += 1
			}
		}
	}

	return ck
}

func (ck *CardsCheck) CardFace() CardFaceType {
	if ck.face != CardFaceNone {
		return ck.face
	}

	var hasStraight bool
	//check straight flush
	n, val := uint16(10), uint16(0)
	for n >= 1 {
		for i := 0; i < len(ck.s4); i++ {
			c := ck.s4[i] >> uint(10-n)
			if c&0x1f == 0x1f && n > val {
				ck.face = StraightFlush
				ck.s = n + uint16(i*13) //ck.s remember the position of base of straight flush
				return ck.face
			}
		}

		//check straight
		if !hasStraight && (ck.s>>uint(10-n))&0x1f == 0x1f {
			hasStraight = true
			ck.s = n
		}

		n--
	}

	//check FourOfAKind
	p31, p32, p21, p22 := 0, 0, 0, 0
	for i := len(ck.kinds) - 1; i >= 0; i-- {
		switch ck.kinds[i] {
		case 4:
			ck.face = FourOfAKind
			ck.s = uint16(i + 1)
			return ck.face
		case 3:
			if p31 == 0 {
				p31 = i + 1
			} else {
				p32 = i + 1
			}
		case 2:
			if p21 == 0 {
				p21 = i + 1
			} else {
				p22 = i + 1
			}
		}
	}

	//check fullhouse
	if p31 > 0 && (p32 > 0 || p21 > 0) {
		ck.face = FullHouse
		//ck.s remember the position of 3-cards in higher 4 bits and 2-cards in lower 4 bits
		if p32 >= 0 {
			ck.s = uint16(p31)<<8 | uint16(p32)
		} else {
			ck.s = uint16(p31)<<8 | uint16(p21)
		}
		return ck.face
	}

	//check flush
	for i, v := range ck.flushs {
		if v >= 5 {
			ck.s = uint16(i) //ck.s remember the face  of flush
			ck.face = Flush
			return ck.face
		}
	}

	//check straight
	if hasStraight {
		ck.face = Straight
		return ck.face
	}

	//check three of a kind
	if p31 > 0 {
		ck.face = ThreeOfAKind
		ck.s = uint16(p31)
		return ck.face
	}

	//check TwoPairs
	if p21 > 0 && p22 > 0 {
		ck.face = TwoPairs
		ck.s = uint16(p21)<<8 | uint16(p22)
		return ck.face
	}

	//check pair
	if p21 > 0 {
		ck.face = Pair
		ck.s = uint16(p21)
		return ck.face
	}

	ck.face = HighCard
	return ck.face
}

func (ck *CardsCheck) Top5Cards() []Card {
	if len(ck.topCards) > 0 {
		return ck.topCards
	}
	if ck.face == CardFaceNone {
		_ = ck.CardFace()
	}
	//return ck.topCards
	if ck.face == StraightFlush {
		if ck.s%13 == 10 {
			ck.topCards = []Card{Card(ck.s), Card(ck.s + 1), Card(ck.s + 2), Card(ck.s + 3), Card(ck.s - 9)}
		} else {
			ck.topCards = []Card{Card(ck.s), Card(ck.s + 1), Card(ck.s + 2), Card(ck.s + 3), Card(ck.s + 4)}
		}
		return ck.topCards
	} else if ck.face == FourOfAKind {
		ck.topCards = []Card{Card(ck.s), Card(ck.s + 13), Card(ck.s + 26), Card(ck.s + 39)}
	} else if ck.face == FullHouse {
		h, l := (ck.s&0xff00)>>8, (ck.s & 0x00ff)
		two := []Card{}
		for i, p := range ck.s4 {
			if p&(0x01<<uint16(12-h)) == 0x01 {
				ck.topCards = append(ck.topCards, Card(int(h)+1+i*13))
			}
			if p&(0x01<<uint16(12-l)) == 0x01 {
				two = append(two, Card(int(l)+1+i*13))
			}
		}
		ck.topCards = append(ck.topCards, two...)
		return ck.topCards
	} else if ck.face == Flush {
		p := ck.s4[ck.s]
		//Ace is the largest
		if (p>>13)&0x01 == 0x01 {
			ck.topCards = append(ck.topCards, Card(ck.s*13+1))
		}
		for i := uint16(1); i <= 12; i++ {
			if (p>>i)&0x01 == 0x01 {
				ck.topCards = append(ck.topCards, Card(ck.s*13+14-i))
			}
			if len(ck.topCards) == 5 {
				break
			}
		}
	} else if ck.face == Straight {
		for i := ck.s; i < ck.s+5; i++ {
			for j, p := range ck.s4 {
				if p>>(14-i)&0x01 == 0x01 {
					ck.topCards = append(ck.topCards, Card(j*13+int(i)))
					break
				}
			}
		}
	} else if ck.face == ThreeOfAKind {
		for i, p := range ck.s4 {
			if p&(0x01<<(12-ck.s)) == 0x01 {
				ck.topCards = append(ck.topCards, Card(int(ck.s)+1+i*13))
			}
		}
	} else if ck.face == TwoPairs {
		h, l := int((ck.s&0xff00)>>8), int(ck.s&0x00ff)
		two := []Card{}
		for i, p := range ck.s4 {
			if p&(0x01<<uint16(len(ck.kinds)-1-h)) == 0x01 {
				ck.topCards = append(ck.topCards, Card(i*13+h+1))
			}
			if p&(0x01<<uint16(len(ck.kinds)-1-l)) == 0x01 {
				two = append(two, Card(i*13+l+1))
			}
		}
		ck.topCards = append(ck.topCards, two...)
	} else if ck.face == Pair {
		for i, p := range ck.s4 {
			if p&(0x01<<uint16(len(ck.kinds)-1-int(ck.s))) == 0x01 {
				ck.topCards = append(ck.topCards, Card(i*13+int(ck.s)+1))
			}
		}
	}

	//pad rest
	if len(ck.topCards) < 5 {
	PAD:
		for i := len(ck.kinds) - 1; i >= 0; i-- {
			if ck.kinds[i] == 0 {
				continue
			}
			for j, p := range ck.s4 {
				if p&(0x01<<uint16(13-i)) == 0x01 {
					c := Card(j*13 + i%13)
					exists := false
					for _, v := range ck.topCards {
						if v == c {
							exists = true
							break
						}
					}
					if !exists {
						ck.topCards = append(ck.topCards, Card(j*13+i%13))
						if len(ck.topCards) == 5 {
							break PAD
						}
					}
				}
			}
		}
	}
	return ck.topCards
}

//returns: big 1, equal 0, small -1
func (ck *CardsCheck) CmpTo(ck2 *CardsCheck) int {
	//compare face first
	if ck.CardFace() != ck2.CardFace() {
		if ck.CardFace() < ck2.CardFace() {
			return 1
		} else {
			return -1
		}
	}

	//both have the same type. compare top cards
	t1 := ck.Top5Cards()
	t2 := ck2.Top5Cards()
	for i := 0; i < 5; i++ {
		if t1[i] > t2[i] {
			return 1
		} else if t1[i] < t2[i] {
			return -1
		}
	}

	return 0
}
