package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/memetea/dzpk"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	iSize := 10000000
	var gens [10000000]*dzpk.CardsCheck
	for i := 0; i < iSize; i++ {
		var cc dzpk.CardCollection
		n := 7
		for n > 0 {
			randK := dzpk.Card(rand.Int31n(52))
			if cc.CardExists(randK) {
				continue
			}
			cc.SetCard(randK)
			n--
		}
		gens[i] = dzpk.NewCardsCheckFromCC(cc)
	}

	//var cardTypeMap [9]int
	t0 := time.Now()
	for i := 0; i < iSize; i++ {
		//t := gens[i].CardFace()
		gens[i].Top5Cards()
		//cardTypeMap[t]++
		//AnalysePoker(gens[i])
	}
	t1 := time.Now()
	fmt.Printf("--------------------The call took %v to run.------------------------\n", t1.Sub(t0))

	// var sum float64
	// for i := 0; i < len(cardTypeMap); i++ {
	// 	sum += float64(cardTypeMap[i])
	// }

	// for i := 0; i < len(cardTypeMap); i++ {
	// 	fmt.Printf("%d percent: %f\n", i, float64(cardTypeMap[i])/sum)
	// }
}
