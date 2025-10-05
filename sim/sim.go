package sim

import (
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
	Value int // resulting value of the action
	Depth int // depth of the action in the game tree (for debugging)
}

type SimState struct {
	SimEvalData []SimEvalData // list of all simulation data
}

func SimulateBJ(hands int, dataset SimDataMap) {
	// run many simulations of the game
	fmt.Println("Starting simulation of", hands, "hands...")

	for i := 1; i <= hands; i++ {
		recentSimStates := single_player_sim()
		//if i%5000 == 0 {
		fmt.Println("Simulated", i, "hands...")

		//print recentSimStates for debugging
		for _, d := range recentSimStates.SimEvalData {
			fmt.Printf("DealerShownScore: %d, DealerScore: %d, PlayerScores: %d, PlayerHandCats: %d, ChoosenAction: %d, Value: %d\n",
				d.DealerStart, d.DealerScore, d.PlayerScores, d.PlayerHandCats, d.ChoosenAction, d.Value)
			if d.ChoosenAction == -1 {
				println("")
			}
		}

		// fmt.Println("Adding data to simulation data structure...")
		// dataset.AddData(recentSimStates)

		// fmt.Println("Saving simulation data to bj_sim_data.json...")
		// dataset.ToJSON()
		
		if i == 1 { // ! remove later - just for testing
			return
		}
	}
}

func single_player_sim() SimState {
	// run a single simulation of the game
	// return the result of the game

	gs := game.StartGame()
	gs.Print()

	simState := SimState{
		SimEvalData: make([]SimEvalData, 0),
	}

	// Start recursive exploration from initial game state
	node_explore(gs, &simState)
	
	fmt.Println("Simulation complete.")
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
func node_explore(gs game.GameState, simState *SimState) (value int) {
	// for any given hand state, explore all possible actions recursively
	// returns the value of the final outcome...

	if gs.HandToPlay >= len(gs.PlayerHand) {
		// Game over for this hand, evaluate outcome
		
		// debuging msg - print dealer score, player scores, and states
		fmt.Printf("END GAME: Dealer Score: %d", gs.DealerScore)
		for i, v := range gs.PlayerScore {
			fmt.Printf("  Player %d Score: %d", i, v)
			fmt.Printf("  Player %d State: %d\n\n", i, gs.State[i])
		}

		total := 0
		for _, v := range gs.HandValues {
			total += v
		}
		return total
		
	}
	
	// Get current hand's legal moves
	currentHandMoves := gs.PlayerMoves[gs.HandToPlay]

	actions_vals := make(map[int]int)
	for i, action := range PlayerActions {

		if (action.actionMask&currentHandMoves != 0) || action.actionInt == 0 { // Stand is always possible
			gsCopy := (&gs).Copy()
			fmt.Println("Action: ", action.actionInt,  "Current player score: ", gsCopy.PlayerScore[gs.HandToPlay])
			gsCopy.ActionCalc(action.actionMask)
			fmt.Println("After action player score: ", gsCopy.PlayerScore[gs.HandToPlay])

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

	// -----------------------------------
	// Determine value to return

	n_acts, sum_val := 0, 0
	// Determine realistic values
	for act, val := range actions_vals {
		if val == -100 {
			continue //
		}
		if gs.PlayerScore[gs.HandToPlay] > 14 {
			// pass
			act = act + 1
			// ! later we may want to account for act probabilities
			// ? currently there is a bias of downstream actions being random...
			// e.g. if score > 14 then use best strategy if score <= 14 loop through all actions
		}
		n_acts++
		sum_val += val
	}
	return sum_val / n_acts
}

// Helper function to categorize player hand
func getHandCategory(hand []game.Card) int {
	// 0: no ace, 1: has ace, 2: split available
	hasAce := false
	canSplit := false
	
	for _, card := range hand {
		if card.Rank == 1 {
			hasAce = true
		}
	}
	
	// Check if split is available (two cards of same rank)
	if len(hand) == 2 && hand[0].Rank == hand[1].Rank {
		canSplit = true
	}
	
	if canSplit {
		return 2
	} else if hasAce {
		return 1
	} else {
		return 0
	}
}