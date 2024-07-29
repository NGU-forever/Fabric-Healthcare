package chaincode

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ManufacturerDrug struct {
	Name           string
	TraceCode      string
	Manufacturer   string
	Price          float64
	ProductionTime string
	InStock        bool
}

type Manufacturer struct {
	Name      string
	Inventory map[string]ManufacturerDrug
	Contact   string
	Channels  map[string]bool
	mu        sync.Mutex
}

type ManufacturerContract struct {
	contractapi.Contract
}

var manufacturers = map[string]*Manufacturer{}

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

// ProduceDrug produces a new drug and generates a trace code for it.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - manufacturerName: the name of the manufacturer producing the drug.
// - drugName: the name of the drug to produce.
// - price: the price of the drug.
//
// This function checks if the manufacturer exists. If not, it returns an error.
// It also checks if the drug already exists in the manufacturer's inventory.
// If the drug already exists, it returns an error. Otherwise, it creates a new drug
// with a generated trace code and the provided details. The drug is then added to the
// manufacturer's inventory and stored in the world state.
//
// Returns:
//   - string: the generated trace code if the operation is successful.
//   - error: nil if the operation is successful, or an error message if it fails,
//     the manufacturer does not exist, or the drug already exists.
func (mc *ManufacturerContract) ProduceDrug(ctx contractapi.TransactionContextInterface, manufacturerName, drugName string, price float64) (string, error) {
	manufacturer, exists := manufacturers[manufacturerName]
	if !exists {
		return "", fmt.Errorf("manufacturer not found")
	}

	manufacturer.mu.Lock()
	defer manufacturer.mu.Unlock()

	if _, exists := manufacturer.Inventory[drugName]; exists {
		return "", fmt.Errorf("drug already exists")
	}

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

func (mc *Manufacturer) GetManufacturers(ctx contractapi.TransactionContextInterface) ([]string, error) {
	var manufacturerList []string
	for manufacturerName := range manufacturers {
		manufacturerList = append(manufacturerList, manufacturerName)
	}
	return manufacturerList, nil
}
