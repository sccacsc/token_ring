package main

import (
	"crypto/sha3"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Token struct {
	Data     string
	DestHash [32]byte
	TTL      int
}

type Node struct {
	ID  int
	In  chan Token
	Out chan Token
}

func hashID(id int) [32]byte {
	return sha3.Sum256([]byte(strconv.Itoa(id)))
}

func (n *Node) run(totalNodes int) {
	for token := range n.In {
		fmt.Printf("[Node %d] received token: data=%q ttl=%d\n",
			n.ID, token.Data, token.TTL)

		// Проверка адресата
		if token.DestHash == hashID(n.ID) {
			fmt.Printf("[Node %d] token accepted: %q\n", n.ID, token.Data)

			// Генерация нового сообщения
			dest := rand.Intn(totalNodes) + 1
			newToken := Token{
				Data:     fmt.Sprintf("msg from %d to %d", n.ID, dest),
				DestHash: hashID(dest),
				TTL:      rand.Intn(totalNodes) + 1,
			}

			fmt.Printf("[Node %d] generated new token → %d\n", n.ID, dest)
			n.Out <- newToken
			continue
		}

		// Пересылка дальше
		token.TTL--
		if token.TTL <= 0 {
			fmt.Printf("[Node %d] token expired\n", n.ID)
			continue
		}

		fmt.Printf("[Node %d] forwarding token\n", n.ID)
		n.Out <- token
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: tokenring <number_of_nodes>")
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n < 2 {
		fmt.Println("Number of nodes must be >= 2")
		return
	}

	rand.Seed(time.Now().UnixNano())

	nodes := make([]Node, n)
	channels := make([]chan Token, n)

	for i := 0; i < n; i++ {
		channels[i] = make(chan Token)
	}

	for i := 0; i < n; i++ {
		nodes[i] = Node{
			ID:  i + 1,
			In:  channels[(i-1+n)%n],
			Out: channels[i],
		}
		go nodes[i].run(n)
	}

	// Первое сообщение от main
	dest := rand.Intn(n) + 1
	initial := Token{
		Data:     "initial message",
		DestHash: hashID(dest),
		TTL:      n,
	}

	fmt.Printf("[Main] sending initial token → %d\n", dest)
	channels[0] <- initial

	select {}
	/*for {
		time.Sleep(time.Second)
	}*/

}
