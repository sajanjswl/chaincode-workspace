/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	shell "github.com/ipfs/go-ipfs-api"
)

const ipfsURL = "http://localhost:5001"

// student document store map
var StudentDocStoreMap map[string]Student

var sh *shell.Shell

// SmartContract provides functions for managing a student
type SmartContract struct {
	contractapi.Contract
	sh *shell.Shell
}

// Student IPFS Record
type StudentIPFSCID struct {
	CID string `json:"cid"`
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
	Key string `json:"key"`
	CID string `json:"cid"`
}

// InitLedger adds a base set of student to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	s.sh = shell.NewShell(ipfsURL)

	students := []Student{
		{FirstName: "Sajan", LastName: "Jaiswal", Branch: "CSE", RegistrationNumber: "1816123", BloodGroup: "A+", MobileNumber: "+917064274923", Address: "White House, Motihari, Bihar"},
		{FirstName: "Abhishek", LastName: "Jaiswal", Branch: "CSE", RegistrationNumber: "1816124", BloodGroup: "B+", MobileNumber: "+918210791275", Address: "MidLand,Dimapur"},
	}

	StudentDocStoreMap = make(map[string]Student)

	for _, student := range students {

		// Map the struct instance to the mapping
		StudentDocStoreMap[student.RegistrationNumber] = student

		// Converting the map into JSON object
		studentAsBytes, _ := json.Marshal(StudentDocStoreMap)

		// Dag PUT operation which will return the CID for futher access or pinning etc.
		cid, err := s.sh.DagPut(studentAsBytes, "json", "cbor")
		if err != nil {
			return fmt.Errorf("failed to put student record to ipfs %s %s", student.RegistrationNumber, err.Error())

		}

		fmt.Println("prinitng cid", cid)
		cidAsBytes, _ := json.Marshal(StudentIPFSCID{CID: cid})

		if err = ctx.GetStub().PutState(student.RegistrationNumber, cidAsBytes); err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

	}

	return nil
}

// RegisterStudent adds a new Student to the world state with given details
func (s *SmartContract) RegisterStudent(ctx contractapi.TransactionContextInterface, registrationNumber, firstName, lastName, branch, bloodGroup, mobileNumber, address string) error {

	sh = shell.NewShell(ipfsURL)
	student := Student{
		FirstName:          firstName,
		LastName:           lastName,
		Branch:             branch,
		RegistrationNumber: registrationNumber,
		BloodGroup:         bloodGroup,
		MobileNumber:       mobileNumber,
		Address:            address,
	}

	// Map the struct instance to the mapping
	StudentDocStoreMap[student.RegistrationNumber] = student

	// Converting the map into JSON object
	studentAsBytes, _ := json.Marshal(StudentDocStoreMap)

	// Dag PUT operation which will return the CID for futher access or pinning etc.
	cid, err := sh.DagPut(studentAsBytes, "json", "cbor")
	if err != nil {
		return fmt.Errorf("failed to put student record to ipfs %s %s", student.RegistrationNumber, err.Error())

	}
	fmt.Println("the cid is ", cid)

	cidAsBytes, _ := json.Marshal(StudentIPFSCID{CID: cid})

	return ctx.GetStub().PutState(student.RegistrationNumber, cidAsBytes)
}

// QueryStudent returns the student stored in the world state with given id
func (s *SmartContract) QueryStudent(ctx contractapi.TransactionContextInterface, registrationNumber string) (*Student, error) {
	cidAsBytes, err := ctx.GetStub().GetState(registrationNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if cidAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", registrationNumber)
	}

	studentIPFSCID := new(StudentIPFSCID)
	_ = json.Unmarshal(cidAsBytes, studentIPFSCID)

	student, err := s.GetDocument(studentIPFSCID.CID, registrationNumber)
	if err != nil {
		return nil, fmt.Errorf("%s does not exist on ipfs %s", registrationNumber, err.Error())
	}
	return &student, nil
}

func (s *SmartContract) GetDocument(ref, key string) (out Student, err error) {
	err = s.sh.DagGet(ref+"/"+key, &out)
	return
}

// QueryAllStudent returns all cars found in world state
// func (s *SmartContract) QueryAllStudent(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
// 	startKey := ""
// 	endKey := ""

// 	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	results := []QueryResult{}

// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()

// 		if err != nil {
// 			return nil, err
// 		}

// 		student := new(Student)
// 		_ = json.Unmarshal(queryResponse.Value, student)

// 		queryResult := QueryResult{Key: queryResponse.Key, Record: student}
// 		results = append(results, queryResult)
// 	}

// 	return results, nil
// }

// ChangeCarOwner updates the owner field of car with given id in world state
// func (s *SmartContract) UpdateMobileNumber(ctx contractapi.TransactionContextInterface, registrationNumber string, mobileNumber string) error {
// 	student, err := s.QueryStudent(ctx, registrationNumber)

// 	if err != nil {
// 		return err
// 	}
// 	student.MobileNumber = mobileNumber

// 	studentAsBytes, _ := json.Marshal(student)

// 	return ctx.GetStub().PutState(registrationNumber, studentAsBytes)
// }

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
