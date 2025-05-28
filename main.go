package main

import (
	"blackjack/game"
)

func main() {
	//fmt.Println("Hello, World!")

	gs := game.StartGame()
	
	gs.Print()

}