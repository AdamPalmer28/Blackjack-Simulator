package game

import (
	"fmt"
	"math/rand"
	"strconv"
)

// ============================================================================
// Deck structures

type Deck struct {
	Cards [52]Card

	Drawn int // number of cards drawn from the deck
	// i.e. next card to be drawn is at index Drawn
}

func newDeck() Deck {
	// make a new deck of cards
	var deck Deck


	for suit := 0; suit < 4; suit++ {
		for rank := 1; rank <= 13; rank++ {
			deck.Cards[suit*13+rank-1] = Card{Suit: suit, Rank: rank}
		}
	}

	// shuffle the deck
	deck.shuffle()

	return deck
}

// Draw a card from the deck
func (deck *Deck) Draw() Card {
	if deck.Drawn >= 52 {
		panic("Error: No more cards to draw from the deck.")
	}

	// get the next card to draw
	card := deck.Cards[deck.Drawn]
	deck.Drawn++

	return card
}
// shuffle the order of the cards in the deck
func (deck *Deck) shuffle() {

	// create a random sample of indices 0 to 51
	indices := rand.Perm(len(deck.Cards))

	// create a new deck to hold the shuffled cards
	var shuffledDeck [52]Card
	for i, index := range indices {
		shuffledDeck[i] = deck.Cards[index]
	}
	// assign the shuffled cards back to the deck
	deck.Cards = shuffledDeck
}

// Print human format deck - only for debugging purposes
func (deck Deck) Print() {
	for i, card := range deck.Cards {
		fmt.Printf("Card %d: %s\n", i+1, card)
	}
}

// ============================================================================
// Card structures

type Card struct {
	Suit int // 0: Hearts, 1: Diamonds, 2: Clubs, 3: Spades
	Rank int // 1: Ace, 2: Two, ..., 10: Ten, 11: Jack, 12: Queen, 13: King
}

func PrintCards(cards []Card) (string) {
	
	card_strings := ""
	for _, card := range cards {
		card_strings += card.String() + " "
	}

	return card_strings
}

// String returns a string representation of the card.
// It formats the card as "Rank of Suit", e.g., "Ace of Hearts", "10 of Diamonds".
func (c Card) String() string {

	// SUIT
	suitToString := func(suit int) string {
		switch suit {
		case 0:
			return "♥️"
		case 1:
			return "♦️"
		case 2:
			return "♣️"
		case 3:
			return "♠️"
		default:
			return "Unknown Suit"
		}
	}

	// RANK
	switch c.Rank {
	case 1:
		return "A" + suitToString(c.Suit)
	case 11:
		return "J" + suitToString(c.Suit)
	case 12:
		return "Q" + suitToString(c.Suit)
	case 13:
		return "K" + suitToString(c.Suit)
	default:
		return strconv.Itoa(c.Rank) + suitToString(c.Suit)
	}
}
