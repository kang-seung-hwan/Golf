/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a reservation
type SmartContract struct {
	contractapi.Contract
}

// GameInfo is the structure that informs the game information
type GameInfo struct {
	GroundID   string `json:"groundID"`
	User1ID    string `json:"user1ID"`
	User2ID    string `json:"user2ID"`
	User3ID    string `json:"user3ID"`
	User4ID    string `json:"user4ID"`
	GameNumber string `json:"gameNumber"`
	GameCode   string `json:"gameCode"`
	IsReady    bool `json:"isReady"`
}

// GameNumberKey is the struct containing a gameNumber key and index
type GameNumberKey struct {
	Key string
	Idx int
}

// HoleScore is the struct that informs hole number, each user score and consensus result.
// All user achive consensus, Validated is true
type HoleScore struct {
	HoleNumber string `json:"holeNumber"`
	User1Score string `json:"user1Score"`
	User2Score string `json:"user2Score"`
	User3Score string `json:"user3Score"`
	User4Score string `json:"user4Score"`
	Validated  bool   `json:"validated"`
}

// Agreement is the struct that informs each user's agreement the score
type Agreement struct {
	HoleNumber string `json:"holeNumber"`
	User1Agree string `json:"user1Agree"`
	User2Agree string `json:"user2Agree"`
	User3Agree string `json:"user3Agree"`
	User4Agree string `json:"user4Agree"`
}

// InitLedger is the init function
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("initLedger")
	return nil
}

/*
// StartGame is the invoke function that starts the gaim with user1, 2, 3, and 4
// params - ground ID, users ID(4)
func (s *SmartContract) StartGame(ctx contractapi.TransactionContextInterface, groundID, user1, user2, user3, user4 string, isReady bool) error {
	fmt.Println("Start Game")

	// get the unique game number
	var gameNumberKey *GameNumberKey
	gameNumberKey = s.GenerateKey(ctx, "latestKey")
	keyidx := strconv.Itoa(gameNumberKey.Idx)
	fmt.Println("Key : " + gameNumberKey.Key + ", Idx : " + keyidx)

	var gameNumber = gameNumberKey.Key + keyidx
	fmt.Println("gameNumberKey is " + gameNumber)

	// create composite key for the gameInfo
	GameCompositeKey, _ := ctx.GetStub().CreateCompositeKey("game", []string{groundID, gameNumber})

	gameInfo := GameInfo{
		GroundID: groundId,
		GameNumber: gameNumber,
		User1ID: user1,
		User2ID: user2,
		User3ID: user3,
		User4ID: user4,
		IsReady: isReady
	}
	
	if gameInfo.User4ID != user4 {
		gameInfo.IsReady = false
	} else {
		gameInfo.IsReady = true
	}

	gameInfoAsBytes, err := json.Marshal(gameInfo)
	if err != nil {
		return fmt.Errorf("gameInfo Marshal Error: %s", err.Error())
	}
	

	// update game number
	gameNumberKeyAsBytes, _ := json.Marshal(gameNumberKey)
	ctx.GetStub().PutState("latestKey", gameNumberKeyAsBytes)

	
	// update gameInfo
	return ctx.GetStub().PutState(GameCompositeKey, gameInfoAsBytes)
}
*/

func (s *SmartContract) StartGame(ctx contractapi.TransactionContextInterface, groundID, user, userNumber, gameCode string) error {
	fmt.Println("Start Game")

	// get the unique game number
	var gameNumberKey *GameNumberKey
	gameNumberKey = s.GenerateKey(ctx, "latestKey")
	keyidx := strconv.Itoa(gameNumberKey.Idx)
	fmt.Println("Key : " + gameNumberKey.Key + ", Idx : " + keyidx)

	var gameNumber = gameNumberKey.Key + keyidx
	fmt.Println("gameNumberKey is " + gameNumber)

	// create composite key for the gameInfo
	GameCompositeKey, _ := ctx.GetStub().CreateCompositeKey("game", []string{groundID, gameNumber})
	gameInfo, err := s.QueryGameInfo(ctx, groundID, gameNumber)
	if gameInfo == nil {
		gameInfo = &GameInfo{}
	}
	gameInfo.GroundID = groundID
	gameInfo.GameNumber = gameNumber
	if userNumber == "1"{
		gameInfo.User1ID = user
	} else if userNumber == "2"{
		gameInfo.User2ID = user
	} else if userNumber == "3" {
		gameInfo.User3ID = user
	} else if userNumber == "4" {
		gameInfo.User4ID = user
	}
	gameInfo.GameCode = gameCode

	if gameInfo.User4ID != user {
		gameInfo.IsReady = false
	} else {
		gameInfo.IsReady = true
	}

	gameInfoAsBytes, err := json.Marshal(gameInfo)
	if err != nil {
		return fmt.Errorf("gameInfo Marshal Error: %s", err.Error())
	}
	

	// update next game number
	if gameInfo.GameCode != gameCode  {
	gameNumberKeyAsBytes, _ := json.Marshal(gameNumberKey)
	ctx.GetStub().PutState("latestKey", gameNumberKeyAsBytes)
	}

	
	// update gameInfo
	return ctx.GetStub().PutState(GameCompositeKey, gameInfoAsBytes)
}


