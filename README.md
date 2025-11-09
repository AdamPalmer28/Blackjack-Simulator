# Blackjack Simulator Documentation

While this Blackjack is a solved game with well-known optimal strategies, this simulator is a computational exercise for myself to implement a complex recursive simulation in Go, exploring all possible game states and collecting statistical data.

## Overview

This Blackjack Simulator is a Go-based application designed to simulate and analyse blackjack gameplay through recursive exploration of all possible game states. The simulator can both run interactive CLI games and perform automated simulations to collect statistical data about optimal blackjack strategies.

## Project Structure

```
blackjack/
├── main.go           # Entry point with CLI interface
├── go.mod           # Go module configuration
├── game/            # Core game logic
│   ├── game.go      # GameState and game mechanics
│   └── cards.go     # Card and Deck structures
├── sim/             # Simulation engine
│   ├── sim.go       # Simulation logic and recursive exploration
│   └── data.go      # Data collection and JSON persistence
└── bj_sim_data.json # Simulation results storage
```

## Core Components

### 1. Game Engine (`game/`)

#### Card System (`cards.go`)

- **`Card`**: Represents a playing card with suit (0-3) and rank (1-13)
  - Suits: Hearts (♥️), Diamonds (♦️), Clubs (♣️), Spades (♠️)
  - Ranks: 1=Ace, 2-10=Number cards, 11=Jack, 12=Queen, 13=King
- **`Deck`**: Contains 52 cards with shuffling and drawing capabilities
  - Tracks drawn cards to prevent reuse
  - Supports deep copying for simulation branching

#### Game State (`game.go`)

The `GameState` struct manages the complete state of a blackjack game:

```go
type GameState struct {
	// Game state
	Deck  Deck
	State []int // 0: active game, 1: player win, 2: dealer win, 3: draw, 4: player bust

	HandToPlay int // player hand to play
	PlayerMoves []int // legal moves for the player (hit, double down, split) - 0b000: no moves, 0b001: hit, 0b010: double down, 0b100: split
	HandValues []int  // value of each hand - used for splitting / doubling down

	// Player's hand
	PlayerHand  [][]Card // allow for splitting hands
	PlayerScore []int
	playerAce   []bool // true if player has an Ace in their hand
	Moves       []int  // legal moves (hit, double down, split)

	// Dealer's hand
	DealerHand       []Card
	DealerShownScore int  // score of the dealer's shown card
	dealerShownAce   bool // true if dealer's shown card is an Ace

	dealerAce   bool // true if dealer has an Ace in their hand
	DealerScore int8
}
```

#### Key Game Mechanics

**Legal Moves (Bitfield System)**:

- `0b001` (1): Hit - Draw another card
- `0b010` (2): Double Down - Double bet, draw one card, end turn
- `0b100` (4): Split - Split pair into two hands
- `0b000` (0): Stand - End turn (always available)

**Action Processing (`ActionCalc`)**:

1. Validates and executes player action
2. Updates scores and game state
3. Handles splitting logic (creates new hands)
4. Manages turn progression
5. Triggers endgame when all hands complete

**Endgame Logic**:

1. Dealer draws until score ≥ 17 (or soft 17+ with ace)
2. Calculates final outcomes:
   - `1`: Player Win
   - `2`: Dealer Win
   - `3`: Draw/Push
   - `4`: Player Bust

**Scoring System**:

- Number cards: Face value
- Face cards (J,Q,K): 10 points
- Aces: 1 point (soft/hard ace logic handled in endgame)

### 2. Simulation Engine (`sim/`)

#### Core Simulation (`sim.go`)

The simulation uses recursive tree exploration to analyze all possible game outcomes:

**`SimEvalData`** - Records each decision point:

```go
type SimEvalData struct {
    DealerStart    int  // Dealer's shown card
    DealerScore    int  // Final dealer score
    PlayerScores   int  // Player hand score
    PlayerHandCats int  // Hand category (0=no ace, 1=has ace, 2=split)
    ChoosenAction  int  // Action taken (0-3)
    Value         int  // Outcome value (+/- bet amount)
    Depth         int  // Tree depth (for debugging)
}
```

