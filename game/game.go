package game

type GameState struct {
	Deck Deck
	State uint // 0: player move, 1: dealer move, 2: player win, 3: dealer win, 4: draw, 5: player bust

	// Player's hand
	PlayerHand  [][]Card // allow for splitting hands
	PlayerScore int
	playerAce   bool // true if player has an Ace in their hand
	Moves uint // legal moves (hit, double down, split)

	// Dealer's hand
	DealerHand       []Card
	dealerShownScore int  // score of the dealer's shown card
	dealerShownAce   bool // true if dealer's shown card is an Ace

	dealerAce   bool // true if dealer has an Ace in their hand
	DealerScore int
}

// ----------------------------------------------------------------------------
// ACTIONS

// Player can hit, stand, double down, or split


// ----------------------------------------------------------------------------
// GAME LOGIC

func (gs *GameState) nextMove(playerMove bool) {

	// Calculate scores

}

func (gs *GameState) calculateScore() {

func StartGame() GameState {

	// Initialize a new game state
	gs := GameState{
		Deck: newDeck(),

		PlayerHand:  make([][]Card, 0), // Start with no player hands
		PlayerScore: 0,
		playerAce:   false,

		DealerHand:       make([]Card, 0),
		DealerScore:      0,
		dealerShownScore: 0,
		dealerAce:        false,
		dealerShownAce:   false,
	}

	// Deal initial cards to player and dealer
	gs.dealInitialCards()

	gs.NextMove()
	return gs
}

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
}

// ============================================================================
// HELPER FUNCTIONS

func printScore([]Card) string {

	// ! not sure how to do this???

	// if has Ace then appear 1/11..
	// if has 10, J, Q, K then appear 10
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
