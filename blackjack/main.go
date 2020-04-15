package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ericywl/gophercises/deck"
)

type game struct {
	cards   []deck.Card
	players []player
}

func (g *game) dispatchCard(playerIdx int, faceDown bool) {
	bjCard := BlackJackCard{
		card:     g.cards[len(g.cards)-1],
		faceDown: faceDown,
	}
	g.players[playerIdx].hand = append(g.players[playerIdx].hand, bjCard)
	g.cards = g.cards[:len(g.cards)-1]
}

func (g *game) dispatchInitialCards() {
	// Dispatch 2 cards to each
	for i := 0; i < 2; i++ {
		for j, p := range g.players {
			g.dispatchCard(j, p.isDealer && i != 0)
		}
	}
}

func (g *game) start() {
	reader := bufio.NewReader(os.Stdin)
	for i, p := range g.players {
		if p.isDealer {
			continue
		}

		fmt.Println(g.players)
		for {
			fmt.Println("Hit (H) or Stand (s)?")
			ans, _ := reader.ReadString('\n')
			ans = strings.TrimSpace(ans)
			if ans == "" || strings.ToLower(ans) == "h" {
				g.dispatchCard(i, p.isDealer)
				fmt.Println(g.players)
				if g.checkPlayerDead(i) || g.checkPlayerMax(i) {
					break
				}

			} else {
				break
			}
		}
	}

	g.end()
}

func (g *game) checkPlayerDead(playerIdx int) bool {
	return g.players[playerIdx].sumHand() > 21
}

func (g *game) checkPlayerMax(playerIdx int) bool {
	return g.players[playerIdx].sumHand() == 21
}

func (g *game) showAllHand() {
	for i := range g.players {
		g.players[i].showHand()
	}
}

func (g *game) end() {
	winnerIdx := len(g.players) - 1
	for i := 0; i < len(g.players)-1; i++ {
		if !g.checkPlayerDead(i) && g.players[i].sumHand() > g.players[winnerIdx].sumHand() {
			winnerIdx = i
		}
	}

	g.showAllHand()
	fmt.Println(g.players)
	if g.checkPlayerDead(winnerIdx) {
		fmt.Printf("\nDraw.\n")
		return
	}

	fmt.Printf("\n" + g.players[winnerIdx].Name() + " wins!\n")
}

func NewGame() error {
	players := []player{
		{name: "Eric", isDealer: false},
		{name: "Bot", isDealer: true},
	}
	cards, err := deck.New(deck.Opts{
		Shuffle:   true,
		NumJokers: 0,
		NumDecks:  1,
	})
	if err != nil {
		return err
	}

	g := game{cards, players}
	g.dispatchInitialCards()
	g.start()
	return nil
}

func main() {
	err := NewGame()
	if err != nil {
		log.Fatalf("Error starting new game: %v", err)
	}
}
