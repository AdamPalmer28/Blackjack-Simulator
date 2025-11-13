# Blackjack Simulator

A computational exercise implementing a recursive blackjack simulator in Go to explore all possible game states and determine optimal playing strategies through statistical analysis.

## Quick Start

### Running Simulations

```bash
go run .
```

This runs simulations and saves results to `bj_sim_data.json`.

### Analyzing Results

Navigate to the `analysis/` directory and run the Jupyter notebook to generate strategy charts:

```bash
cd analysis
pip install -r requirements.txt
jupyter notebook blackjack_analysis.ipynb
```

## Analysis and Findings

The analysis produces optimal strategy heatmaps for three types of hands:
1. **Hard Hands** (no ace)
2. **Soft Hands** (with ace)
3. **Pairs** (split available)

### Key Findings

Based on the simulation data:

#### Overall Statistics
- **Most common optimal action**: Hit (37.9% of situations)
- **Average expected value**: +0.015 (slight player advantage with optimal play)
- **Best possible situation**: Splitting 2s against dealer's 6 (EV: +1.15)
- **Worst situation**: Hard 16 vs dealer 10 (EV: -0.79)

#### Strategy by Hand Type

**Hard Hands (No Ace)**
- Average EV: -0.11 (player disadvantage)
- Strategy: Hit aggressively on low totals (≤11), stand on 17+
- Most challenging: Hard 12-16 vs dealer's high cards (7-10, Ace)

**Soft Hands (With Ace)**
- Average EV: +0.13 (player advantage)
- Strategy: More aggressive with doubling, especially vs weak dealer cards (4-6)
- Key advantage: Flexibility of the ace allows risk-taking

**Pairs (Split Available)**
- Average EV: +0.15 (best overall category)
- Strategy: Split aggressively (59% of situations)
- Always split: Aces and 8s
- Never split: 10s and 5s
- Best opportunities: Low pairs vs dealer's weak cards

#### Dealer Impact

Dealer's shown card significantly affects player advantage:
- **Worst for player**: Dealer shows 10 (avg EV: -0.38)
- **Best for player**: Dealer shows 6 (avg EV: +0.18)
- **Weak dealer cards** (4-6): Player should be more aggressive with doubles and splits
- **Strong dealer cards** (9, 10, Ace): Conservative play with focus on not busting

### Performance Insights

The recursive simulation explores all possible game paths to determine mathematically optimal strategies. The results align closely with traditional basic strategy charts, validating the simulation's accuracy while providing precise expected value calculations for each decision point.

## Project Structure

```
blackjack/
├── main.go              # Entry point with CLI and simulation modes
├── go.mod              # Go module configuration
├── game/               # Core game logic
│   ├── game.go         # GameState and mechanics
│   └── cards.go        # Card and Deck structures
├── sim/                # Simulation engine
│   ├── sim.go          # Recursive exploration logic
│   └── data.go         # Data collection and persistence
├── analysis/           # Python analysis tools
│   ├── blackjack_analysis.ipynb  # Jupyter notebook
│   ├── requirements.txt          # Python dependencies
│   └── README.md                 # Analysis documentation
└── bj_sim_data.json    # Simulation results
```

## Technical Implementation

### Core Features

- **Recursive State Exploration**: Explores all possible game outcomes from each decision point
- **Deep Copy System**: Ensures independent simulation branches
- **Bitfield Move System**: Efficient legal move representation
- **Multi-Hand Support**: Handles complex splitting scenarios
- **Statistical Data Collection**: Records expected values and trial counts

### Game Mechanics

- Standard casino rules (dealer hits on 16, stands on 17)
- Actions: Hit, Stand, Double Down, Split
- Proper ace handling (soft/hard conversion)
- Bet value tracking for expected value calculation

### Data Structure

Simulation results are stored hierarchically:
```
[Dealer Score][Player Score][Hand Category][Action] → {Expected Value, Trials}
```

- **Dealer Score**: 1-10 (1=Ace, 10=10/Face)
- **Player Score**: 2-20
- **Hand Category**: 0=Hard, 1=Soft, 2=Pair
- **Action**: 0=Hit, 1=Stand, 2=Double, 3=Split

## Dependencies

- **Go 1.x+**: Core simulator (no external dependencies)
- **Python 3.7+**: Analysis tools
  - jupyter
  - pandas
  - matplotlib
  - seaborn
  - numpy

## Academic Context

While blackjack is a solved game with well-documented optimal strategies, this project serves as a computational exercise in:
- Recursive algorithm implementation
- State space exploration
- Statistical data collection and analysis
- Go programming language proficiency
- Data visualization techniques

The simulator provides empirical validation of theoretical blackjack strategies through comprehensive state space exploration and statistical analysis.
