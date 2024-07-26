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

	// GenerateTraceCode 生成一个唯一的溯源码，包含药品名称、厂家、价格和生产时间
	func GenerateTraceCode(drugName, manufacturer, price, productionTime string) string {
		return fmt.Sprintf("%s-%s-%s-%s-%d", drugName, manufacturer, price, productionTime, rand.Int())
	}

	// DecodeTraceCode 解码溯源码并返回药品信息
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
