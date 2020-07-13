package uno

import "fmt"

type Players struct {
	Head *Node
	len  int
}

type Node struct {
	Player *Player
	next   *Node
}

func NewPlayers() *Players {
	return &Players{nil, 0}
}

func (ps *Players) IsEmpty() bool {
	return ps.Head == nil
}

func (ps *Players) Push(p *Player) {
	curr := ps.Head
	if curr == nil {
		ps.Head = &Node{p, nil}
		ps.Head.next = ps.Head
		ps.len = 1
		return
	}
	for curr.next != ps.Head {
		curr = curr.next
	}
	curr.next = &Node{p, ps.Head}
	ps.len++
}

func (ps *Players) GetPlayer(name string) *Player {
	curr := ps.Head
	if ps.Head == nil {
		return nil
	}
	fmt.Println(curr, curr.Player, curr.Player.Name)
	for curr.Player.Name != name {
		curr = curr.next
	}
	return curr.Player
}

func (ps *Players) Remove(p *Player) {
	curr := ps.Head
	if curr == nil {
		return
	}
	if curr.Player == p {
		if ps.len == 1 {
			ps.Head = nil
		} else {
			ps.Head = curr.next
		}
		ps.len = 0
		return
	}
	flag := true
	for curr.next != ps.Head {
		if curr.next.Player == p {
			flag = false
			curr.next = curr.next.next
			ps.len--
		}
	}
	if flag {
		fmt.Println("not found")
	}
}

func (ps *Players) Pop() *Player {
	curr := ps.Head
	for curr.next.next != ps.Head {
		curr = curr.next
	}
	p := curr.next.Player
	curr.next = ps.Head
	ps.len--
	return p
}

func (ps *Players) IsInside(name string) bool {
	curr := ps.Head
	for curr.next != ps.Head {
		if curr.Player.Name == name {
			return true
		}
	}
	return false
}

func (n *Node) Next() *Node {
	return n.next
}

func (ps *Players) Reverse() {
	curr := ps.Head
	var prev *Node = nil
	var next *Node = nil
	for curr != nil {
		next = curr.next
		curr.next = prev
		prev = curr
		curr = next
	}
	ps.Head = prev.next

}
