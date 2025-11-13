# Blackjack Strategy Analysis

This directory contains tools for analyzing the blackjack simulation data and generating optimal strategy charts.

## Setup

Install the required Python packages:

```bash
pip install -r requirements.txt
```

## Running the Analysis

Open the Jupyter notebook:

```bash
jupyter notebook blackjack_analysis.ipynb
```

Or run all cells to generate the analysis outputs.

## Outputs

The notebook generates the following files:

### CSV Files
- **blackjack_analysis_data.csv**: Complete dataset with all expected values for each action
- **blackjack_summary.csv**: Summary dataset showing the best action for each game state

### Strategy Heatmaps
Three heatmap images showing optimal play for each hand type:
1. **strategy_heatmap_0_normal_hand_(hard).png**: Hard hands (no ace)
2. **strategy_heatmap_1_soft_hand_(has_ace).png**: Soft hands (with ace)
3. **strategy_heatmap_2_pair_(split_available).png**: Pairs (split available)

### Additional Visualizations
- **ev_distributions.png**: Expected value distributions by hand category
- **action_distribution.png**: Action frequency by hand category

## Understanding the Heatmaps

- **Y-axis**: Player's score (2-20)
- **X-axis**: Dealer's shown card (1-10, where 1=Ace)
- **Cell colors**: Best action to take
  - Red (H): Hit
  - Teal (S): Stand
  - Yellow (D): Double down
  - Light green (P): Split

## Data Structure

The simulation data (`../bj_sim_data.json`) has the following structure:
- **Dealer Shown Score**: 1-10 (1=Ace, 10=10/Face cards)
- **Player Score**: 2-20
- **Hand Categories**:
  - 0: Normal/Hard hand (no ace)
  - 1: Soft hand (has ace, score ≤ 11)
  - 2: Split available (pair)
- **Actions**:
  - 0: Hit
  - 1: Stand
  - 2: Double down
  - 3: Split

Each game state stores the expected value and number of trials for each possible action.
