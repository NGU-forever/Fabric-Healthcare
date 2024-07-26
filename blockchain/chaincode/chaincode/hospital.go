/*
SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// HospitalDrug 表示医院中的药品信息
type HospitalDrug struct {
	Name         string
	TraceCode    string
	InStock      bool
	HospitalName string
}

// MedicalReport 表示体检报告
type MedicalReport struct {
	ID          int
	PatientName string
	Symptoms    string
	NeededDrugs []string
}

// Hospital 表示医院信息
type Hospital struct {
	Name      string
	Contact   string
	Reports   map[int]MedicalReport
	Inventory map[string]HospitalDrug
	Patients  map[string]bool
	Channels  map[string]bool
	mu        sync.Mutex
}

// DrugInfo 表示药品信息
type DrugInfo struct {
	Name           string
	TraceCode      string
	Manufacturer   string
	Price          float64
	ProductionTime string
	InStock        bool
}

// HospitalContract 定义了医院合约
type HospitalContract struct {
	contractapi.Contract
}

// 全局变量，存储所有医院信息
var hospitals = map[string]*Hospital{}

// var manufacturers = map[string]*Manufacturer{} // For simplicity in this example

// CreateHospital 创建一个新的医院
func (hc *HospitalContract) CreateHospital(ctx contractapi.TransactionContextInterface, name, contact string) error {
	hospital := &Hospital{
		Name:      name,
		Contact:   contact,
		Reports:   make(map[int]MedicalReport),
		Inventory: make(map[string]HospitalDrug),
		Patients:  make(map[string]bool),
		Channels:  make(map[string]bool),
	}

	hospitals[name] = hospital
	hospitalJSON, err := json.Marshal(hospital)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(name, hospitalJSON)
}

// ModifyReport 创建体检报告
func (hc *HospitalContract) modifyReport(ctx contractapi.TransactionContextInterface, hospitalName, patientName, symptoms string, neededDrugs []string) (int, error) {
	hospital := hospitals[hospitalName]
	hospital.mu.Lock()
	defer hospital.mu.Unlock()
	// 检查病人是否在 patients 名单中
	if _, exists := hospital.Patients[patientName]; !exists {
		return 0, fmt.Errorf("patient not found in hospital's list")
	}

	reportID := len(hospital.Reports) + 1
	report := MedicalReport{
		ID:          reportID,
		PatientName: patientName,
		Symptoms:    symptoms,
		NeededDrugs: neededDrugs,
	}

	hospital.Reports[reportID] = report
	hospitalJSON, err := json.Marshal(hospital)
	if err != nil {
		return 0, err
	}

	return reportID, ctx.GetStub().PutState(hospitalName, hospitalJSON)
}

// BuyDrug 从厂家购买药品
func (hc *HospitalContract) BuyDrug(ctx contractapi.TransactionContextInterface, hospitalName, manufacturerName, drugName string) (string, error) {
	hospital := hospitals[hospitalName]
	manufacturer := manufacturers[manufacturerName]

	manufacturer.mu.Lock()
	defer manufacturer.mu.Unlock()

	if drug, exists := manufacturer.Inventory[drugName]; exists && drug.InStock {
		drug.InStock = false // 更新制造商库存中的药品状态为已售出
		hospital.Inventory[drugName] = HospitalDrug{
			Name:         drug.Name,
			TraceCode:    drug.TraceCode,
			HospitalName: hospitalName,
			InStock:      true,
		}

		hospitalJSON, err := json.Marshal(hospital)
		if err != nil {
			return "", err
		}

		manufacturerJSON, err := json.Marshal(manufacturer)
		if err != nil {
			return "", err
		}

		if err := ctx.GetStub().PutState(hospitalName, hospitalJSON); err != nil {
			return "", err
		}

		if err := ctx.GetStub().PutState(manufacturerName, manufacturerJSON); err != nil {
			return "", err
		}

		return drug.TraceCode, nil
	}

	return "", fmt.Errorf("drug not available")
}

// TraceDrug 根据溯源码追踪药品
func (hc *HospitalContract) TraceDrug(ctx contractapi.TransactionContextInterface, traceCode string) (DrugInfo, error) {
	drugName, manufacturerName, price, productionTime, err := DecodeTraceCode(traceCode)
	if err != nil {
		return DrugInfo{}, err
	}

	return DrugInfo{
		Name:           drugName,
		TraceCode:      traceCode,
		Manufacturer:   manufacturerName,
		Price:          price,
		ProductionTime: productionTime,
		InStock:        false,
	}, nil
}
