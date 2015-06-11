package dzpk

//cards operation.
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var cardsShift [52]uint = [52]uint{
	61, 60, 59, 58, 57, 56, 55, 54, 53, 52, 51, 50, 49,
	45, 44, 43, 42, 41, 40, 39, 38, 37, 36, 35, 34, 33,
	29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17,
	13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
}

//mask for card [1, 13]. index from 0
var cardsMask [13]uint16 = [13]uint16{
	0x01 << 13,
	0x01 << 12,
	0x01 << 11,
	0x01 << 10,
	0x01 << 9,
	0x01 << 8,
	0x01 << 7,
	0x01 << 6,
	0x01 << 5,
	0x01 << 4,
	0x01 << 3,
	0x01 << 2,
	0x01 << 1,
}

//card index
const (
	Card_K    = 12
	Card_Q    = 11
	Card_J    = 10
	Card_10   = 9
	Card_9    = 8
	Card_8    = 7
	Card_7    = 6
	Card_6    = 5
	Card_5    = 4
	Card_4    = 3
	Card_3    = 2
	Card_2    = 1
	Card_A    = 0
	Card_None = -1
)

type CardFace int

const (
	Spade   CardFace = 0
	Hearts  CardFace = 1
	Clubs   CardFace = 2
	Diamond CardFace = 3
)

//牌型
type CardFaceType int8

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

