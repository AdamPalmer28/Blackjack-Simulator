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

// Complete structure for simulation data
// made from many SimData
type SimData struct {
	// DealerScore  int // dealer shown score
	// PlayerScores int // player score
	// PlayerOptions int // 0: no ace, 1: has ace, 2: split
	// ChosenAction int // 0: hit, 1: stand, 2: double down, 3: split
	ExpectedValue int // resulting value of the action
	Trials int // number of trials for this data set
}

// map of SimData dealer score, player score, player ace boolean, chosen action
type SimDataMap map[int]map[int]map[int]map[int]SimData


func CreateSimDataStructure() SimDataMap {
	ds := make(map[int]map[int]map[int]map[int]SimData)

	for i := 1; i <= 10; i++ { // dealer shown score
		for j := 2; j <= 20; j++ { // player score

			loopList := []int{0}
			if j < 12 && j > 2{
				loopList = append(loopList, 1)
			}
			if j%2 == 0 {
				loopList = append(loopList, 2)
			}
			for _, k := range loopList { // player ace boolean

				player_options := 3
				if k == 2 {
					player_options = 4
				}

				for l := 0; l < player_options; l++ { // actions [hit, stand, double down, split]
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

func (sdm SimDataMap) AddData(data SimState) {
	// add data to the simulation data structure
	for _, d := range data.SimEvalData {
		dealerScore := d.DealerScore
		playerScore := d.PlayerScores
		playerHandCat := d.PlayerHandCats
		chosenAction := d.ChoosenAction

		// update the expected value and trials
		simData := sdm[dealerScore][playerScore][playerHandCat][chosenAction]
		simData.ExpectedValue = (simData.ExpectedValue*simData.Trials + d.Value) / (simData.Trials + 1)
		simData.Trials++
		sdm[dealerScore][playerScore][playerHandCat][chosenAction] = simData
	}
}


// ----------------------------------------------------------------------------

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

// load data from JSON file to SimDataMap
func LoadFromJSON(filename string) (SimDataMap, error) {
	file, err := os.ReadFile(filename)
	if err != nil {	
		return nil, err
	}
	var sdm SimDataMap
	err = json.Unmarshal(file, &sdm)
	if err != nil {
		return nil, err
	}
	return sdm, nil
}
