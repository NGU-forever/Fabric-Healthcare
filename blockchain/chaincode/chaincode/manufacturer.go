/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ManufacturerDrug 表示药品信息
type ManufacturerDrug struct {
	Name           string
	TraceCode      string
	Manufacturer   string
	Price          float64
	ProductionTime string
	InStock        bool
}

// Manufacturer 表示厂家信息
type Manufacturer struct {
	Name      string
	Inventory map[string]ManufacturerDrug
	Contact   string
	Channels  map[string]bool
	mu        sync.Mutex
}

// ManufacturerContract 定义了厂家合约
type ManufacturerContract struct {
	contractapi.Contract
}

// 全局变量，存储所有厂家信息
var manufacturers = map[string]*Manufacturer{}

// CreateManufacturer 创建一个新的厂家
func (mc *ManufacturerContract) CreateManufacturer(ctx contractapi.TransactionContextInterface, name, contact string) error {
	manufacturer := &Manufacturer{
		Name:      name,
		Contact:   contact,
		Inventory: make(map[string]ManufacturerDrug),
		Channels:  make(map[string]bool),
	}

	manufacturers[name] = manufacturer
	manufacturerJSON, err := json.Marshal(manufacturer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(name, manufacturerJSON)
}

// ProduceDrug 生产药品并生成溯源码
func (mc *ManufacturerContract) ProduceDrug(ctx contractapi.TransactionContextInterface, manufacturerName, drugName string, price float64) (string, error) {
	manufacturer := manufacturers[manufacturerName]
	manufacturer.mu.Lock()
	defer manufacturer.mu.Unlock()

	productionTime := time.Now().Format(time.RFC3339)
	traceCode := GenerateTraceCode(drugName, manufacturerName, fmt.Sprintf("%.2f", price), productionTime)
	drug := ManufacturerDrug{
		Name:           drugName,
		TraceCode:      traceCode,
		Manufacturer:   manufacturerName,
		Price:          price,
		ProductionTime: productionTime,
		InStock:        true,
	}

	manufacturer.Inventory[drugName] = drug
	manufacturerJSON, err := json.Marshal(manufacturer)
	if err != nil {
		return "", err
	}

	return traceCode, ctx.GetStub().PutState(manufacturerName, manufacturerJSON)
}
