/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the vehicle structure, with 4 properties.  Structure tags are used by encoding/json library
type Vehicle struct {
	Type         string `json:"type"`
	Manufactor   string `json:"manufactor"`
	Model        string `json:"model"`
	Registration string `json:"Registration"`
}
// Define the log structure, with 3 properties.  Structure tags are used by encoding/json library
type Log struct {
	Date     string `json:"date"`
	Mileage  string `json:"mileage"`
	Comment  string `json:"comment"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryVehicle" {
		return s.queryVehicle(APIstub, args)
	} else if function == "queryMileage" {
		return s.queryMileage(APIstub, args)
	} else if function == "createVehicle" {
		return s.createVehicle(APIstub, args)
	} else if function == "createMileage" {
		return s.createMileage(APIstub, args)
	} else if function == "queryAllVehicles" {
		return s.queryAllVehicles(APIstub)
	} else if function == "queryMileageHistory" {
		return s.queryMileageHistory(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name: " + function)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	vehicles := []Vehicle{
		Vehicle{Type: "Kfz", Manufactor: "VW", Model: "California", Registration: "ZWS-IT 001"},
		Vehicle{Type: "Krad", Manufactor: "BMW", Model: "R 1100 R", Registration: "ZWS-IT 002"},
	}
	logs := []Log{
		Log{Date: "2014-02-04", Mileage: "23.0", Comment: "Tanken"}, // first log of first vehicle
		Log{Date: "2010-05-14", Mileage: "86010.0", Comment: "Tanken"}, // first log of second vehicle
	}

	i := 0
	for i < len(vehicles) {
		fmt.Println("i is ", i)
		vehicleAsBytes, _ := json.Marshal(vehicles[i])
		APIstub.PutState("FIN-"+fmt.Sprintf("%03d", i+1), vehicleAsBytes)
		fmt.Println("Added", vehicles[i])
		logAsBytes, _ := json.Marshal(logs[i])
		APIstub.PutState("LOG-"+fmt.Sprintf("%03d", i+1), logAsBytes)
		fmt.Println("Added", logs[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryVehicle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	if !(strings.HasPrefix(args[0], "FIN-")) {
		return shim.Error("Wrong prefix. Expecting FIN-")
	}

	vehicleAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(vehicleAsBytes)
}

func (s *SmartContract) queryMileage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	if !(strings.HasPrefix(args[0], "FIN-")) {
		return shim.Error("Wrong prefix. Expecting FIN-")
	}

	var key = strings.Replace(args[0], "FIN-", "LOG-", 1)
	logAsBytes, _ := APIstub.GetState(key)
	return shim.Success(logAsBytes)
}

func (s *SmartContract) createVehicle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	if !(strings.HasPrefix(args[0], "FIN-")) {
		return shim.Error("Wrong prefix. Expecting FIN-")
	}

	var vehicle = Vehicle{Type: args[1], Manufactor: args[2], Model: args[3], Registration: args[4]}

	vehicleAsBytes, _ := json.Marshal(vehicle)
	APIstub.PutState(args[0], vehicleAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) createMileage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if !(strings.HasPrefix(args[0], "FIN-")) {
		return shim.Error("Wrong prefix. Expecting FIN-")
	}

	var log = Log{Date: args[1], Mileage: args[2], Comment: args[3]}
	var key = strings.Replace(args[0], "FIN-", "LOG-", 1)

	logAsBytes, _ := json.Marshal(log)
	APIstub.PutState(key, logAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllVehicles(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "FIN-001"
	endKey := "FIN-999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllVehicles:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryMileageHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	if !(strings.HasPrefix(args[0], "FIN-")) {
		return shim.Error("Wrong prefix. Expecting FIN-")
	}

	var key = strings.Replace(args[0], "FIN-", "LOG-", 1)
	resultsIterator, err := APIstub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")


	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryMileageHistory:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
