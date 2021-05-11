/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a student
type SmartContract struct {
	contractapi.Contract
}

// Student  describes basic details of what makes up a student
type Student struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Branch             string `json:"branch"`
	RegistrationNumber string `json:"registrationNumber"`
	BloodGroup         string `json:"bloodGroup"`
	MobileNumber       string `json:"mobileNumber"`
	Address            string `json:"address"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Student
}

// InitLedger adds a base set of student to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	students := []Student{
		{FirstName: "Sajan", LastName: "Jaiswal", Branch: "CSE", RegistrationNumber: "1816123", BloodGroup: "A+", MobileNumber: "+917064274923", Address: "White House, Motihari, Bihar"},
		{FirstName: "Abhishek", LastName: "Jaiswal", Branch: "CSE", RegistrationNumber: "1816124", BloodGroup: "B+", MobileNumber: "+918210791275", Address: "MidLand,Dimapur"},
	}

	for _, student := range students {
		studentAsBytes, _ := json.Marshal(student)
		err := ctx.GetStub().PutState(student.RegistrationNumber, studentAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// RegisterStudent adds a new Student to the world state with given details
func (s *SmartContract) RegisterStudent(ctx contractapi.TransactionContextInterface, registrationNumber, firstName, lastName, branch, bloodGroup, mobileNumber, address string) error {
	student := Student{
		FirstName:          firstName,
		LastName:           lastName,
		Branch:             branch,
		RegistrationNumber: registrationNumber,
		BloodGroup:         bloodGroup,
		MobileNumber:       mobileNumber,
		Address:            address,
	}

	studentAsBytes, _ := json.Marshal(student)

	return ctx.GetStub().PutState(registrationNumber, studentAsBytes)
}

// QueryStudent returns the student stored in the world state with given id
func (s *SmartContract) QueryStudent(ctx contractapi.TransactionContextInterface, registrationNumber string) (*Student, error) {
	studentAsBytes, err := ctx.GetStub().GetState(registrationNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if studentAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", registrationNumber)
	}

	student := new(Student)
	_ = json.Unmarshal(studentAsBytes, student)

	return student, nil
}

// QueryAllStudent returns all cars found in world state
func (s *SmartContract) QueryAllStudent(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		student := new(Student)
		_ = json.Unmarshal(queryResponse.Value, student)

		queryResult := QueryResult{Key: queryResponse.Key, Record: student}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) UpdateMobileNumber(ctx contractapi.TransactionContextInterface, registrationNumber string, mobileNumber string) error {
	student, err := s.QueryStudent(ctx, registrationNumber)

	if err != nil {
		return err
	}
	student.MobileNumber = mobileNumber

	studentAsBytes, _ := json.Marshal(student)

	return ctx.GetStub().PutState(registrationNumber, studentAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create student chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting student chaincode: %s", err.Error())
	}
}
