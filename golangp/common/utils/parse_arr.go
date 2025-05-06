package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/toheart/functrace"
	"strings"
)

func ParseTags(tagsStr string) ([]string, error) {
	defer functrace.Trace([]interface {
	}{tagsStr})()

	if tagsStr == "" {
		return []string{}, errors.New("Invalid tags: empty string")
	}
	var tags []string

	if strings.HasPrefix(tagsStr, "[") && strings.HasSuffix(tagsStr, "]") {

		err := json.Unmarshal([]byte(tagsStr), &tags)
		if err != nil {
			return nil, fmt.Errorf("Invalid JSON array: %v", err)
		}
	} else {

		tags = strings.Split(tagsStr, ",")
	}

	return tags, nil
}
