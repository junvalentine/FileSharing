package main

import (
	"fmt"
	"time"
	"bufio"
	"os"
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
		fmt.Printf("1. discover [hostname]\n")
		fmt.Printf("2. ping [hostname]\n")
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
		} else if input[:4] == "ping" {
			fmt.Printf("incoming\n")
			//To do
		} else if len(input) >= 8 && input[:8] == "discover" {
			fmt.Printf("incoming\n")
			//To do
		} else{
			fmt.Printf("Incorrect command\n")
		}
	}

}
