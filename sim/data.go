/*
How to record the blackjack data for the simulation

first layer key - dealer score
second layer key - player score, record soft playerscores seperately for ace boolean (only for score <= 10)-
third layer key - options [hit, stand, double down, split]
values - [expected value, number of trials]
*/
package sim

import (
	"encoding/json"
	"fmt"
	"os"
)

// individual simulation data piece
type SimEvalData struct {
	DealerScore  int
	PlayerScores int
	PlayerAce bool
	ChoosenAction int // 0: hit, 1: stand, 2: double down, 3: split
	Value int // resulting value of the action
}



// Complete structure for simulation data
// made from many SimData
type SimData struct {
	// DealerScore  int // dealer score
	// PlayerScores int // player score
	// PlayerAce bool // true if player has an Ace in their hand
	// ChosenAction int // 0: hit, 1: stand, 2: double down, 3: split
	ExpectedValue int // resulting value of the action
	Trials int // number of trials for this data set
}

// map of SimData dealer score, player score, player ace boolean, chosen action
type SimDataMap map[int]map[int]map[int]map[int]SimData


func CreateSimDataStructure() SimDataMap {
	ds := make(map[int]map[int]map[int]map[int]SimData)

	for i := 2; i <= 20; i++ { // dealer score
		for j := 2; j <= 20; j++ { // player score

			limit := 1
			if j < 12 { // ace doesn't matter for player scores above 11
				limit = 2
			}
			for k := 0; k < limit; k++ { // player ace boolean
				for l := 0; l < 4; l++ { // actions [hit, stand, double down, split]
					// Initialize the data structure with zero values
					simData := SimData{
						ExpectedValue: 0,
						Trials:       0,
					}
					// Add simData to the simulation data structure
					//fmt.Println(i, j, k == 1, l)
					if _, ok := ds[i]; !ok {
						ds[i] = make(map[int]map[int]map[int]SimData)
					}
					if _, ok := ds[i][j]; !ok {
						ds[i][j] = make(map[int]map[int]SimData)
					}
					if _, ok := ds[i][j][k]; !ok {
						ds[i][j][k] = make(map[int]SimData)
					}
					ds[i][j][k][l] = simData
				}
			}
			
		}
	}

	return ds
}

// SimDataMap to JSON "bj_sim_data.json"
func (sdm SimDataMap) ToJSON() ([]byte, error) {
	
	data, err := json.Marshal(sdm)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile("bj_sim_data.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil, err
	}

	return data, nil
}
