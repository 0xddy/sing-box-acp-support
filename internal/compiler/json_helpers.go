package compiler

import (
	"encoding/json"
)

func applyOptionalString(target map[string]any, key string, value string) {
	if value != "" {
		target[key] = value
	}
}

func applyStringList(target map[string]any, key string, values []string) {
	if len(values) > 0 {
		target[key] = values
	}
}

func applyUint32List(target map[string]any, key string, values []uint32) {
	if len(values) > 0 {
		target[key] = values
	}
}

func applyInt32List(target map[string]any, key string, values []int32) {
	if len(values) > 0 {
		target[key] = values
	}
}

func applyNestedStruct(target map[string]any, key string, value any) {
	compiled := structMap(value)
	if len(compiled) > 0 {
		target[key] = compiled
	}
}

func mergeStruct(target map[string]any, value any) {
	for key, fieldValue := range structMap(value) {
		target[key] = fieldValue
	}
}

func structMap(value any) map[string]any {
	if value == nil {
		return nil
	}
	encoded, err := json.Marshal(value)
	if err != nil || string(encoded) == "null" || string(encoded) == "{}" {
		return nil
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil
	}
	return decoded
}
