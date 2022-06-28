package dql

import "encoding/json"

func transformer(filters string) (DSLFilter, error) {
	dsl := DSLFilter{}
	err := json.Unmarshal([]byte(filters), &dsl)
	if err != nil {
		return DSLFilter{}, err
	}
	return dsl, err
}

func Transformer(filters string) (Filters, error) {
	dsl := Filters{}
	err := json.Unmarshal([]byte(filters), &dsl)
	if err != nil {
		return Filters{}, err
	}
	return dsl, err
}
