package main

import (
	"fmt"
	"strings"
)

type player struct {
	name     string
	hand     []BlackJackCard
	isDealer bool
}

func (p player) Name() string {
	var sb strings.Builder
	if p.isDealer {
		sb.WriteString("Dealer(")
	} else {
		sb.WriteString("Player(")
	}

	sb.WriteString(p.name + ")")
	return sb.String()
}

func (p player) String() string {
	return p.Name() + ": " + fmt.Sprint(p.hand)
}

func (p player) sumHand() uint8 {
	var sum uint8
	for _, c := range p.hand {
		sum += c.value()
	}

	return sum
}

func (p player) showHand() {
	for i := range p.hand {
		p.hand[i].faceDown = false
	}
}