**Recursive Exploration (`node_explore`)**:

1. **Base Case**: When all hands are played, record final outcomes
2. **Recursive Case**: For each legal action:
   - Create deep copy of game state
   - Execute action via `ActionCalc`
   - Recursively explore resulting state
   - Record action and outcome data

**Action Evaluation Strategy**:

- **Early Game (Score ≤ 14)**: Explores all possible actions
- **Late Game (Score > 14)**: Uses strategic action selection
- **Value Calculation**: Returns average expected value across explored branches

#### Data Collection (`data.go`)

**`SimDataMap`** - Hierarchical data structure:

```go
// Structure: [DealerScore][PlayerScore][HandCategory][Action] -> SimData
type SimDataMap map[int]map[int]map[int]map[int]SimData
```

**Data Dimensions**:

- **Dealer Score**: 1-10 (shown card)
- **Player Score**: 2-20
- **Hand Categories**:
  - 0: Hard hand (no ace)
  - 1: Soft hand (has ace, score ≤ 11)
  - 2: Split opportunity (pair)
- **Actions**: 0=Hit, 1=Stand, 2=Double, 3=Split

**Persistence**:

- JSON serialization to `bj_sim_data.json`
- Accumulates expected values and trial counts
- Supports incremental data collection

### 3. User Interface (`main.go`)

#### CLI Game Mode

Interactive blackjack with:

- Visual card display with Unicode suits
- Legal move validation
- Real-time score calculation
- Multi-hand support (splitting)
- End-game result display

#### Simulation Mode

Automated analysis featuring:

- Configurable simulation count
- Real-time progress reporting
- Debug output showing decision trees
- Statistical data collection

## Key Features

### Deep Copy System

Ensures simulation branches don't interfere:

- **GameState**: Copies all slices and nested structures
- **Deck**: Preserves card order and draw state
- Critical for accurate tree exploration

### Bit-field Move System

Efficient legal move representation:

```go
// Example: Can Hit and Double Down
moves := 0b011  // Binary: Hit(1) + Double(2)

// Check if action is legal
if moves & 0b010 != 0 {  // Can double down
    // Execute double down
}
```

### Multi-Hand Support

Handles complex splitting scenarios:

- Dynamic hand array expansion
- Independent score/state tracking per hand
- Sequential turn processing
- Proper bet value management

### Statistical Analysis

Comprehensive data collection:

- Expected value calculation
- Action frequency analysis
- Outcome distribution tracking
- Strategy optimization data

## Usage Examples

### Running Simulations

```bash
go run .
# Runs 10 hands by default, saves to bj_sim_data.json
```

### Interactive Play

```go
// In main.go, uncomment CLI section for interactive mode
blackjackCLI()
```

### Analyzing Results

The simulation produces data showing:

- Optimal actions for each game state
- Expected values for different strategies
- Win/loss probabilities
- Impact of different hand compositions

## Technical Implementation Details

### Ace Handling

- **During Play**: Aces count as 1 point
- **Endgame**: Automatically converts to 11 if beneficial
- **Soft Hands**: Tracked separately in data collection

### Splitting Logic

- Creates two new hands from split pair
- Each hand gets one new card immediately
- Independent action sequences for each hand
- Proper bet multiplication

### Dealer AI

- Standard casino rules: Hit on 16, stand on 17
- Soft 17 handling (hits on soft 17)
- Automatic play after all player hands complete

### Memory Management

- Deep copying prevents state corruption
- Efficient slice operations
- Garbage collection friendly design

## Future Enhancements

Potential improvements could include:

- Machine learning integration for strategy optimization
- Parallel simulation processing
- Advanced statistical analysis tools
- Web-based visualization interface
- Card counting simulation capabilities

## Dependencies

- **Go 1.x+**: Core language runtime
- **Standard Library Only**: No external dependencies
  - `encoding/json`: Data persistence
  - `math/rand`: Card shuffling
  - `bufio`: CLI input handling

This simulator provides a robust foundation for blackjack analysis, combining accurate game simulation with comprehensive data collection capabilities.
