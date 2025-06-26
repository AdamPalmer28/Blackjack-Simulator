package game

import "strconv"

type GameState struct {

	// Game state
	Deck  Deck
	State []int // 0: player move, 1: dealer move, 2: player win, 3: dealer win, 4: draw, 5: player bust
	Value []int // value of each hand - used for splitting / doubling down
	
	HandToPlay int // player hand to play
	PlayerMoves []int // legal moves for the player (hit, double down, split) - 0b000: no moves, 0b001: hit, 0b010: double down, 0b100: split

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
		active_turn = false // player has stood, no more actions
	}
	// draw new card -- need to make this a function
	if playerMove == 0b001 { // hit

	} else if playerMove == 0b010 { // double down
		active_turn = false // player has doubled down, no more actions
		// double value
		// draw one card


	} else if playerMove == 0b100 { // split
		// Player splits their hand into two hands
		// This will be handled in the next turn
	} else {
		panic("Error: Invalid player move: " + strconv.Itoa(playerMove))
	}

	// ----------- AFTER ACTION -----------

	// if legal moves go back to user...
	if !active_turn {
		gs.HandToPlay++
		
		if gs.HandToPlay >= len(gs.PlayerHand) {
			// All player hands have been played, now it's the dealer's turn
			// TODO
			return
		}	
		
	}

	// determine legal moves for the next player to act
	gs.calcPlayMoves()

	return // return back to the game loop
}

func (gs *GameState) calcPlayMoves() {
	// Calculate the legal moves for the player based on their hand
	// This will be called after last players action

	// Legal moves are double, split, hit (111) note stand is always possible

	// Reset moves
	playerMove := gs.HandToPlay
	legalMoves := 0b001

	if gs.PlayerScore[playerMove] < 21 {
		// player can double
		legalMoves |= 0b010 
	}
	if gs.PlayerHand[playerMove][0].Rank == gs.PlayerHand[playerMove][1].Rank {
		// player can split
		legalMoves |= 0b100
	}

	// Add legal moves to the player's moves
	gs.PlayerMoves[playerMove] = legalMoves
}

func (gs *GameState) UpdatePlayerState() {
	// player score, legal moves
	ind := gs.HandToPlay

	// Calculate the player's score
	gs.PlayerScore[ind] = calculateScore(gs.PlayerHand[ind])
	// Check if player has an Ace
	gs.playerAce[ind] = false
	for _, card := range gs.PlayerHand[ind] {
		if card.Rank == 1 {
			gs.playerAce[ind] = true
			break
		}
	}

	// Calculate legal moves for the player
	gs.calcPlayMoves()
}

func StartGame() GameState {

	// Initialize a new game state
	gs := GameState{
		Deck: newDeck(),
		// State of play
		HandToPlay: 0,
		PlayerMoves: make([]int, 0), // legal moves (hit, double down, split)

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
	// calculate dealers state


	return gs
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
