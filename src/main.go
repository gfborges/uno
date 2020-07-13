package main

import (
	"fmt"
	"log"
	"net/http"
	"uno/src/uno"

	"github.com/gorilla/websocket"
)

var up = websocket.Upgrader{}

var game = uno.NewGame()

var turns = make(map[string]*websocket.Conn)

var stack = make(map[string]*websocket.Conn)

type message struct {
	Card    uno.Card `json:"card"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	MyTurn  bool     `json:"myTurn"`
	Draws   int      `json:"draws"`
}

var yourTurn = message{Content: "It is your turn !", MyTurn: true}

var notYourTurn = message{Content: "It isn't your turn !", MyTurn: false}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws/game", handleGame)
	http.HandleFunc("/ws/stack", handleStack)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	ws, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()
	var msg message
	err = ws.ReadJSON(&msg)
	if err != nil {
		log.Println("ws/game ReadJSON", err)
		return
	}
	name := msg.Name
	p := game.AddPlayer(name)
	defer game.RemovePlayer(p)
	hand := game.GiveHand(p)
	for _, card := range hand {
		ws.WriteJSON(writeCard(*card, name))
	}
	changeStack()
	for !p.HasWon() {
		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Println("ws/game ReadJSON", name, err)
			return
		}
		if p != game.CurrPlayer().Player {
			continue
		}
		fmt.Print(name)
		if msg.Card.Id == -2 {
			fmt.Println(" passed the turn")
			game.PassTurn()
		}
		if msg.Card.Id == -1 {
			fmt.Println(" drawed a card")
			cards := game.DrawCard()
			for _, card := range cards {
				fmt.Println(card)
				ws.WriteJSON(writeCard(*card, name))
			}
		}
		if msg.Card.Id == -3 {
			fmt.Println(" drawed a card")
			cards := game.DrawCard()
			for _, card := range cards {
				fmt.Println(card)
				ws.WriteJSON(writeCard(*card, name))
			}
			game.PassTurn()
		}
		if msg.Card.Id >= 0 {
			fmt.Println(" played a card")
			game.Play(msg.Card.Id, msg.Card.Color)
		}
		changeStack()
	}
	fmt.Println(name, " won !")
	sendWinner(name)
	if game.IsEmpty() {
		game = uno.NewGame()
	}
}

func handleStack(w http.ResponseWriter, r *http.Request) {
	ws, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws/stack Upgrade", err)
	}
	defer ws.Close()
	var msg message
	ws.ReadJSON(&msg)
	name := msg.Name
	stack[name] = ws
	msg.Card = game.PeekStack()
	msg.Draws = game.Draws
	if game.CurrPlayer() == nil {
		msg.Name = "Your"
		msg.MyTurn = true
	} else {
		msg.Name = game.CurrPlayer().Player.Name
	}
	ws.WriteJSON(msg)
	for {
		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Println("ws/stack ReadJSON", err)
			delete(stack, name)
			break
		}
		msg.Card = game.PeekStack()
		fmt.Println("Draws ", game.Draws)
		msg.Draws = game.Draws
		msg.Name = game.CurrPlayer().Player.Name
		ws.WriteJSON(msg)
	}
}

func changeStack() {
	msg := message{Card: game.PeekStack()}
	fmt.Println("Stack card: ", msg.Card)
	fmt.Println("Turn: ", game.CurrPlayer().Player.Name)
	for name := range stack {
		msg.Draws = game.Draws
		msg.Name = game.CurrPlayer().Player.Name
		if msg.Name == name {
			msg.MyTurn = true
			msg.Name = "Your"
		}
		err := stack[name].WriteJSON(msg)
		if err != nil {
			log.Println("ws/stack WriteJSON", err)
			stack[name].Close()
			delete(stack, name)
		}
	}
}

func sendWinner(name string) {
	msg := message{Content: fmt.Sprintf("%s %s", name, "won the game"), Name: "*"}
	for player := range stack {
		err := stack[player].WriteJSON(msg)
		if err != nil {
			log.Println("ws/stack WriteJSON ", err)
		}
	}
}

func writeCard(card uno.Card, name string) message {
	return message{Name: name, Content: "", Card: card}
}
