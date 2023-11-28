package main

import (
	"fmt"
	"time"
	"bufio"
	"os"
	"strings"
)

func main() {
	var port string
	fmt.Printf("Please enter a port number: ")
	fmt.Scanf("%s", &port)

	var loc string
	fmt.Printf("Please enter your local repository location: ")
	fmt.Scanf("%s", &loc)

	start := time.Now()
	p := MakePeer(loc, ":"+ port)
	t1 := time.Now()
	elapsed := t1.Sub(start)

	p.Welcome()
	fmt.Printf("Total start time: %v\n", elapsed)
	
	t2 := time.Now()
	p.ConnectServer()
	t3 := time.Now()
	peerConnectServerTime := t3.Sub(t2)
	fmt.Printf("Peer connect to server time : %v\n", peerConnectServerTime)

	for true {

		fmt.Printf("\nPlease enter a command: \n")
		fmt.Printf("1. publish [lname] [fname]\n")
		fmt.Printf("2. fetch [fname]\n")
		fmt.Printf("3. exit\n")

		var input string
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		// fmt.Printf("You entered: %s", input)

		if len(input) < 4 {
			fmt.Printf("Incorrect command\n")
		} else if input[:4] == "exit"{
			fmt.Printf("Server shutting down\n")
			break
		} else if len(input) >= 5 && input[:5] == "fetch" {
			words := strings.Split(input, " ")
			if len(words) != 2 {
				fmt.Printf("Incorrect command\n")
			} else {
				p.SearchForFile(strings.TrimSpace(words[1]))
			}
			//To do
		} else if len(input) >= 7 && input[:7] == "publish" {
			words := strings.Split(input, " ")
			if len(words) != 3 {
				fmt.Printf("Incorrect command\n")
			} else {
				filePath:=strings.TrimSpace(words[1])+strings.TrimSpace(words[2])
				_, err := os.Stat(filePath)
				if os.IsNotExist(err) {
					fmt.Printf("File not exist in your local file system")
					continue
				}
				p.RegisterFile(strings.TrimSpace(words[2]),strings.TrimSpace(words[1]))
				fmt.Printf("Register file %s%s\n", strings.TrimSpace(words[1]), strings.TrimSpace(words[2]))
			}
			//To do
		} else{
			fmt.Printf("Incorrect command\n")
		}
	}

}