// QueryGameInfo returns the gameInfo stored in the world state with given IDs and gameNumber
// params - ground ID, user's ID, and unique game number
// returns the GameInfo
func (s *SmartContract) QueryGameInfo(ctx contractapi.TransactionContextInterface, groundID, gameNumber string) (*GameInfo, error) {
	fmt.Println("QueryGameInfo")

	// get the Composite key for the gameInfo
	GameCompositeKey, _ := ctx.GetStub().CreateCompositeKey("game", []string{groundID, gameNumber})
	gameInfoAsBytes, err := ctx.GetStub().GetState(GameCompositeKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if gameInfoAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", gameNumber)
	}

	gameInfo := new(GameInfo)
	_ = json.Unmarshal(gameInfoAsBytes, gameInfo)

	return gameInfo, nil
}

// SetScore is the invoke function that sets the user's score with given hole and user number
// params - unique game number, hole's number, user's Number(ordered 1,2,3,4) and score
func (s *SmartContract) SetScore(ctx contractapi.TransactionContextInterface, gameNumber, holeNumber, userNumber, score string) error {
	fmt.Println("SetScore")

	// create composite key for the hole score
	HoleScoreCompositeKey, _ := ctx.GetStub().CreateCompositeKey("holeScore", []string{gameNumber, holeNumber})

	// get the holeScore
	holeScore, err := s.QueryScore(ctx, gameNumber, holeNumber)
	// initialize the hole score
	if holeScore == nil {
		holeScore = &HoleScore{
			HoleNumber: holeNumber,
			Validated:  false,
		}
	}
	// update holeScore with given userNumber and score
	if userNumber == "1" {
		holeScore.User1Score = score
	} else if userNumber == "2" {
		holeScore.User2Score = score
	} else if userNumber == "3" {
		holeScore.User3Score = score
	} else if userNumber == "4" {
		holeScore.User4Score = score
	}

	holeScoreAsBytes, err := json.Marshal(holeScore)
	if err != nil {
		return fmt.Errorf("holseScore Marshal Error: %s", err.Error())
	}

	return ctx.GetStub().PutState(HoleScoreCompositeKey, holeScoreAsBytes)
}

// QueryScore is the query function that returns the HoleScore
// params - gameNumber, holeNumber
// returns the HoleScore
func (s *SmartContract) QueryScore(ctx contractapi.TransactionContextInterface, gameNumber, holeNumber string) (*HoleScore, error) {
	fmt.Println("QueryScore")

	// create composite key for the holeScore
	HoleScoreCompositeKey, _ := ctx.GetStub().CreateCompositeKey("holeScore", []string{gameNumber, holeNumber})
	holeScoreAsBytes, err := ctx.GetStub().GetState(HoleScoreCompositeKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if holeScoreAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", gameNumber)
	}

	holeScore := new(HoleScore)
	_ = json.Unmarshal(holeScoreAsBytes, holeScore)

	return holeScore, nil
}

// AgreeScore is the invoke function that updates the Agreement with given userNumber and agreement
// params - gameNumber, holeNumber, userNumber, and user's agreement status
func (s *SmartContract) AgreeScore(ctx contractapi.TransactionContextInterface, gameNumber, holeNumber, userNumber, isAgreed string) error {
	fmt.Println("AgreeScore")

	// create composite key for the agreement
	AgreementCompositeKey, _ := ctx.GetStub().CreateCompositeKey("agreement", []string{gameNumber, holeNumber})

	// get the Agreement
	agreement, err := s.QueryAgreement(ctx, gameNumber, holeNumber)
	// initialize the Agreement
	if agreement == nil {
		agreement = &Agreement{}
	}
	// update the Agreement with given user number and agreement status
	if userNumber == "1" {
		agreement.User1Agree = isAgreed
	} else if userNumber == "2" {
		agreement.User2Agree = isAgreed
	} else if userNumber == "3" {
		agreement.User3Agree = isAgreed
	} else if userNumber == "4" {
		agreement.User4Agree = isAgreed
	}

	// when all user agree the score, call the function ValidateScore.
	if agreement.User1Agree == "agree" && agreement.User2Agree == "agree" && agreement.User3Agree == "agree" && agreement.User4Agree == "agree" {
		fmt.Println("All user agrees the score")
		s.ValidateScore(ctx, gameNumber, holeNumber)
	}
	agreementAsBytes, err := json.Marshal(agreement)
	if err != nil {
		return fmt.Errorf("holseScore Marshal Error: %s", err.Error())
	}

	return ctx.GetStub().PutState(AgreementCompositeKey, agreementAsBytes)
}

