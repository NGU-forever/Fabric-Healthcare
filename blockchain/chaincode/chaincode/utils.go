/*
SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// GenerateTraceCode generates a unique trace code that includes the drug name, manufacturer, price, and production time.
// Parameters:
// - drugName: the name of the drug.
// - manufacturer: the name of the manufacturer.
// - price: the price of the drug.
// - productionTime: the time the drug was produced.
//
// This function concatenates the provided parameters with a randomly generated integer to create a unique trace code.
//
// Returns:
// - string: the generated trace code.
func GenerateTraceCode(drugName, manufacturer, price, productionTime string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%d", drugName, manufacturer, price, productionTime, rand.Int())
}

// DecodeTraceCode decodes a trace code to extract the drug name, manufacturer, price, and production time.
// Parameters:
// - traceCode: the trace code to decode.
//
// This function splits the trace code into its constituent parts and converts the price from string to float.
//
// Returns:
// - string: the drug name.
// - string: the manufacturer name.
// - float64: the price of the drug.
// - string: the production time.
// - error: an error message if the trace code format is invalid or the price conversion fails.
func DecodeTraceCode(traceCode string) (string, string, float64, string, error) {
	parts := strings.Split(traceCode, "-")
	if len(parts) < 5 {
		return "", "", 0, "", fmt.Errorf("invalid trace code format")
	}

	drugName := parts[0]
	manufacturer := parts[1]
	price, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return "", "", 0, "", fmt.Errorf("invalid price format")
	}
	productionTime := parts[3]

	return drugName, manufacturer, price, productionTime, nil
}
