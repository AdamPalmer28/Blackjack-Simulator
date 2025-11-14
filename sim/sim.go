package sim

import (
	"blackjack/config"
	"blackjack/game"
	"fmt"
	"time"
)

// individual simulation data piece
type SimEvalData struct {
	DealerStart   int
	DealerScore   int
	PlayerScores  int
	PlayerHandCats int // 0: no ace, 1: has ace, 2: split available
	ChoosenAction int // 0: hit, 1: stand, 2: double down, 3: split
	Value float32 // resulting value of the action
	Depth int // depth of the action in the game tree (for debugging)
}

type SimState struct {
	SimEvalData []SimEvalData // list of all simulation data
}


func SimulateBJ(hands int, dataset SimDataMap) {
	// run many simulations of the game
	startTime := time.Now()
	fmt.Printf("Starting simulation of %d million hands at %s...\n", hands / 1_000_000, startTime.Format("15:04:05"))

	debugMode := config.IsDebugMode()

	for i := 1; i <= hands; i++ {

		recentSimStates := single_player_sim(&dataset)

		//fmt.Println("Adding data to simulation data structure...")
		dataset.AddData(recentSimStates)

		if debugMode {
		for _, d := range recentSimStates.SimEvalData {
			fmt.Printf("DSS: %d, DS: %d, S: %d, cat: %d, Act: %d, V: %f\n",
				d.DealerStart, d.DealerScore, d.PlayerScores, d.PlayerHandCats, d.ChoosenAction, d.Value)
				if d.ChoosenAction == -1 {
					println("")
				}
			}
		} 
		if i%1_000_000 == 0 {
				elapsed := time.Since(startTime)
				handsPerSecond := float64(i) / elapsed.Seconds()
				progress := float64(i) / float64(hands) * 100
				estTimeRemaining := float64(hands - i) / handsPerSecond 

				fmt.Printf("Progress: %.2f%% \n", progress)
				fmt.Printf("Simulated %d million hands (%.2f hands/sec, elapsed: %s)\n", 
					i/1_000_000, handsPerSecond, elapsed.Round(time.Second))
				fmt.Printf("Est. time remaining: %.2f seconds\n\n", estTimeRemaining)


				//fmt.Println("Saving simulation data to bj_sim_data.json...")
				dataset.ToJSON()
			
		}
	
	}
	
	totalElapsed := time.Since(startTime)
	finalRate := float64(hands) / totalElapsed.Seconds()
	fmt.Printf("Simulation completed! Total time: %s (%.2f hands/sec)\n", 
		totalElapsed.Round(time.Millisecond), finalRate)
}

func single_player_sim(dataset *SimDataMap) SimState {
	// run a single simulation of the game
	// return the result of the game

	gs := game.StartGame()
	if config.IsDebugMode() {
		gs.Print()
	}	

	simState := SimState{
		SimEvalData: make([]SimEvalData, 0),
	}

	// Start recursive exploration from initial game state
	node_explore(gs, &simState, dataset)
	
	if config.IsDebugMode() {
	fmt.Println("Simulation complete.")
	}
	return simState
}


// Try all possible actions for current hand
// Stand (action 0) - always available
var PlayerActions = []struct {
	actionInt   int // 0: Stand, 1: Hit, 2: Double Down, 3: Split
	actionMask   int
}{
	{0, 0b000}, // Stand
	{1, 0b001}, // Hit
	{2, 0b010}, // Double Down
	{3, 0b100}, // Split
}

