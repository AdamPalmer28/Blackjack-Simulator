package main

import (
	"blackjack/config"
	"blackjack/game"
	"blackjack/sim"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// dataset := sim.CreateSimDataStructure() // create the simulation data structure
// dataset.ToJSON()
func main() {
	// Initialize configuration
	config.Init()

	fmt.Println("Welcome to the Blackjack Simulator!")

	dataset, _ := sim.LoadFromJSON("bj_sim_data.json")
	

	fmt.Println("Beginning simulating bj hands...")

	sim.SimulateBJ(1000, dataset) // run 100,000 simulations and save to dataset



	// cli menu for blackjack
	// for {
	// 		blackjackCLI()
	// 		fmt.Println("Function has ended!")
	// }

}

func mainMenuCLI() {
	fmt.Println("Welcome to the Blackjack Simulator!")
	fmt.Println("1. Start Game")
	fmt.Println("2. Exit")
	fmt.Print("Choose an option: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1] // Remove newline character
	choice, err := strconv.Atoi(input)
	if err != nil || (choice != 1 && choice != 2) {
		fmt.Println("Invalid input. Please enter 1 or 2.")
		return
	}

	if choice == 1 {
		main() // Start the game
	} else {
		fmt.Println("Exiting the game. Goodbye!")
	}
}

func blackjackCLI() {
	fmt.Println("=================================================")
	fmt.Println("Welcome to the Blackjack CLI!")
	fmt.Println("This is a simple command-line interface for playing Blackjack.")

	reader := bufio.NewReader(os.Stdin)
	
	gs := game.StartGame()// Initialize the game state
	
	for { // ! START OF TURN LOOP LOGIC
		// --------------------------------------------
		gs.Print() // Display the initial game state
		ind := gs.HandToPlay
		fmt.Printf("\n--- Hand %d ---\n", ind+1)

		moves := gs.PlayerMoves[ind]
		// 1. Hit, 2. Stand, 3. Double Down, 4. Split
		fmt.Println("Available moves:")
		if moves&0b001 != 0 {
			fmt.Println("1. Hit")
			fmt.Println("2. Stand")
		}
		if moves&0b010 != 0 {
			fmt.Println("3. Double Down")
		}						
		if moves&0b100 != 0 {
			fmt.Println("4. Split")
		}
		input, _ := reader.ReadString('\n')
		playerMove, err := strconv.Atoi(string(input[0]))
		fmt.Println("You chose:", playerMove)
		if err != nil || playerMove < 1 || playerMove > 4 {
			fmt.Println("Invalid input. Please enter a number between 1 and 4.")
			continue
		}

		switch playerMove {
		case 1: // Hit
			gs.ActionCalc(0b001) // Hit is always possible
		case 2: // Stand
			gs.ActionCalc(0b000) // Stand is always possible
		case 3: // Double Down
			if moves&0b010 != 0 {
				gs.ActionCalc(0b010) 
			} else {
				fmt.Println("You cannot double down at this time.")
				continue
			}
		case 4: // Split
			if moves&0b100 != 0 {
				gs.ActionCalc(0b100)
			} else {
				fmt.Println("You cannot split at this time.")
				continue
			}
		default:
			fmt.Println("Invalid move. Please choose a valid option.")
		}
		// --------------------------------------------
		// ! END OF USER INPUT LOGIC

		// after hand is done
		if gs.HandToPlay + 1 > len(gs.PlayerHand) {
			bjEndGame(gs)
			
			// enter to continue
			fmt.Println("Press Enter to return to the main menu...")
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')
			return
		}

		fmt.Println("Turn loop ended - restarting game loop...\n")
	}
}

func bjEndGame(gs game.GameState) {
	fmt.Println("\n\n-------------------------\nGame is over")
	fmt.Println("Dealer hand (", gs.DealerScore, "):", game.PrintCards(gs.DealerHand))
	fmt.Println("States: ", gs.State)

	// game is over
	for i, hand := range gs.PlayerHand {
		fmt.Printf("\n--- Hand %d ---", i+1)
		
		
		state := gs.State[i]
		switch state {
		case 1:
			fmt.Println("Win")
		case 2:
			fmt.Println("Loss")
		case 3:
			fmt.Println("Draw")
		case 4:
			fmt.Println("Bust")
		default:
			fmt.Println("ERROR state: ", state)
		}
		
		fmt.Println(game.PrintCards(hand))
		fmt.Println("Hand value: ", gs.HandValues[i])
		fmt.Println("Score: ", gs.PlayerScore[i])
		fmt.Println("\n")
	}

}