package uno

import (
	"fmt"
)

type Player struct {
	cards map[int]*Card
	Name  string
}

func NewPlayer(name string) *Player {
	cards := make(map[int]*Card)
	return &Player{cards, name}
}

func (p *Player) HasWon() bool {
	return len(p.cards) == 0
}

func (p *Player) PlayCard(c int) *Card {
	card := p.cards[c]
	delete(p.cards, c)
	return card
}

func (p *Player) Draw(card *Card) {
	p.cards[card.Id] = card
}

func (p Player) String() string {
	return fmt.Sprintf("%s %v", p.Name, p.cards)
}