func (ct CardFaceType) String() string {
	switch ct {
	// case RoyalFlush:
	// 	return "RoyalFlush"
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

// [0, 51].
type Card uint32

func (c Card) Face() CardFace {
	return CardFace(c / 13)
}
func (c Card) Value() uint32 {
	return uint32(c % 13)
}
func (c Card) String() string {
	var v string
	switch c.Value() {
	case Card_A:
		v = "A"
	case Card_J:
		v = "J"
	case Card_Q:
		v = "Q"
	case Card_K:
		v = "K"
	default:
		v = strconv.Itoa(int(c.Value() + 1))
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

type SortByValue []Card

func (a SortByValue) Len() int           { return len(a) }
func (a SortByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByValue) Less(i, j int) bool { return a[i].Value() < a[j].Value() }

func init() {
	rand.Seed(time.Now().UnixNano())
}

func genNCards(n int) []Card {
	//rand.Seed(time.Now().UnixNano())
	var cards []Card
	for len(cards) < n {
		// random card in [0, 52)
		c := Card(rand.Int31n(52))
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
		return Card(0), fmt.Errorf("Invalid card:%v", card)
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
	default:
		return Card(0), fmt.Errorf("Invalid card:%v", card)
	}

	v, err := strconv.Atoi(card[1:])
	if err != nil {
		switch card[1] {
		case 'K':
			v = Card_K
		case 'Q':
			v = Card_Q
		case 'J':
			v = Card_J
		case 'A':
			v = Card_A
		default:
			return Card(0), fmt.Errorf("Invalid card:%v", card)
		}
	} else {
		v = v - 1
	}

	return Card(int(face)*13 + v), nil
}

//represents a collection of cards in 64 bits
//Spade     [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K, A] at bit [61, 60, 59, 58, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48]
//Hearts    [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K, A] at bit [45, 44, 43, 42, 41, 40, 39, 38, 37, 36, 35, 34, 33, 32]
//Clubs     [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K, A] at bit [29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16]
//Diamond   [A, 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K, A] at bit [13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0]
type CardCollection uint64

func (cc *CardCollection) SetCard(c Card) {
	*cc = (*cc) | (0x01 << cardsShift[int(c)])
	if c.Value() == Card_A {
		*cc = *cc | (0x01 << uint(3-c.Face()) * 16)
	}
}

func (cc *CardCollection) CardExists(c Card) bool {
	return ((*cc) & (0x01 << cardsShift[int(c)])) != 0
}

func NewCardCollection(cc []Card) CardCollection {
	var cardcc CardCollection
	for _, p := range cc {
		cardcc.SetCard(p)
	}
	return cardcc
}

type CardsCheck struct {
	face     CardFaceType
	s        uint16    // 用于判断顺子. 如果存在同花顺， 保存同花顺最小的牌.
	s4       [4]uint16 // 用于判断同花顺
	kinds    [15]byte  // 相同大小牌的计数. index 0-12 放牌出现的张数; index 13 复制一份Card_A的张数; index 14 放最大牌位置
	flushs   [4]byte   // 相同花色的计数
	topCards []Card    // 最大的5张牌
}

func NewCardsCheck(cc []Card) *CardsCheck {
	ck := &CardsCheck{face: CardFaceNone}
	for _, c := range cc {
		f, v := c.Face(), c.Value()
		if byte(v) > ck.kinds[14] {
			if byte(v) == Card_A {
				ck.kinds[14] = Card_K + 1
			} else {
				ck.kinds[14] = byte(v)
			}
		}
		ck.flushs[f]++
		ck.kinds[v]++
		ck.s4[f] = ck.s4[f] | (0x01 << uint(13-v))
	}
	for i := 0; i < len(ck.s4); i++ {
		ck.s4[i] = ck.s4[i] | (ck.s4[i] >> 13)
	}
	//copy ace count
	ck.kinds[13] = ck.kinds[Card_A]
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
		for j := Card_A; j < Card_K; j++ {
			if p&cardsMask[j] == cardsMask[j] {
				ck.kinds[j] += 1
				ck.flushs[i] += 1
				if byte(j) > ck.kinds[14] {
					if byte(j) == Card_A {
						ck.kinds[14] = Card_K + 1
					} else {
						ck.kinds[14] = byte(j)
					}
				}
			}
		}
	}
	//copy ace count
	ck.kinds[13] = ck.kinds[Card_A]
	return ck
}

func (ck *CardsCheck) CardCollection() CardCollection {
	var cc CardCollection
	cc = cc | CardCollection(ck.s4[0])<<48 | CardCollection(ck.s4[1])<<32 | CardCollection(ck.s4[2])<<16 | CardCollection(ck.s4[3])
	return cc
}

func (ck *CardsCheck) CardFace() CardFaceType {
	if ck.face != CardFaceNone {
		return ck.face
	}
	//return ck.face
	//check straight flush
	n, val := Card_10, 0
	for n >= Card_A {
		for i := 0; i < len(ck.s4); i++ {
			//_ = val
			if (ck.s4[i]>>uint(9-n))&0x1f == 0x1f && n >= val {
				ck.face = StraightFlush
				ck.s = uint16(i*13 + n) //ck.s remember the position of base of straight flush
				return ck.face
			}
		}

		n--
	}
	//return ck.face

	//check FourOfAKind
	p3, p21, p22 := -1, -1, -1
	for i := 13; i >= 0; i-- {
		switch ck.kinds[i] {
		case 4:
			ck.face = FourOfAKind
			ck.s = uint16(i % 13)
			return ck.face
		case 3:
			if p3 == -1 {
				p3 = i % 13
			} else {
				ck.s = uint16(p3)<<8 | uint16(i)
				ck.face = FullHouse
				return ck.face
			}
		case 2:
			if p21 == -1 {
				p21 = i % 13
			} else {
				p22 = i
			}
		}
	}

	//check fullhouse
	if p3 >= 0 && p21 >= 0 {
		ck.face = FullHouse
		//ck.s remember the position of 3-cards in higher 4 bits and 2-cards in lower 4 bits
		ck.s = uint16(p3)<<8 | uint16(p21)
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
	n = Card_10
	for n >= Card_A {
		//check straight
		if (ck.s>>uint(9-n))&0x1f == 0x1f {
			ck.face = Straight
			ck.s = uint16(n)
			return ck.face
		}
		n--
	}

	//check three of a kind
	if p3 >= 0 {
		ck.face = ThreeOfAKind
		ck.s = uint16(p3)
		return ck.face
	}

	//check TwoPairs
	if p21 >= 0 && p22 >= 0 {
		ck.face = TwoPairs
		ck.s = uint16(p21)<<8 | uint16(p22)
		return ck.face
	}

	//check pair
	if p21 >= 0 {
		ck.face = Pair
		ck.s = uint16(p21)
		return ck.face
	}

	ck.face = HighCard
	return ck.face
}

//returns the top cards(high cards not included.)
//for example: there're 4 cards returned when CardFaceType is FourOfAKind
func (ck *CardsCheck) TopCards() []Card {
	if len(ck.topCards) > 0 {
		return ck.topCards
	}
	if ck.face == CardFaceNone {
		_ = ck.CardFace()
	}
	//return ck.topCards
	if ck.face == StraightFlush {
		if ck.s%13 == Card_10 {
			ck.topCards = []Card{Card(ck.s), Card(ck.s + 1), Card(ck.s + 2), Card(ck.s + 3), Card(ck.s - 9)}
		} else {
			ck.topCards = []Card{Card(ck.s), Card(ck.s + 1), Card(ck.s + 2), Card(ck.s + 3), Card(ck.s + 4)}
		}
		return ck.topCards
	} else if ck.face == FourOfAKind {
		ck.topCards = []Card{Card(ck.s), Card(ck.s + 13), Card(ck.s + 26), Card(ck.s + 39)}
		ck.kinds[ck.s], ck.kinds[13] = 0, ck.kinds[0]
		return ck.topCards[:4]
	} else if ck.face == FullHouse {
		h, l := (ck.s&0xff00)>>8, (ck.s & 0x00ff)
		hMask, lMask := cardsMask[h%13], cardsMask[l%13]
		two := []Card{}
		for i, p := range ck.s4 {
			if p&hMask == hMask {
				ck.topCards = append(ck.topCards, Card(i*13+int(h)))
			}
			if p&lMask == lMask {
				two = append(two, Card(i*13+int(l)))
			}
		}
		ck.topCards = append(ck.topCards, two...)
		return ck.topCards
	} else if ck.face == Flush {
		p := ck.s4[ck.s]
		//Ace is the largest
		if p&0x01 == 0x01 {
			ck.topCards = append(ck.topCards, Card(ck.s*13))
		}
		for i := Card_K; i >= Card_2; i-- {
			if p&cardsMask[i] == cardsMask[i] {
				ck.topCards = append(ck.topCards, Card(int(ck.s*13)+i))
			}
			if len(ck.topCards) == 5 {
				return ck.topCards
			}
		}
	} else if ck.face == Straight {
		for i := ck.s; i < ck.s+5; i++ {
			for j, p := range ck.s4 {
				if p&cardsMask[i%13] == cardsMask[i%13] {
					ck.topCards = append(ck.topCards, Card(j*13+int(i)))
					break
				}
			}
			if len(ck.topCards) == 5 {
				return ck.topCards
			}
		}
	} else if ck.face == ThreeOfAKind {
		mask := cardsMask[ck.s]
		for i, p := range ck.s4 {
			if p&mask == mask {
				ck.topCards = append(ck.topCards, Card(i*13+int(ck.s)))
			}
		}
		ck.kinds[ck.s], ck.kinds[13] = 0, ck.kinds[0]
		return ck.topCards[:3]
	} else if ck.face == TwoPairs {
		h, l := int((ck.s&0xff00)>>8), int(ck.s&0x00ff)
		hMask, lMask := cardsMask[h], cardsMask[l]
		two := []Card{}
		for i, p := range ck.s4 {
			if p&hMask == hMask {
				ck.topCards = append(ck.topCards, Card(i*13+h))
			}
			if p&lMask == lMask {
				two = append(two, Card(i*13+l))
			}
		}
		ck.topCards = append(ck.topCards, two...)
		ck.kinds[h], ck.kinds[l], ck.kinds[13] = 0, 0, ck.kinds[0]
		return ck.topCards[:4]
	} else if ck.face == Pair {
		mask := cardsMask[ck.s]
		for i, p := range ck.s4 {
			if p&mask == mask {
				ck.topCards = append(ck.topCards, Card(i*13+int(ck.s)))
			}
		}
		ck.kinds[ck.s], ck.kinds[13] = 0, ck.kinds[0]
		return ck.topCards[:2]
	}

	return ck.topCards
}

func (ck *CardsCheck) Top5Cards() []Card {
	if len(ck.topCards) == 0 {
		_ = ck.TopCards()
	}
	//pad rest
	if len(ck.topCards) < 5 {
		for i := ck.kinds[14]; i > 0; i-- {
			if ck.kinds[i] == 0 {
				continue
			}
			for j, p := range ck.s4 {
				if ((p >> uint(13-i)) & 0x01) == 0x01 {
					ck.topCards = append(ck.topCards, Card(j*13+int(i)%13))
					if len(ck.topCards) == 5 {
						return ck.topCards
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
		if ck.CardFace() > ck2.CardFace() {
			return 1
		} else {
			return -1
		}
	}

	//both have the same type. compare top cards
	t1 := ck.TopCards()
	t2 := ck2.TopCards()
	for i := 0; i < len(t1); i++ {
		if t1[i] > t2[i] {
			return 1
		} else if t1[i] < t2[i] {
			return -1
		}
	}

	startIndex := len(t1)
	t1 = ck.Top5Cards()
	t2 = ck.Top5Cards()
	for i := startIndex; i < 5; i++ {
		if t1[i] > t2[i] {
			return 1
		} else if t1[i] < t2[i] {
			return -1
		}
	}

	return 0
}
