package main

import (
	"fmt"
	"time"
	"bufio"
	"os"
	"strings"
	"strconv"
)

func main() {

	start := time.Now()
	m := MakeServer()
	t1 := time.Now()
	elapsed := t1.Sub(start)

	m.Welcome()
	fmt.Printf("Total start time: %v\n", elapsed)
	
	for true {

		fmt.Printf("\nPlease enter a command: \n")
		fmt.Printf("1. discover [hostname/PeerID]\n")
		fmt.Printf("2. ping [hostname/PeerID]\n")
		fmt.Printf("3. list\n")
		fmt.Printf("4. exit\n")

		var input string
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		// fmt.Printf("You entered: %s", input)

		if len(input) < 4 {
			fmt.Printf("Incorrect command\n")
		} else if input[:4] == "exit"{
			fmt.Printf("Server shutting down\n")
			break
		} else if input[:4] == "ping" {
			words := strings.Split(input, " ")
			
			if len(words) != 2 {
				fmt.Printf("Incorrect command\n")
			} else {
				peerID, err := strconv.Atoi(strings.TrimSpace(words[1]))
				if err != nil {
					fmt.Printf("Invalid PeerID")
					continue
				}
				m.PingPeer(peerID)
			}
		} else if input[:4] == "list"{
			m.ListPeers()
		} else if len(input) >= 8 && input[:8] == "discover" {
			words := strings.Split(input, " ")
			
			if len(words) != 2 {
				fmt.Printf("Incorrect command\n")
			} else {
				peerID, err := strconv.Atoi(strings.TrimSpace(words[1]))
				if err != nil {
					fmt.Printf("Invalid PeerID")
					continue
				}
				m.DiscoverFile(peerID)
			}
			//To do
		} else{
			fmt.Printf("Incorrect command\n")
		}
	}

}
