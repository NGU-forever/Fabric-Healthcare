package chaincode

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type HospitalDrug struct {
	Name         string
	TraceCode    string
	InStock      bool
	HospitalName string
}

type MedicalReport struct {
	ID          int
	PatientName string
	Symptoms    string
	NeededDrugs []string
}

type Hospital struct {
	Name      string
	Contact   string
	Reports   map[int]MedicalReport
	Inventory map[string]HospitalDrug
	Patients  map[string]Patient
	Channels  map[string]bool
	mu        sync.Mutex
}

type DrugInfo struct {
	Name           string
	TraceCode      string
	Manufacturer   string
	Price          float64
	ProductionTime string
}

type HospitalContract struct {
	contractapi.Contract
}

var hospitals = map[string]*Hospital{}

// CreateHospital creates a new hospital record in the ledger.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - name: the name of the hospital to be created.
// - contact: the contact information of the hospital.
//
// This function first checks if a hospital with the given name already exists.
// If it does, it returns an error indicating the hospital already exists.
// If not, it creates a new Hospital struct, initializes its fields, and stores it
// in the hospitals map. It then serializes the hospital struct to JSON and stores
// it in the ledger using PutState.
//
// Returns:
//   - error: nil if the operation is successful, or an error message if it fails or
//     the hospital already exists.
func (hc *HospitalContract) CreateHospital(ctx contractapi.TransactionContextInterface, name, contact string) error {
	if _, exists := hospitals[name]; exists {
		return fmt.Errorf("hospital already exists")
	}

	hospital := &Hospital{
		Name:      name,
		Contact:   contact,
		Reports:   make(map[int]MedicalReport),
		Inventory: make(map[string]HospitalDrug),
		Patients:  make(map[string]Patient),
		Channels:  make(map[string]bool),
	}

	hospitals[name] = hospital
	hospitalJSON, err := json.Marshal(hospital)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(name, hospitalJSON)
}

// ModifyReport modifies a medical report for a patient in a hospital.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - hospitalName: the name of the hospital where the report is being created.
// - patientName: the name of the patient for whom the report is being created.
// - symptoms: the symptoms reported by the patient.
// - neededDrugs: a list of drugs needed by the patient.
//
// This function first checks if the hospital exists. If not, it returns an error.
// Then it checks if the patient exists in the hospital's records. If not, it returns an error.
// If both exist, it creates a new medical report, assigns it an ID, and stores it in the hospital's reports map.
// The function then serializes the hospital struct to JSON and updates the hospital record in the ledger.
//
// Returns:
// - int: the ID of the created report if the operation is successful.
// - error: nil if the operation is successful, or an error message if it fails.
func (hc *HospitalContract) ModifyReport(ctx contractapi.TransactionContextInterface, hospitalName, patientName, symptoms string, neededDrugs []string) (int, error) {
	hospital, hospitalExists := hospitals[hospitalName]
	if !hospitalExists {
		return 0, fmt.Errorf("hospital not found")
	}

	hospital.mu.Lock()
	defer hospital.mu.Unlock()

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

// BuyDrug allows a hospital to purchase a drug from a manufacturer.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - hospitalName: the name of the hospital buying the drug.
// - manufacturerName: the name of the manufacturer selling the drug.
// - drugName: the name of the drug being purchased.
//
// This function first checks if the hospital and manufacturer exist. If either does not exist, it returns an error.
// If both exist, it checks if the drug is available in the manufacturer's inventory and if it's in stock.
// If the drug is available and in stock, it updates the drug's status to 'not in stock' in the manufacturer's inventory
// and adds the drug to the hospital's inventory. The hospital and manufacturer records are then updated in the ledger.
//
// Returns:
// - string: the trace code of the purchased drug if the operation is successful.
// - error: nil if the operation is successful, or an error message if it fails.
func (hc *HospitalContract) BuyDrug(ctx contractapi.TransactionContextInterface, hospitalName, manufacturerName, drugName string) (string, error) {
	hospital, hospitalExists := hospitals[hospitalName]
	if !hospitalExists {
		return "", fmt.Errorf("hospital not found")
	}

	manufacturer, manufacturerExists := manufacturers[manufacturerName]
	if !manufacturerExists {
		return "", fmt.Errorf("manufacturer not found")
	}

	manufacturer.mu.Lock()
	defer manufacturer.mu.Unlock()

	if drug, exists := manufacturer.Inventory[drugName]; exists && drug.InStock {
		drug.InStock = false
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

// TraceDrug traces the drug based on its trace code.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - traceCode: the unique trace code of the drug.
//
// This function decodes the trace code to retrieve the drug's information including
// the drug's name, manufacturer, price, and production time. It returns this information
// encapsulated in a DrugInfo struct.
//
// Returns:
// - DrugInfo: a struct containing the traced drug's information if the operation is successful.
// - error: nil if the operation is successful, or an error message if it fails.
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
	}, nil
}

// GetPatients returns a list of patients in a hospital.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - hospitalName: the name of the hospital.
//
// This function first checks if the hospital exists. If not, it returns an error.
// If the hospital exists, it retrieves the list of patients from the hospital's
// Patients map and returns it.
//
// Returns:
// - []string: a list of patient names if the operation is successful.
// - error: nil if the operation is successful, or an error message if it fails.
func (hc *HospitalContract) GetPatients(ctx contractapi.TransactionContextInterface, hospitalName string) ([]string, error) {
	hospital, hospitalExists := hospitals[hospitalName]
	if !hospitalExists {
		return nil, fmt.Errorf("hospital not found")
	}

	var patientList []string
	for patient := range hospital.Patients {
		patientList = append(patientList, patient)
	}

	return patientList, nil
}

// // GetHospitals retrieves a list of all hospitals.
// // Parameters:
// // - ctx: the transaction context provided by Hyperledger Fabric.
// //
// // This function iterates through the hospitals map and appends each hospital's name to a slice.
// // It returns the list of hospital names.
// //
// // Returns:
// // - []string: a slice containing the names of all hospitals if the operation is successful.
// // - error: nil if the operation is successful, or an error message if it fails.
func (hc *HospitalContract) GetHospitals(ctx contractapi.TransactionContextInterface) ([]string, error) {
	var hospitalList []string
	for hospitalName := range hospitals {
		hospitalList = append(hospitalList, hospitalName)
	}
	return hospitalList, nil
}
