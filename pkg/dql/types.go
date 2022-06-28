package dql

/*
E.g:
// type LOGIC OPERATOR
// filters ARRAY of elements
{
  "filters": {
    "type": "AND",
    "filters": [
      {
        "type": "TERM",
        "field": "name",
        "value": "pepe"
      },
      {
        "type": "TERM",
        "field": "age",
        "value": 20
      },
      {
        "type": "OR",
        "filters": [
          {
            "type": "TERM",
            "field": "province",
            "value": "La Habana"
          },
          {
            "type": "TERM",
            "field": "province",
            "value": "Las Tunas"
          }
        ]
      }
    ]
  }
}

*/

type DSLFilter struct {
	Filter Filter `json:"filters,omitempty"`
}

type Filter struct {
	Type    string    `json:"type,omitempty"`
	Filters []Filters `json:"filters,omitempty"`
}

type Filters struct {
	Type       string      `json:"type,omitempty"`
	Field      string      `json:"field,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	SubFilters []Filters   `json:"filters,omitempty"`
}

type Operator interface {
	AndOP(filters []Filters) interface{}
	OrOP(filters []Filters) interface{}
	NotOP(filters []Filters) interface{}
}
