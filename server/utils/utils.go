package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func StringToInt64Array(s string) ([]int64, error) {
	var result []int64
	parts := strings.Split(s, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		num, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse '%s' as int64: %w", trimmed, err)
		}
		result = append(result, num)
	}
	return result, nil
}

func StrToInt64(s string) (int64, error) {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("conversion failed: %w", err)
	}
	return result, nil
}
