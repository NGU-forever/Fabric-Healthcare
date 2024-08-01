/*
SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PatientContract
type PatientContract struct {
	contractapi.Contract
}

// 全局变量，存储所有病人信息
var patients = map[string]*Patient{}

// Patient shows the info of patients
type Patient struct {
	Name      string
	BirthDate string
	Height    float64
	Weight    float64
	Gender    string
	Contact   string
}

// CreatePatient creates a new patient record.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - name: the name of the patient.
// - birthDate: the birth date of the patient.
// - height: the height of the patient.
// - weight: the weight of the patient.
// - gender: the gender of the patient.
// - contact: the contact information of the patient.
//
// This function checks if a patient with the given name already exists. If so, it returns an error.
// If not, it creates a new Patient struct, adds it to the global patients map, and stores it in the world state.
//
// Returns:
// - error: nil if the operation is successful, or an error message if the patient already exists or there is an error during storage.
func (pc *PatientContract) CreatePatient(ctx contractapi.TransactionContextInterface, name, birthDate string, height, weight float64, gender, contact string) error {
	if _, exists := patients[name]; exists {
		return fmt.Errorf("patient already exists")
	}

	patient := &Patient{
		Name:      name,
		BirthDate: birthDate,
		Height:    height,
		Weight:    weight,
		Gender:    gender,
		Contact:   contact,
	}

	patients[name] = patient
	patientJSON, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(name, patientJSON)
}

func (pc *PatientContract) GetPatients(ctx contractapi.TransactionContextInterface) ([]string, error) {
	var patientList []string
	for patientName := range patients {
		patientList = append(patientList, patientName)
	}
	return patientList, nil
}
