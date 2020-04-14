package main

import (
	"fmt"
	"log"

	"github.com/ericywl/gophercises/deck/deck"
)

func main() {
	cards, err := deck.New(deck.Opts{
		SortLessFn: nil,
		Shuffle:    false,
		NumJokers:  2,
		Filter: []deck.Card{
			{
				Suit:  deck.Wildcard,
				Value: 1,
			},
		},
		NumDecks: 2,
	})

	if err != nil {
		log.Printf("Error generating deck: %v", err)
		return
	}

	fmt.Println(cards)
}
