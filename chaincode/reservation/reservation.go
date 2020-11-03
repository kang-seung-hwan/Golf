/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"math/rand"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a reservation
type SmartContract struct {
	contractapi.Contract
}

// Ground information
type Ground struct {
	GroundID           string `json:"groundID"`
	GroundName         string `json:"groundName"`
	AvailableTimeStart uint   `json:"availableTimeStart"`
	AvailableTimeEnd   uint   `json:"availableTimeEnd"`
	TotalHole          uint   `json:"totalHole"`
	// HolesInfo          map[uint]*HoleInfo `json:"holesInfo"`
}

// HoleInfo is the struct that informs number of par and difficulty
// type HoleInfo struct {
// 	ParNumber  uint   `json:"parNumber"`
// 	Difficulty string `json:"difficulty"`
// }

// Reservation is the sturct that desribes the reservation information.
type Reservation struct {
	GroundID          string    `json:"groundID"`
	UserID            string    `json:"userID"`
	Begin             time.Time `json:"begin"`
	End               time.Time `json:"end"`
	ReservationNumber string    `json:"reservationNumber"`
	GameCode          int		`json:"gameCode"`
}

// ReservationKey is the struct containing a reservation key and index
type ReservationKey struct {
	Key string
	Idx int
}

// InitLedger adds a base set of grounds to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	ground := Ground{
		GroundID:           "Ground01",
		GroundName:         "수성",
		AvailableTimeStart: 9,
		AvailableTimeEnd:   18,
		TotalHole:          34,
		// HolesInfo:          make(map[uint]*HoleInfo),
	}
	groundCompositeKey, _ := ctx.GetStub().CreateCompositeKey("ground", []string{"Ground01"})

	groundAsBytes, _ := json.Marshal(ground)
	err := ctx.GetStub().PutState(groundCompositeKey, groundAsBytes)
	if err != nil {
		return fmt.Errorf("Failed to put the world state. %s", err.Error())
	}

	return nil
}

// CreateGround adds a new ground to the world state with given details
// params - groundID, grond name, start time, end time, and total hole number
func (s *SmartContract) CreateGround(ctx contractapi.TransactionContextInterface, groundID string, name string, startTime uint, endTime uint, totalHole uint) error {
	fmt.Println("CreateGround called")

	// create composite key for the ground
	groundCompositeKey, _ := ctx.GetStub().CreateCompositeKey("ground", []string{groundID})
	ground := Ground{
		GroundID:           groundID,
		GroundName:         name,
		AvailableTimeStart: startTime,
		AvailableTimeEnd:   endTime,
		TotalHole:          totalHole,
		// HolesInfo:          make(map[uint]*HoleInfo),
	}

	groundAsBytes, _ := json.Marshal(ground)

	return ctx.GetStub().PutState(groundCompositeKey, groundAsBytes)
}

// QueryGround returns the ground stored in the world state with given groundID
// params - groundID
// returns the Ground
func (s *SmartContract) QueryGround(ctx contractapi.TransactionContextInterface, groundID string) (*Ground, error) {
	groundCompositeKey, _ := ctx.GetStub().CreateCompositeKey("ground", []string{groundID})
	groundAsBytes, err := ctx.GetStub().GetState(groundCompositeKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if groundAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", groundID)
	}

	ground := new(Ground)
	_ = json.Unmarshal(groundAsBytes, ground)

	return ground, nil
}

// QueryAllGround returns all grounds found in world state
// returns the array of Ground
func (s *SmartContract) QueryAllGround(ctx contractapi.TransactionContextInterface) ([]*Ground, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("ground", []string{})

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var grounds []*Ground

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var ground Ground

		_ = json.Unmarshal(queryResponse.Value, &ground)

		grounds = append(grounds, &ground)
	}

	return grounds, nil
}

// parseTime is the parsing funciton
// params - string of time
// returns the time object
func parseTime(timeString string) time.Time {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		fmt.Printf("Parsing Time Erorr!!!!: %s", err.Error())
	}
	return t
}

func createRandomCode() int {
	s1:=rand.NewSource(time.Now().UnixNano())
	r1:=rand.New(s1)
	r := r1.Intn(9999)

	return r
}

