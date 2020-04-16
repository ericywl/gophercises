package deck

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"errors"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrAllWildcard occurs when two wildcards are provided for Card{} in opts.Filter
	ErrAllWildcard = errors.New("cannot filter with double wildcard")
	// ErrInvalidSuit occurs when user provides invalid suit for Card{} in opts.Filter
	ErrInvalidSuit = errors.New("invalid suit")
	// ErrInvalidValue occurs when user provides invalid value for Card{} in opts.Filter
	ErrInvalidValue = errors.New("invalid value")
)

const (
	Wildcard uint8 = iota
	Spade
	Diamond
	Heart
	Club
	Joker
)

// Card contains a Suit and a Value.
type Card struct {
	Suit  uint8
	Value uint8
}

func (c Card) String() string {
	if c.Suit == Joker {
		return "Joker"
	}

	var sb strings.Builder
	switch c.Suit {
	case Spade:
		sb.WriteString("Spade")
	case Diamond:
		sb.WriteString("Diamond")
	case Heart:
		sb.WriteString("Heart")
	case Club:
		sb.WriteString("Club")
	}

	sb.WriteString(":")

	switch c.Value {
	case 1:
		sb.WriteString("A")
	case 11:
		sb.WriteString("J")
	case 12:
		sb.WriteString("Q")
	case 13:
		sb.WriteString("K")
	default:
		sb.WriteString(strconv.FormatUint(uint64(c.Value), 10))
	}

	return sb.String()
}

// Opts defines a set of options that can be passed when creating deck of cards.
type Opts struct {
	SortLessFn func(a, b int) bool
	Shuffle    bool
	NumJokers  uint8
	NumDecks   uint8
	Filter     []Card
}

func buildFilterMap(filter []Card) (map[Card]bool, error) {
	filterMap := map[Card]bool{}
	if filter == nil {
		return filterMap, nil
	}

	for _, c := range filter {
		if c.Suit == Wildcard && c.Value == Wildcard {
			return nil, ErrAllWildcard
		}

		if c.Suit > Joker {
			return nil, ErrInvalidSuit
		}

		if c.Value > 13 {
			return nil, ErrInvalidValue
		}

		if c.Suit == Wildcard {
			for s := Spade; s <= Club; s++ {
				filterMap[Card{s, c.Value}] = true
			}
			continue
		}

		if c.Value == Wildcard {
			for v := uint8(1); v <= 13; v++ {
				filterMap[Card{c.Suit, v}] = true
			}

			continue
		}

		filterMap[c] = true
	}

	return filterMap, nil
}

func seed() {
	var b [8]byte
	_, err := cryptoRand.Read(b[:])
	if err != nil {
		rand.Seed(time.Now().UTC().UnixNano())
		return
	}

	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

// New creates new deck of cards
func New(opts Opts) ([]Card, error) {
	seed()
	var cards []Card

	filterMap, err := buildFilterMap(opts.Filter)
	if err != nil {
		return nil, err
	}

	numDecks := opts.NumDecks
	if numDecks == 0 {
		numDecks = 1
	}

	for i := uint8(0); i < numDecks; i++ {
		for s := Spade; s <= Club; s++ {
			for v := uint8(1); v <= 13; v++ {
				card := Card{s, v}
				if _, ok := filterMap[card]; !ok {
					cards = append(cards, card)
				}
			}
		}
	}

	if opts.NumJokers > 0 {
		for i := uint8(0); i < opts.NumJokers; i++ {
			cards = append(cards, Card{Joker, 14})
		}
	}

	if opts.SortLessFn != nil {
		sort.Slice(cards, opts.SortLessFn)
	}

	if opts.Shuffle {
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	return cards, nil
}
