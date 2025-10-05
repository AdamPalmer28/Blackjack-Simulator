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

		fmt.Println("Adding data to simulation data structure...")
		dataset.AddData(recentSimStates)

		fmt.Println("Saving simulation data to bj_sim_data.json...")
		dataset.ToJSON()
		
		if i == 0 { // ! remove later - just for testing
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
	
	return simState
}

// ! NEED TO REDO THIS ENTIRE FUNCTION WITHOUT USING AI...
func node_explore(gs game.GameState, simState *SimState) {
	// for any given hand state, explore all possible actions recursively
	
	// Check if game is f
	// inished (all hands played)
	if gs.HandToPlay >= len(gs.PlayerHand) {
		// Game is over, record final results
		for i := 0; i < len(gs.State); i++ {
			simData := SimEvalData{
				DealerStart:    gs.DealerShownScore,
				DealerScore:    gs.DealerScore,
				PlayerScores:   gs.PlayerScore[i],
				PlayerHandCats: getHandCategory(gs.PlayerHand[i]),
				ChoosenAction:  -1, // Final state, no action chosen
				Value:         gs.HandValues[i],
			}
			simState.SimEvalData = append(simState.SimEvalData, simData)
		}
		return gs.HandValues[i]
	}

	// Get current hand's legal moves
	currentHandMoves := gs.PlayerMoves[gs.HandToPlay]
	
	// Try all possible actions for current hand
	// Stand (action 0) - always available
	actions := []struct {
		actionInt   int // 0: Stand, 1: Hit, 2: Double Down, 3: Split
		actionMask   int
	}{
		{0, 0b000}, // Stand
		{1, 0b001}, // Hit
		{2, 0b010}, // Double Down
		{3, 0b100}, // Split
	}

	for _, action := range actions {

		if (action.actionMask&currentHandMoves != 0) || action.actionInt == 0 { // Stand is always possible
			gsCopy := (&gs).Copy()
			gsCopy.ActionCalc(action.actionMask)

			simData := SimEvalData{
				DealerStart:    gs.DealerShownScore,
				DealerScore:    gs.DealerScore,
				PlayerScores:   gs.PlayerScore[gs.HandToPlay],
				PlayerHandCats: getHandCategory(gs.PlayerHand[gs.HandToPlay]),
				ChoosenAction:  action.actionInt,
				Value:          0,  // Will be updated when game finishes
			}
			simState.SimEvalData = append(simState.SimEvalData, simData)

			// Continue exploring from this state
			node_explore(gsCopy, simState)
		}
	}

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