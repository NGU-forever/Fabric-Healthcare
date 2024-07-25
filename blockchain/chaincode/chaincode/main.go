/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 创建并注册所有合约
func main() {
	manufacturerContract := new(ManufacturerContract)
	hospitalContract := new(HospitalContract)
	patientContract := new(PatientContract)

	chaincode, err := contractapi.NewChaincode(manufacturerContract, hospitalContract, patientContract)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
