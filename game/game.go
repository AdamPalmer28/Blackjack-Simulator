package game

import (
	"strconv"
)

type GameState struct {

	// Game state
	Deck  Deck
	State []int // 0: player move, 1: dealer move, 2: player win, 3: dealer win, 4: draw, 5: player bust
	Value []int // value of each hand - used for splitting / doubling down
	
	HandToPlay int // player hand to play
	PlayerMoves []int // legal moves for the player (hit, double down, split) - 0b000: no moves, 0b001: hit, 0b010: double down, 0b100: split
	HandValues []int

	// Player's hand
	PlayerHand  [][]Card // allow for splitting hands
	PlayerScore []int
	playerAce   []bool // true if player has an Ace in their hand
	Moves       []int  // legal moves (hit, double down, split)

	// Dealer's hand
	DealerHand       []Card
	dealerShownScore int  // score of the dealer's shown card
	dealerShownAce   bool // true if dealer's shown card is an Ace

	dealerAce   bool // true if dealer has an Ace in their hand
	DealerScore int
}



// ----------------------------------------------------------------------------
// GAME LOGIC

// ! Worker function does the inreactions with gamestate
func (gs *GameState) ActionCalc(playerMove int) {
	// Acts as next step in the game logic - playerMove if player has not stood yet
	active_turn := true // true if player can still act

	if playerMove == 0 { // stand
		active_turn = false 

	} else if playerMove == 0b001 { // hit
		gs.drawCard(gs.HandToPlay)

	} else if playerMove == 0b010 { // double down
		active_turn = false
		gs.HandValues[gs.HandToPlay] *= 2
		gs.drawCard(gs.HandToPlay)

	} else if playerMove == 0b100 { // split
		// Player splits their hand into two hands
		// This will be handled in the next turn

		// TODO:
	} else {
		panic("Error: Invalid player move: " + strconv.Itoa(playerMove))
	}

	// ----------- AFTER ACTION -----------

	// if legal moves go back to user...
	if !active_turn {
		gs.HandToPlay++
		
		if gs.HandToPlay + 1 >= len(gs.PlayerHand) {
			// All player hands have been played, now it's the dealer's turn
			// TODO game_end
			return
		}	
		
	}
	gs.UpdatePlayerState()
	return // return back to the game loop
}


// Initialize a new game state
func StartGame() GameState {

	gs := GameState{
		Deck: newDeck(),
		// State of play
		HandToPlay: 0,
		PlayerMoves: make([]int, 0), // legal moves (hit, double down, split)
		HandValues: make([]int, 1),

		// Player Hands 
		PlayerHand:  make([][]Card, 0), // Start with no player hands
		PlayerScore: make([]int, 0),
		playerAce:   make([]bool, 0),

		// Dealer Hands
		DealerHand:       make([]Card, 0),
		DealerScore:      0,
		dealerShownScore: 0,
		dealerAce:        false,
		dealerShownAce:   false, // ? do I need this?
	}

	// Deal initial cards to player and dealer
	gs.dealInitialCards()
	
	// update player states
	gs.UpdatePlayerState()

	// TODO
	// calculate initial dealers state


	return gs
}


func (gs *GameState) calcPlayMoves() {
	// Calculate the legal moves for the player based on their hand
	// Legal moves are double, split, hit (111) note stand is always possible
	playerMove := gs.HandToPlay

	// Reset moves
	legalMoves := 0b001

	if gs.PlayerScore[playerMove] < 21 {
		// player can double
		legalMoves |= 0b010
	}
	if gs.PlayerHand[playerMove][0].Rank == gs.PlayerHand[playerMove][1].Rank {
		// player can split
		legalMoves |= 0b100
	}
	gs.PlayerMoves[playerMove] = legalMoves
}

func (gs *GameState) UpdatePlayerState() {
	// player score, ace, legal moves
	ind := gs.HandToPlay

	// score
	gs.PlayerScore[ind] = calculateScore(gs.PlayerHand[ind])
	
	// Check if player has an Ace
	gs.playerAce[ind] = false
	for _, card := range gs.PlayerHand[ind] {
		if card.Rank == 1 {
			gs.playerAce[ind] = true
			break
		}
	}
	// legal moves for the player
	gs.calcPlayMoves()
}

func (gs *GameState) endGame() {
	// TODO: runs once game is over
	// Computes dealer hand/moves + final state computation
}

// --------------------------
// gamestate helper functions
// --------------------------

func (gs *GameState) dealInitialCards() {
	// Deal two cards to the player

	playerHand := make([]Card, 0)
	for i := 0; i < 2; i++ {
		playerHand = append(playerHand, gs.Deck.Draw())
	}
	gs.PlayerHand = append(gs.PlayerHand, playerHand)
	// Deal two cards to the dealer
	for i := 0; i < 2; i++ {
		gs.DealerHand = append(gs.DealerHand, gs.Deck.Draw())
	}

	// Update player score and ace status
	gs.PlayerScore = append(gs.PlayerScore, calculateScore(playerHand))
	gs.playerAce = append(gs.playerAce, false) // Initialize ace status
	gs.PlayerMoves = append(gs.PlayerMoves, 0b001) // Player can hit or stand initially
	gs.HandValues = append(gs.HandValues, 1)
}

func (gs *GameState) drawCard(hand_ind int) {
	// draw card into hand
	new_card := gs.Deck.Draw()
	gs.PlayerHand[hand_ind] = append(gs.PlayerHand[hand_ind], new_card)
}

func calculateScore(hand []Card) int {
	// Calculate the score of a hand
	score := 0
	for _, card := range hand {
		if card.Rank > 10 {
			// J, Q, K are worth 10 points
			score += 10
			continue
		}
		score += card.Rank // Add the rank of the card
	}
	return score
}



// ============================================================================
// HELPER FUNCTIONS

func printScore(hand []Card) string {
	// Display score of a hand
	score := 0
	ace := false
	for _, card := range hand {
		if card.Rank > 10 {
			score += 10 // J, Q, K
			continue
		} else if card.Rank == 1 {
			ace = true // Ace can be 1 or 11
		}
		score += card.Rank // Add the rank of the card
	}
	x := strconv.Itoa(score) // Start with the score
	if ace && score <= 11 {
		x += "/" + strconv.Itoa(score+10)
	}
	return x
}

func (gs GameState) Print() {
	// Print the player's hand
	
	print("Player Hand(s):\n")
	for i := 0; i < len(gs.PlayerHand); i++ {
		println(i+1, "("+printScore(gs.PlayerHand[i])+"):", PrintCards(gs.PlayerHand[i]))

	}
	println("")
	println("Dealer (" + printScore(gs.DealerHand) + "):")
	// only print the first card of the dealer's hand
	println(gs.DealerHand[0].String(), " ?")

}
