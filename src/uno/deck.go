package uno

import (
	"math/rand"
	"time"
)

const DECKSIZE = 108

type Deck struct {
	cards []*Card
}

func NewEmptyDeck() *Deck {
	cards := make([]*Card, 0, DECKSIZE)
	return &Deck{cards}
}

func NewDeck() *Deck {
	i := 0
	cards := make([]*Card, DECKSIZE, DECKSIZE)
	// number cards
	for _, color := range COLORS {
		for _, num := range NUMBERS {
			if num > 0 {
				cards[i] = NewCard(color, num, i)
				i++
			}
			cards[i] = NewCard(color, num, i)
			i++
		}
	}
	// action cards
	for _, color := range COLORS {
		for _, act := range ACTIONS {
			if act != WILD && act != PLUS4 {
				cards[i] = NewCard(color, act, i)
				cards[i+1] = NewCard(color, act, i)
				i += 2
			} else {
				cards[i] = NewCard(NO_COLOR, act, i)
				i++
			}
		}
	}
	deck := &Deck{cards}
	Shuffle(deck)
	return deck
}

func (deck *Deck) Peek() Card {
	if len(deck.cards) == 0 {
		return NILCARD
	}
	return *deck.cards[len(deck.cards)-1]
}

func (deck *Deck) Pop() *Card {
	last := len(deck.cards) - 1
	if last < 0 {
		return nil
	}
	card := deck.cards[last]
	deck.cards = deck.cards[:last]
	return card
}

func (deck *Deck) Push(card *Card) {
	deck.cards = append(deck.cards, card)
}

func (deck *Deck) Burry(card *Card) {
	deck.cards = append(deck.cards, nil)
	copy(deck.cards[1:], deck.cards)
	deck.cards[0] = card
}

func (deck *Deck) IsEmpty() bool {
	return len(deck.cards) == 0
}

func Shuffle(deck *Deck) {
	cards := deck.cards
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
}

func (deck *Deck) Reshuffle(stack *Deck, leftOver int) {
	for len(stack.cards) > 0 {
		card := stack.Pop()
		if card.IsWild() {
			card.Color = NO_COLOR
		}
		deck.Push(card)
	}
	for i := 0; i < leftOver; i++ {
		stack.Push(deck.Pop())
	}
	Shuffle(deck)
}
