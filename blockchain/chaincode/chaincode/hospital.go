package chaincode

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type HospitalContract struct {
	contractapi.Contract
}

var hospitals = map[string]*Hospital{}

type Hospital struct {
	Name      string
	Contact   string
	Reports   map[int]MedicalReport
	Inventory map[string]HospitalDrug
	Patients  map[string]string
	Channels  map[string]bool
	mu        sync.Mutex
}

type HospitalDrug struct {
	Name         string
	TraceCode    string
	HospitalName string
}

type MedicalReport struct {
	ID          int
	PatientName string
	Symptoms    string
	NeededDrugs []string
}

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
		Patients:  make(map[string]string),
		Channels:  make(map[string]bool),
	}

	hospitals[name] = hospital
	hospitalJSON, err := json.Marshal(hospital)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(name, hospitalJSON)
}

// CreatePatientRecord creates a new patient record in the hospital's Patients map.
func (hc *HospitalContract) CreatePatientRecord(ctx contractapi.TransactionContextInterface, hospitalName string, patientName string) error {
	hospital, hospitalExists := hospitals[hospitalName]
	if !hospitalExists {
		return fmt.Errorf("hospital not found")
	}

	hospital.mu.Lock()
	defer hospital.mu.Unlock()

	hospital.Patients[patientName] = patientName
	hospitalJSON, err := json.Marshal(hospital)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hospitalName, hospitalJSON)
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
	for patientName := range hospital.Patients {
		patientList = append(patientList, patientName)
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

// AddDrugToHospitalInventory adds a drug to the hospital's inventory.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - hospitalName: the name of the hospital.
// - drugName: the name of the drug being added.
// - traceCode: the trace code of the drug.
//
// This function checks if the hospital exists. If not, it returns an error.
// It then adds the drug to the hospital's inventory.
func (hc *HospitalContract) AddDrugToHospitalInventory(ctx contractapi.TransactionContextInterface, hospitalName, drugName, traceCode string) error {
	hospital, hospitalExists := hospitals[hospitalName]
	if !hospitalExists {
		return fmt.Errorf("hospital not found")
	}

	hospital.mu.Lock()
	defer hospital.mu.Unlock()

	hospital.Inventory[traceCode] = HospitalDrug{
		Name:         drugName,
		TraceCode:    traceCode,
		HospitalName: hospitalName,
	}

	hospitalJSON, err := json.Marshal(hospital)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hospitalName, hospitalJSON)
}

// RemoveDrugFromHospitalInventory removes a drug from the hospital's inventory and returns its trace code.
// Parameters:
// - ctx: the transaction context provided by Hyperledger Fabric.
// - hospitalName: the name of the hospital selling the drug.
// - drugName: the name of the drug being sold.
//
// This function checks if the hospital exists. If not, it returns an error.
// If the drug is available in the hospital's inventory, it is removed and its trace code is returned.
//
// Returns:
// - string: the trace code if the operation is successful.
// - error: nil if the operation is successful, or an error message if it fails.
func (hc *HospitalContract) RemoveDrugFromHospitalInventory(ctx contractapi.TransactionContextInterface, hospitalName, drugName string) (string, error) {
	hospital, hospitalExists := hospitals[hospitalName]
	if !hospitalExists {
		return "", fmt.Errorf("hospital not found")
	}

	hospital.mu.Lock()
	defer hospital.mu.Unlock()

	for traceCode, drug := range hospital.Inventory {
		if drug.Name == drugName {
			delete(hospital.Inventory, traceCode)

			hospitalJSON, err := json.Marshal(hospital)
			if err != nil {
				return "", err
			}

			if err := ctx.GetStub().PutState(hospitalName, hospitalJSON); err != nil {
				return "", err
			}

			return traceCode, nil // 返回溯源码
		}
	}

	return "", fmt.Errorf("drug not available")
}

// check hospital is valid or not
func (hc *HospitalContract) ValidHospital(ctx contractapi.TransactionContextInterface, hospitalName string) (bool, error) {
	hospitalJSON, err := ctx.GetStub().GetState(hospitalName)
	if err != nil {
		return false, fmt.Errorf("failed to read hospital from world state: %v", err)
	}
	if hospitalJSON == nil {
		return false, nil // Hospital does not exist
	}

	return true, nil // Hospital exists
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
func (hc *HospitalContract) ViewReport(ctx contractapi.TransactionContextInterface, patientName, hospitalName string, reportID int) (MedicalReport, error) {
	hospital := hospitals[hospitalName]
	if hospital == nil {
		return MedicalReport{}, fmt.Errorf("hospital not found")
	}

	if report, exists := hospital.Reports[reportID]; exists && report.PatientName == patientName {
		return report, nil
	}

	return MedicalReport{}, fmt.Errorf("report not found")
}
