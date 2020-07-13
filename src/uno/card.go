package uno

import "fmt"

const BLUE = "Blue"
const RED = "Red"
const GREEN = "Green"
const YELLOW = "Yellow"
const NO_COLOR = "Wild"

var COLORS = [...]string{BLUE, RED, GREEN, YELLOW}

var NUMBERS = [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

const BLOCK = 10
const RETURN = 11
const PLUS2 = 12
const WILD = 13
const PLUS4 = 14

var ACTIONS = [...]int{BLOCK, RETURN, PLUS2, WILD, PLUS4}
var SACTIONS = [...]string{"-B", "-R", "+2", "-W", "+4"}

var SYMBOLS = map[int]string{
	0:  "0",
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "B",
	11: "R",
	12: "+2",
	13: "",
	14: "+4",
}

var CBLOCK = Card{NO_COLOR, BLOCK, SYMBOLS[BLOCK], true, 0}
var CRETURN = Card{NO_COLOR, RETURN, SYMBOLS[RETURN], true, 0}
var CPLUS2 = Card{NO_COLOR, PLUS2, SYMBOLS[PLUS2], true, 0}
var CWILD = Card{NO_COLOR, WILD, SYMBOLS[WILD], true, 0}
var CPLUS4 = Card{NO_COLOR, PLUS4, SYMBOLS[PLUS4], true, 0}
var NILCARD = Card{"*", -1, "*", false, -1}

type Card struct {
	Color  string `json:"color"`
	Number int    `json:"number"`
	Symbol string `json:"symbol"`
	Action bool   `json:"action"`
	Id     int    `json:"id"`
}

func NewCard(Color string, number, id int) *Card {
	return &Card{Color, number, SYMBOLS[number], number > 9, id}
}

func (c Card) IsAction() bool {
	return c.Action
}

func (c Card) IsWild() bool {
	return c.Action && (c.Number == WILD || c.Number == PLUS4)
}

func (c Card) IsPlus() bool {
	return c.Action && (c.Number == PLUS2 || c.Number == PLUS4)
}

func (c *Card) Similar(c2 Card) bool {
	return c2.Color == c.Color || c.Number == c2.Number || c2.Color == NO_COLOR || c.Color == NO_COLOR
}

func (c *Card) Equals(c2 Card) bool {
	return (c.Color == c2.Color || c.Color == NO_COLOR || c2.Color == NO_COLOR) && c.Number == c2.Number
}

func (c Card) String() string {
	return fmt.Sprintf("%s %s", c.Color, c.Symbol)
}
