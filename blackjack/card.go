package main

import "github.com/ericywl/gophercises/deck"

type BlackJackCard struct {
	card     deck.Card
	faceDown bool
}

func (c BlackJackCard) String() string {
	if c.faceDown {
		return "*****:*"
	}

	return c.card.String()
}

func (c BlackJackCard) value() uint8 {
	if c.card.Value > 10 {
		return 10
	}

	return c.card.Value
}