// QueryAgreement query function that returns the Agreement
// params - gameNumber, holeNumber
// returns the Agreement
func (s *SmartContract) QueryAgreement(ctx contractapi.TransactionContextInterface, gameNumber, holeNumber string) (*Agreement, error) {
	fmt.Println("QueryAgreement")

	// create composite key for the Agreement
	AgreementCompositeKey, _ := ctx.GetStub().CreateCompositeKey("agreement", []string{gameNumber, holeNumber})
	agreementAsBytes, err := ctx.GetStub().GetState(AgreementCompositeKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if agreementAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", gameNumber)
	}

	agreement := new(Agreement)
	_ = json.Unmarshal(agreementAsBytes, agreement)

	return agreement, nil
}

// ValidateScore is the invoke function that update the HoleScore's Validated
// params - gameNumber, holeNumber
func (s *SmartContract) ValidateScore(ctx contractapi.TransactionContextInterface, gameNumber, holeNumber string) error {
	fmt.Println("ValidateScore")

	// All user must be agree the score
	agreement, err := s.QueryAgreement(ctx, gameNumber, holeNumber)
	if agreement == nil {
		return fmt.Errorf("all must be agree the score, %s", err.Error())
	}

	if agreement.User1Agree == "agree" && agreement.User2Agree == "agree" && agreement.User3Agree == "agree" && agreement.User4Agree == "agree" {
		fmt.Println("All user agrees the score")
		s.ValidateScore(ctx, gameNumber, holeNumber)
	}

	// create composite key for the holeScore
	HoleScoreCompositeKey, _ := ctx.GetStub().CreateCompositeKey("holeScore", []string{gameNumber, holeNumber})
	holeScore, err := s.QueryScore(ctx, gameNumber, holeNumber)
	if holeScore == nil {
		return nil
	}

	// update Validated
	holeScore.Validated = true
	holeScoreAsBytes, err := json.Marshal(holeScore)
	if err != nil {
		return fmt.Errorf("holseScore Marshal Error: %s", err.Error())
	}

	return ctx.GetStub().PutState(HoleScoreCompositeKey, holeScoreAsBytes)
}

// QueryTotalGameScore is the query function that returns the total game score
// params - gameNumber
func (s *SmartContract) QueryTotalGameScore(ctx contractapi.TransactionContextInterface, gameNumber string) ([]*HoleScore, error) {
	fmt.Println("QueryTotalGameScore")
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("holeScore", []string{gameNumber})

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var holesScore []*HoleScore

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var holeScore HoleScore

		_ = json.Unmarshal(queryResponse.Value, &holeScore)
		if holeScore.Validated {
			holesScore = append(holesScore, &holeScore)
		}
	}

	return holesScore, nil
}

// GenerateKey is the function that generate unique game key
func (s *SmartContract) GenerateKey(ctx contractapi.TransactionContextInterface, key string) *GameNumberKey {
	fmt.Println("GenerateKey")
	var isFirst bool = false

	gameNumberKeyAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	var gameNumberKey *GameNumberKey
	if gameNumberKeyAsBytes == nil {
		isFirst = true
		gameNumberKey = &GameNumberKey{
			Key: "GAME",
			Idx: 1,
		}
	}

	var tempIdx string
	if !isFirst {
		err = json.Unmarshal(gameNumberKeyAsBytes, &gameNumberKey)
		if err != nil {
			fmt.Println(err.Error())
		}
		tempIdx = strconv.Itoa(gameNumberKey.Idx)
		gameNumberKey.Idx = gameNumberKey.Idx + 1
	}

	fmt.Println("Last GameNumber is " + gameNumberKey.Key + " : " + tempIdx)

	return gameNumberKey
}

// main function
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}

}