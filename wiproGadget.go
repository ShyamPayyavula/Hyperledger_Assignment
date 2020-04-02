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
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Gadget structure, with 4 properties.  Structure tags are used by encoding/json library
type Gadget struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "Gadget" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "Gadget"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryGadget" {
		return s.queryGadget(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createGadget" {
		return s.createGadget(APIstub, args)
	} else if function == "queryAllGadgets" {
		return s.queryAllGadgets(APIstub)
	} else if function == "changeGadgetOwner" {
		return s.changeGadgetOwner(APIstub, args)
	} else if function == "getGadgetHistory" {
		return s.getGadgetHistory(APIstub, args)
	} else if function == "deleteGadget" {
		return s.deleteGadget(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.") 
}

func (s *SmartContract) queryGadget(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	} 

	gadgetAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(gadgetAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	gadgets := []Gadget{
                Gadget{Make: "Samsung", Model: "GalS10", Colour: "blue", Owner: "Shyam"},
                Gadget{Make: "Apple", Model: "Ipod", Colour: "black", Owner: "Sravan"},
                Gadget{Make: "Apple", Model: "AirMac", Colour: "Silver", Owner: "Pavan"},
                Gadget{Make: "Sony", Model: "Bravia10", Colour: "Black", Owner: "Karthick"},
                Gadget{Make: "SkullCandy", Model: "EarPhones", Colour: "black", Owner: "Jhon"},
                Gadget{Make: "Nokia", Model: "N72", Colour: "Black", Owner: "Lalitha"},
                Gadget{Make: "MI", Model: "MI20", Colour: "Gray", Owner: "Jack"},
                Gadget{Make: "Samsung", Model: "GalS10Plus", Colour: "DarkBlack", Owner: "Akram"},
	}

	i := 0
	for i < len(gadgets) {
		fmt.Println("i is ", i)
		gadgetAsBytes, _ := json.Marshal(gadgets[i])
		APIstub.PutState("GADGET"+strconv.Itoa(i), gadgetAsBytes)
		fmt.Println("Added", gadgets[i])
		i = i + 1
	}

	return shim.Success(nil)
}

// This function is used to Create a new Gadget record and update 
func (s *SmartContract) createGadget(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var gadget = Gadget{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}

	gadgetAsBytes, _ := json.Marshal(gadget)
	APIstub.PutState(args[0], gadgetAsBytes)

	return shim.Success(nil)
}

// This function is used to query all the gadgets at a time 

func (s *SmartContract) queryAllGadgets(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "GADGET0"
	endKey := "GADGET999"

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

	fmt.Printf("- queryAllGadgets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//This function is used to Change the owner of the Gadget

func (s *SmartContract) changeGadgetOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	gadgetAsBytes, _ := APIstub.GetState(args[0])
	gadget := Gadget{}

	json.Unmarshal(gadgetAsBytes, &gadget)
	gadget.Owner = args[1]

	gadgetAsBytes, _ = json.Marshal(gadget)
	APIstub.PutState(args[0], gadgetAsBytes)

	return shim.Success(nil)
}




//Gets the History of the Gadgets  

func (s *SmartContract) getGadgetHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gadgetID := args[0]

	fmt.Printf("- start getGadgetHistory for : %s\n", gadgetID)

	resultsIterator, err := APIstub.GetHistoryForKey(gadgetID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getGadgetHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

//Using this function to Delete the Gadget from the record But will update in the gethistory as deleted 
func (s *SmartContract) deleteGadget(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var jsonResp string
	var gadgetJSON Gadget
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	gadgetID := args[0]
	// to maintain the color~name index, we need to read the marble first and get its color
	valAsbytes, err := APIstub.GetState(gadgetID) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gadgetID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"gadget does not exist: " + gadgetID + "\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal([]byte(valAsbytes), &gadgetJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + gadgetID + "\"}"
		return shim.Error(jsonResp)
	}

	err = APIstub.DelState(gadgetID) //remove the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract 
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