// GenerateKey is the function that generate Reservation key
// params - key
// returns ReservationKey
func (s *SmartContract) GenerateKey(ctx contractapi.TransactionContextInterface, key string) *ReservationKey {
	fmt.Println("GenerateKey")
	var isFirst bool = false

	reservationKeyAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	var reservationKey *ReservationKey

	// initialize ReservationKey
	if reservationKeyAsBytes == nil {
		isFirst = true
		reservationKey = &ReservationKey{
			Key: "RESERVE",
			Idx: 1,
		}
	}

	var tempIdx string
	// when exist the key, increase the index
	if !isFirst {
		err = json.Unmarshal(reservationKeyAsBytes, &reservationKey)
		if err != nil {
			fmt.Println(err.Error())
		}
		tempIdx = strconv.Itoa(reservationKey.Idx)
		reservationKey.Idx = reservationKey.Idx + 1
	}

	fmt.Println("Last ReservationKey is " + reservationKey.Key + " : " + tempIdx)

	return reservationKey
}

// ReserveGround is the invoke function that makes a reservation the ground
// params - groundID, userID, begin and end time of the play
func (s *SmartContract) ReserveGround(ctx contractapi.TransactionContextInterface, groundID string, userID string, begin string, end string)error {
	fmt.Println("ReserveGround called")
	var reservationKey *ReservationKey
	reservationKey = s.GenerateKey(ctx, "latestKey")
	keyidx := strconv.Itoa(reservationKey.Idx)
	fmt.Println("Key : " + reservationKey.Key + ", Idx : " + keyidx)

	var reservationNumer = reservationKey.Key + keyidx
	fmt.Println("reservationKey is " + reservationNumer)

	reservationCompositeKey, _ := ctx.GetStub().CreateCompositeKey("reservation", []string{groundID, userID, reservationNumer})

	// parse the time
	beginTime := parseTime(begin)
	endTime := parseTime(end)

	randomGameCode := createRandomCode()

	// check the validation
	isPossible, err := validateReservation(ctx, groundID, beginTime, endTime)
	if err != nil {
		return fmt.Errorf("validate Error: %s", err.Error())
	}
	if isPossible {
		// create the Reservation
		reservation := Reservation{
			GroundID:          groundID,
			UserID:            userID,
			Begin:             beginTime,
			End:               endTime,
			ReservationNumber: reservationNumer,
			GameCode:          randomGameCode,
		}

		reservationAsBytes, err := json.Marshal(reservation)
		if err != nil {
			return fmt.Errorf("reservation Marshal Error: %s", err.Error())
		}

		err = ctx.GetStub().SetEvent("newReservation", reservationAsBytes)
		if err != nil {
			return fmt.Errorf("event Error: %s", err.Error())
		}

		reservationKeyAsBytes, _ := json.Marshal(reservationKey)
		ctx.GetStub().PutState("latestKey", reservationKeyAsBytes)

		ctx.GetStub().PutState(reservationCompositeKey, reservationAsBytes)
	} else {
		fmt.Println("That time is already reserved")
		return nil
	}
	return nil
}

func (s *SmartContract) UserConfirmReservation(ctx contractapi.TransactionContextInterface, userID string) ([]*Reservation, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("reservation", []string{userID})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var reservations []*Reservation

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var reservation Reservation

		_ = json.Unmarshal(queryResponse.Value, &reservation)

		reservations = append(reservations, &reservation)
	}

	return reservations, nil
}


// ConfirmReservation is the query function that confirms the reservation status given groundID and userID
// params - groundID, userID
// returns the array of reservations
func (s *SmartContract) ConfirmReservation(ctx contractapi.TransactionContextInterface, groundID string, userID string) ([]*Reservation, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("reservation", []string{groundID, userID})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var reservations []*Reservation

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var reservation Reservation

		_ = json.Unmarshal(queryResponse.Value, &reservation)

		reservations = append(reservations, &reservation)
	}

	return reservations, nil
}

// validateReservation is the function that validates the reservation according to given time
// params - groundID, begin and end time
// returns the true or false
func validateReservation(ctx contractapi.TransactionContextInterface, groundID string, beginTime, endTime time.Time) (bool, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("reservation", []string{groundID})
	if err != nil {
		return false, fmt.Errorf("%s", err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return false, fmt.Errorf("%s", err.Error())
		}

		var reservation Reservation

		_ = json.Unmarshal(queryResponse.Value, &reservation)
		// if (beginTime.After(reservation.Begin) && beginTime.Before(reservation.End)) || (endTime.After(reservation.Begin) && endTime.Before(reservation.End) || (beginTime.Equal(reservation.Begin) || endTime.Equal(reservation.End))) {
		// 	return false, nil
		// }
		if !(beginTime.After(reservation.End) || endTime.Before(reservation.Begin)) {
			return false, nil
		}
	}

	return true, nil
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