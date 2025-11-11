package sim

import (
	"blackjack/config"
	"blackjack/game"
	"fmt"
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
	fmt.Println("Starting simulation of", hands, "hands...")

	debugMode := config.IsDebugMode()

	for i := 1; i <= hands; i++ {

		recentSimStates := single_player_sim()

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
		if i%100 == 0 {
				fmt.Println("Simulated", i, "hands...")
		
		
				fmt.Println("Saving simulation data to bj_sim_data.json...")
				dataset.ToJSON()
			
		}
	
	}
}

func single_player_sim() SimState {
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
	node_explore(gs, &simState)
	
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
func node_explore(gs game.GameState, simState *SimState) (value float32) {
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
			value = node_explore(gsCopy, simState)

			simData := SimEvalData{
				DealerStart:    gs.DealerShownScore,
				DealerScore:    gs.DealerScore,
				PlayerScores:   gs.PlayerScore[gs.HandToPlay],
				PlayerHandCats: getHandCategory(gs.PlayerHand[gs.HandToPlay]),
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

	n_acts, sum_val := 0, float32(0)

	// Determine returned score
	//     currently this is just mean...
	for act, val := range actions_vals {
		if val == -100 {
			continue //
		}
		if gs.PlayerScore[gs.HandToPlay] > 14 {
			// pass
			act = act + 1 // dummy operation
			// ! TODO - we may want to account for act probabilities
			// ? currently there is a bias of downstream actions being random...
			// e.g. if score > 14 then use best strategy if score <= 14 loop through all actions
		}
		n_acts++
		sum_val += val
	}
	final_val := float32(sum_val) / float32(n_acts)

	if config.IsDebugMode() {
		fmt.Printf("returned V: %f / %d = %f\n", sum_val, n_acts, float64(sum_val)/float64(n_acts))
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