/*
SPDX-License-Identifier: Apache-2.0
*/

/*
SPDX-License-Identifier: Apache-2.0
*/

package main

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

// var hospitals = map[string]*Hospital{}

// CreatePatient 创建一个新的病人
func (pc *PatientContract) CreatePatient(ctx contractapi.TransactionContextInterface, name, birthDate string, height, weight float64, gender, contact string) error {
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

// ConductExamination 进行体检并在医院创建病人记录和空白体检报告
func (pc *PatientContract) ConductExamination(ctx contractapi.TransactionContextInterface, patientName, hospitalName string) (int, error) {
	hospital := hospitals[hospitalName]
	hospital.mu.Lock()
	defer hospital.mu.Unlock()

	// 在医院创建病人记录
	hospital.Patients[patientName] = false
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

// ViewReport 查看体检报告
func (pc *PatientContract) ViewReport(ctx contractapi.TransactionContextInterface, patientName, hospitalName string, reportID int) (MedicalReport, error) {
	hospital := hospitals[hospitalName]
	if report, exists := hospital.Reports[reportID]; exists && report.PatientName == patientName {
		return report, nil
	}

	return MedicalReport{}, fmt.Errorf("report not found")
}

// BuyDrug 从医院购买药品
func (pc *PatientContract) BuyDrug(ctx contractapi.TransactionContextInterface, patientName, hospitalName, drugName string) (string, error) {
	hospital := hospitals[hospitalName]
	hospital.mu.Lock()
	defer hospital.mu.Unlock()

	if drug, exists := hospital.Inventory[drugName]; exists {
		// drug.HospitalName = hospitalName // 记录购买的医院名称

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
