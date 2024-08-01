package chaincode

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ManufacturerContract struct {
	contractapi.Contract
}

var manufacturers = map[string]*Manufacturer{}

type Manufacturer struct {
	Name      string
	Inventory map[string]ManufacturerDrug
	Contact   string
	Channels  map[string]bool
	mu        sync.Mutex
}

type ManufacturerDrug struct {
	Name           string
	TraceCode      string
	Manufacturer   string
	Price          float64
	ProductionTime string
}

// CreateManufacturer creates a new manufacturer.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - name: the name of the manufacturer to create.
// - contact: the contact information of the manufacturer.
//
// This function checks if the manufacturer already exists. If it does, it returns an error.
// Otherwise, it creates a new manufacturer with the provided name and contact information,
// and initializes its inventory and channels. The manufacturer is then stored in the world state.
//
// Returns:
// - error: nil if the operation is successful, or an error message if it fails or the manufacturer already exists.
func (mc *ManufacturerContract) CreateManufacturer(ctx contractapi.TransactionContextInterface, name, contact string) error {
	if _, exists := manufacturers[name]; exists {
		return fmt.Errorf("manufacturer already exists")
	}

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

// AddDrugToMnfcInventory adds a new drug to the manufacturer's inventory.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - manufacturerName: the name of the manufacturer producing the drug.
// - drugName: the name of the drug to add.
// - traceCode: the trace code of the drug.
// - price: the price of the drug.
//
// This function checks if the manufacturer exists. If not, it returns an error.
// It then adds the drug with the provided trace code to the manufacturer's inventory
// and stores the updated manufacturer in the world state.
//
// Returns:
// - error: nil if the operation is successful, or an error message if it fails.
func (mc *ManufacturerContract) AddDrugToMnfcInventory(ctx contractapi.TransactionContextInterface, manufacturerName, drugName, traceCode string, price float64) error {
	manufacturer, exists := manufacturers[manufacturerName]
	if !exists {
		return fmt.Errorf("manufacturer not found")
	}

	manufacturer.mu.Lock()
	defer manufacturer.mu.Unlock()

	drug := ManufacturerDrug{
		Name:           drugName,
		TraceCode:      traceCode,
		Manufacturer:   manufacturerName,
		Price:          price,
		ProductionTime: time.Now().Format(time.RFC3339),
	}

	manufacturer.Inventory[traceCode] = drug
	manufacturerJSON, err := json.Marshal(manufacturer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(manufacturerName, manufacturerJSON)
}

func (mc *ManufacturerContract) GetManufacturers(ctx contractapi.TransactionContextInterface) ([]string, error) {
	var manufacturerList []string
	for manufacturerName := range manufacturers {
		manufacturerList = append(manufacturerList, manufacturerName)
	}
	return manufacturerList, nil
}

// RemoveDrugFromMnfcInventory removes a drug from the inventory and returns its trace code.
func (mc *ManufacturerContract) RemoveDrugFromMnfcInventory(ctx contractapi.TransactionContextInterface, manufacturerName, drugName string) (string, error) {
	manufacturer, exists := manufacturers[manufacturerName]
	if !exists {
		return "", fmt.Errorf("manufacturer not found")
	}

	manufacturer.mu.Lock()
	defer manufacturer.mu.Unlock()

	for traceCode, drug := range manufacturer.Inventory {
		if drug.Name == drugName {
			delete(manufacturer.Inventory, traceCode)

			manufacturerJSON, err := json.Marshal(manufacturer)
			if err != nil {
				return "", err
			}

			if err := ctx.GetStub().PutState(manufacturerName, manufacturerJSON); err != nil {
				return "", err
			}

			return traceCode, nil // 找到药品并删除后立即返回
		}
	}

	return "", fmt.Errorf("drug not available")
}
