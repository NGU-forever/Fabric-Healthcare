/*
SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PatientContract 定义了病人合约
type PatientContract struct {
	contractapi.Contract
}

// Patient 表示病人信息
type Patient struct {
	Name      string
	BirthDate string
	Height    float64
	Weight    float64
	Gender    string
	Contact   string
}

// PatientDrug 表示病人拥有的药品信息
type PatientDrug struct {
	Name         string
	TraceCode    string
	HospitalName string
}

// DrugInfoPatient 表示药品信息
type DrugInfoPatient struct {
	Name           string
	TraceCode      string
	Manufacturer   string
	Price          float64
	ProductionTime string
	HospitalName   string
}

// 全局变量，存储所有病人信息
var patients = map[string]*Patient{}

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

// ConductExamination conducts a medical examination and creates a patient record and a blank medical report in the specified hospital.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - patientName: the name of the patient.
// - hospitalName: the name of the hospital.
//
// This function checks if the hospital and patient exist. If either does not exist, it returns an error.
// It then creates a blank medical report and adds the patient to the hospital's patient list.
//
// Returns:
// - int: the ID of the newly created medical report.
// - error: nil if the operation is successful, or an error message if the hospital or patient does not exist or there is an error during storage.
func (pc *PatientContract) ConductExamination(ctx contractapi.TransactionContextInterface, patientName, hospitalName string) (int, error) {
	// 检查医院是否存在
	hospital := hospitals[hospitalName]
	if hospital == nil {
		return 0, fmt.Errorf("hospital not found")
	}

	hospital.mu.Lock()
	defer hospital.mu.Unlock()

	// 检查病人是否已存在
	patient, exists := patients[patientName]
	if !exists {
		return 0, fmt.Errorf("patient not found")
	}

	// 在医院创建病人记录
	hospital.Patients[patientName] = *patient
	hospitalJSON, err := json.Marshal(hospital)
	if err != nil {
		return 0, err
	}
	if err := ctx.GetStub().PutState(hospitalName, hospitalJSON); err != nil {
		return 0, err
	}

	// 创建空白体检报告
	reportID := len(hospital.Reports) + 1
	report := MedicalReport{
		ID:          reportID,
		PatientName: patientName,
		Symptoms:    "",
		NeededDrugs: []string{},
	}

	hospital.Reports[reportID] = report
	hospitalJSON, err = json.Marshal(hospital)
	if err != nil {
		return 0, err
	}
	if err := ctx.GetStub().PutState(hospitalName, hospitalJSON); err != nil {
		return 0, err
	}

	return reportID, nil
}

// ViewReport retrieves a medical report for a patient from a specified hospital.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - patientName: the name of the patient.
// - hospitalName: the name of the hospital.
// - reportID: the ID of the medical report to retrieve.
//
// This function checks if the hospital exists and if the medical report exists and belongs to the specified patient.
// If either condition is not met, it returns an error. Otherwise, it returns the medical report.
//
// Returns:
// - MedicalReport: the medical report if found.
// - error: nil if the operation is successful, or an error message if the hospital or report does not exist.
func (pc *PatientContract) ViewReport(ctx contractapi.TransactionContextInterface, patientName, hospitalName string, reportID int) (MedicalReport, error) {
	hospital := hospitals[hospitalName]
	if hospital == nil {
		return MedicalReport{}, fmt.Errorf("hospital not found")
	}

	if report, exists := hospital.Reports[reportID]; exists && report.PatientName == patientName {
		return report, nil
	}

	return MedicalReport{}, fmt.Errorf("report not found")
}

// BuyDrug allows a patient to purchase a drug from a specified hospital.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - patientName: the name of the patient (not used in the current logic, but included for completeness).
// - hospitalName: the name of the hospital.
// - drugName: the name of the drug to purchase.
//
// This function checks if the hospital exists and if the specified drug is available in the hospital's inventory.
// If either condition is not met, it returns an error. If the drug is available, it updates its status to sold and stores the updated information in the world state.
//
// Returns:
// - string: the trace code of the purchased drug.
// - error: nil if the operation is successful, or an error message if the hospital or drug does not exist or there is an error during storage.
func (pc *PatientContract) BuyDrug(ctx contractapi.TransactionContextInterface, patientName, hospitalName, drugName string) (string, error) {
	hospital := hospitals[hospitalName]
	if hospital == nil {
		return "", fmt.Errorf("hospital not found")
	}

	hospital.mu.Lock()
	defer hospital.mu.Unlock()

	if drug, exists := hospital.Inventory[drugName]; exists && drug.InStock {
		drug.InStock = false // 更新医院库存中的药品状态为已售出

		hospitalJSON, err := json.Marshal(hospital)
		if err != nil {
			return "", err
		}

		if err := ctx.GetStub().PutState(hospitalName, hospitalJSON); err != nil {
			return "", err
		}

		return drug.TraceCode, nil
	}

	return "", fmt.Errorf("drug not available")
}

// TraceDrug 根据溯源码追踪药品
func (pc *PatientContract) TraceDrug(ctx contractapi.TransactionContextInterface, traceCode string) (DrugInfoPatient, error) {
	drugName, manufacturerName, price, productionTime, err := DecodeTraceCode(traceCode)
	if err != nil {
		return DrugInfoPatient{}, err
	}

	// 遍历所有医院以找到对应的药品
	for _, hospital := range hospitals {
		for _, drug := range hospital.Inventory {
			if drug.TraceCode == traceCode {
				return DrugInfoPatient{
					Name:           drugName,
					TraceCode:      traceCode,
					Manufacturer:   manufacturerName,
					Price:          price,
					ProductionTime: productionTime,
					HospitalName:   hospital.Name,
				}, nil
			}
		}
	}

	return DrugInfoPatient{}, fmt.Errorf("drug not found")
}

func (pc *PatientContract) GetPatients(ctx contractapi.TransactionContextInterface) ([]string, error) {
	var patientList []string
	for patientName := range manufacturers {
		patientList = append(patientList, patientName)
	}
	return patientList, nil
}
