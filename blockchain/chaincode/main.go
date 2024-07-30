package main

import (
	"log"

	"chaincode/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 创建并注册所有合约
func main() {
	manufacturerContract := new(chaincode.ManufacturerContract)
	hospitalContract := new(chaincode.HospitalContract)
	patientContract := new(chaincode.PatientContract)

	cc, err := contractapi.NewChaincode(manufacturerContract, hospitalContract, patientContract)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := cc.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
