package dql

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoOperator struct{}

func (m *MongoOperator) AndOP(filters []Filters) interface{} {
	query := bson.M{}
	for _, filter := range filters {

		switch filter.Type {
		case "TERM":
			query[filter.Field] = filter.Value
		case "TERMS":
			query[filter.Field] = bson.M{"$in": filter.Value}
		case "OR":
			subQueries := primitive.A{}
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					subQueries = append(subQueries, bson.M{sf.Field: sf.Value})
				} else {

					subQueries = append(subQueries, m.AndOP(sf.SubFilters))
				}
			}
			query["$or"] = subQueries
		case "AND":
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					query[sf.Field] = sf.Value
				} else {
					m.AndOP(sf.SubFilters)
				}
			}
		case "NOT":
			query[filter.Field] = bson.M{"$not": filter.Value}
		}
	}
	return query
}

func (m *MongoOperator) OrOP(filters []Filters) interface{} {
	queryOr := bson.M{}
	elementsOr := bson.A{}
	for _, filter := range filters {

		switch filter.Type {
		case "TERM":
			query := bson.M{}

			query[filter.Field] = filter.Value
			elementsOr = append(elementsOr, query)

		case "TERMS":
			query := bson.M{}

			query[filter.Field] = bson.M{"$in": filter.Value}
			elementsOr = append(elementsOr, query)

		case "OR":
			query := bson.M{}

			subQueries := primitive.A{}
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					subQueries = append(subQueries, bson.M{sf.Field: sf.Value})
				} else {
					subQueries = append(subQueries, m.OrOP(sf.SubFilters))
				}
			}
			query["$or"] = subQueries
			elementsOr = append(elementsOr, query)

		case "AND":
			query := bson.M{}

			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					query[sf.Field] = sf.Value
				}
				m.OrOP(sf.SubFilters)
			}
			elementsOr = append(elementsOr, query)
		case "NOT":
			query := bson.M{filter.Field: bson.M{"$not": filter.Value}}
			elementsOr = append(elementsOr, query)
		}
	}
	queryOr["$or"] = elementsOr

	return queryOr
}

func (m *MongoOperator) NotOP(filters []Filters) interface{} {
	queryNor := bson.M{}
	elementsNor := bson.A{}
	for _, filter := range filters {

		switch filter.Type {
		case "TERM":
			query := bson.M{}

			query[filter.Field] = filter.Value
			elementsNor = append(elementsNor, query)

		case "TERMS":
			query := bson.M{}

			query[filter.Field] = bson.M{"$in": filter.Value}
			elementsNor = append(elementsNor, query)

		case "OR":
			query := bson.M{}

			subQueries := primitive.A{}
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					subQueries = append(subQueries, bson.M{sf.Field: sf.Value})
				} else {
					subQueries = append(subQueries, m.NotOP(sf.SubFilters))
				}
			}
			query["$or"] = subQueries
			elementsNor = append(elementsNor, query)

		case "AND":
			query := bson.M{}

			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					query[sf.Field] = sf.Value
				}
				m.NotOP(sf.SubFilters)
			}
			elementsNor = append(elementsNor, query)

		}
	}
	queryNor["$nor"] = elementsNor
	return queryNor
}

func ToMongoDB(filters string) (query bson.M, err error) {
	dsl, err := transformer(filters)
	if err != nil {
		return nil, err
	}
	mop := MongoOperator{}
	switch dsl.Filter.Type {
	case "AND":
		query = mop.AndOP(dsl.Filter.Filters).(bson.M)
	case "OR":
		query = mop.OrOP(dsl.Filter.Filters).(bson.M)
	case "NOT":
		query = mop.NotOP(dsl.Filter.Filters).(bson.M)
	}
	return
}
