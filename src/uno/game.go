package uno

import "fmt"

type Game struct {
	deck    *Deck
	stack   *Deck
	players *Players
	Draws   int
	turn    *Node
}

func NewGame() *Game {
	deck := NewDeck()
	stack := NewEmptyDeck()
	card := deck.Pop()
	for card.IsAction() {
		deck.Burry(card)
		card = deck.Pop()
	}
	stack.Push(card)
	players := NewPlayers()
	return &Game{deck, stack, players, 0, nil}
}

func (g *Game) AddPlayer(name string) *Player {
	p := NewPlayer(name)
	g.players.Push(p)
	fmt.Println(g.players.len, "players in game")
	if g.players.len == 1 {
		g.turn = g.players.Head
	}
	return p
}

func (g *Game) GiveHand(p *Player) []*Card {
	cards := make([]*Card, 0, 7)
	for i := 0; i < 7; i++ {
		card := g.deck.Pop()
		cards = append(cards, card)
		p.cards[card.Id] = card
	}
	return cards
}

func (g *Game) PeekStack() Card {
	if c := g.stack.Peek(); c != NILCARD {
		return c
	}
	c := g.deck.Pop()
	g.stack.Push(c)
	return *c
}

func (g *Game) CurrPlayer() *Node {
	return g.turn
}

func (g *Game) DrawCard() []*Card {
	var cards []*Card
	var n = 1
	if n < g.Draws {
		n = g.Draws
	}
	fmt.Println("draw action ", g.Draws)
	for ; n > 0; n-- {

		card := g.deck.Pop()
		if card == nil {
			g.deck.Reshuffle(g.stack, 1)
			card = g.deck.Pop()
		}
		cards = append(cards, card)
		g.turn.Player.Draw(card)
	}
	if g.Draws > 0 {
		g.turn = g.turn.next
		g.Draws = 0
	}
	return cards

}

func (g *Game) Play(c int, color string) bool {
	card := g.turn.Player.PlayCard(c)
	if card == nil {
		fmt.Println(g.turn.Player)
		return false
	}
	if g.Draws > 0 && !card.IsPlus() {
		return false
	}
	if !card.Similar(g.stack.Peek()) {
		fmt.Println("not similar ", card)
		return false
	}
	fmt.Println(card)
	g.stack.Push(card)
	g.turn = g.turn.Next()
	// aftermath
	if g.turn.Player.HasWon() {
		g.players.Remove(g.turn.Player)
	}
	g.SolveAction()
	if card.IsWild() {
		card.Color = color
	}
	return true
}

func (g *Game) SolveAction() {
	card := g.stack.Peek()
	if card.Equals(CBLOCK) {
		g.turn = g.turn.next
	} else if card.Equals(CRETURN) {
		if g.players.len == 2 {
			g.turn = g.turn.next
		}
		g.players.Reverse()
	} else if card.Equals(CPLUS2) {
		g.Draws += 2
	} else if card.Equals(CPLUS4) {
		g.Draws += 4
	}
}

func (g *Game) RemovePlayer(p *Player) {
	if p == nil {
		return
	}
	g.players.Remove(p)
	fmt.Println(p, g.turn)
	if g.turn.Player == p {
		g.turn = g.turn.next
	}
	if g.turn.next.Player == p {
		g.turn = nil
	}
	for _, card := range p.cards {
		g.deck.Burry(card)
	}
	if g.players.len == 0 {
		g.deck.Reshuffle(g.stack, 0)
	}
}

func (g *Game) PassTurn() {
	g.turn = g.turn.next
}

func (g *Game) IsEmpty() bool {
	return g.players.IsEmpty()
}