// ! I have rewritten this but not working properly...
func node_explore(gs game.GameState, simState *SimState, dataset *SimDataMap) (value float32) {
	// for any given hand state, explore all possible actions recursively
	// returns the value of the final outcome...

	if gs.HandToPlay >= len(gs.PlayerHand) {
		// !GAME OVER - will exit here
		// for this hand, evaluate outcome
		
		if config.IsDebugMode() {
			fmt.Printf("END GAME: DS: %d", gs.DealerScore)
			for i, v := range gs.PlayerScore {
				fmt.Printf("  P %d S: %d", i, v)
				fmt.Printf("  State: %d\n", gs.State[i])
			}
		}
		total := float32(0)

		for ind, v := range gs.HandValues {
			if config.IsDebugMode() {
				fmt.Printf("V%d: %d \n\n", ind, v)
			}
			total += float32(v)
		}
		return total
		
	}
	// Get current hand's legal moves
	currentHandMoves := gs.PlayerMoves[gs.HandToPlay]

	// Get hand category once before the loop
	hand_cat := getHandCategory(gs.PlayerHand[gs.HandToPlay])

	// ! MAIN LOOP
	if config.IsDebugMode() {
		fmt.Println("new loop  ",len(simState.SimEvalData))
	}
	actions_vals := make(map[int]float32) // map of action index to value
	for i, action := range PlayerActions {
		// do all actions...

		if (action.actionMask&currentHandMoves != 0) || action.actionInt == 0 { // Stand is always possible
			gsCopy := (&gs).Copy()
			if config.IsDebugMode() {
				fmt.Printf("Act: %d || S: %d", action.actionInt, gsCopy.PlayerScore[gs.HandToPlay])
			}
			gsCopy.ActionCalc(action.actionMask)

			if config.IsDebugMode() {
				fmt.Printf(" || Post S: %d\n", gsCopy.PlayerScore[gs.HandToPlay])
			}

			// Continue exploring from this state
			value = node_explore(gsCopy, simState, dataset)

			simData := SimEvalData{
				DealerStart:    gs.DealerShownScore,
				DealerScore:    gs.DealerScore,
				PlayerScores:   gs.PlayerScore[gs.HandToPlay],
				PlayerHandCats: hand_cat,
				ChoosenAction:  action.actionInt,
				Value:          value,  // Will be updated when game finishes
			}
			simState.SimEvalData = append(simState.SimEvalData, simData)

		} else {
			// Action not possible, skip
			value = -100
		}
		actions_vals[i] = value
	}
	if config.IsDebugMode() {
		fmt.Printf("<> FIN All Act  %d\n", gs.PlayerScore[gs.HandToPlay])
	}
	// ----------------------------------
	// Determine value to return
	// ! FUNCTION EXIT - if game not over

	// Find the best action based on expected values from the dataset
	var best_action int = 0  // default to stand
	var best_expected_value float32 = -1000
	
	// Check if we have data for this state in our dataset
	if dealerMap, ok := (*dataset)[gs.DealerShownScore]; ok {
		if playerMap, ok := dealerMap[gs.PlayerScore[gs.HandToPlay]]; ok {
			if categoryMap, ok := playerMap[hand_cat]; ok {
				// Find action with highest expected value
				for action, simData := range categoryMap {
					if simData.Trials > 0 && simData.ExpectedValue > best_expected_value {
						best_expected_value = simData.ExpectedValue
						best_action = action
					}
				}
			}
		}
	}

	// Determine returned score - use the best action's value
	final_val := actions_vals[best_action]

	if config.IsDebugMode() {
		fmt.Printf("Best action: %d, returned V: %f\n", best_action, final_val)
	}


	return final_val
}

// Helper function to categorize player hand
func getHandCategory(hand []game.Card) int {
	// 0: no ace, 1: has ace, 2: split available
	hasAce := false
	canSplit := false
	rank_sum := 0
	
	for _, card := range hand {
		if card.Rank == 1 {
			hasAce = true
		}
		rank_sum += card.Rank
	}

	// Check if split is available (two cards of same rank) 
	if len(hand) == 2 && hand[0].Rank == hand[1].Rank {
		canSplit = true
	}
	
	if canSplit {
		return 2
	} else if hasAce && rank_sum <= 11 {
		return 1
	} else {
		return 0
	}
}