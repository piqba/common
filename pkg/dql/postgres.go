package dql

import (
	sq "github.com/Masterminds/squirrel"
)

type PgOperator struct{}

func (p PgOperator) AndOP(filters []Filters) interface{} {
	query := sq.And{}
	for _, filter := range filters {

		switch filter.Type {
		case "TERM":
			query = append(query, sq.Eq{filter.Field: filter.Value})
		case "TERMS":
			query = append(query, sq.Eq{filter.Field: filter.Value})
		case "OR":
			subQueries := sq.Or{}
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					subQueries = append(subQueries, sq.Eq{sf.Field: sf.Value})
				} else {
					subQueries = append(subQueries, p.AndOP(sf.SubFilters).(sq.Eq))
				}
			}
			query = append(query, subQueries)
		case "AND":
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					query = append(query, sq.Eq{sf.Field: sf.Value})
				} else {
					p.AndOP(sf.SubFilters)
				}
			}
		case "NOT":
			query = append(query, sq.NotEq{filter.Field: filter.Value})
		}
	}
	return query
}

func (p PgOperator) OrOP(filters []Filters) interface{} {
	query := sq.Or{}
	for _, filter := range filters {

		switch filter.Type {
		case "TERM":
			query = append(query, sq.Eq{filter.Field: filter.Value})
		case "TERMS":
			query = append(query, sq.Eq{filter.Field: filter.Value})
		case "OR":
			subQueries := sq.Or{}
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					subQueries = append(subQueries, sq.Eq{sf.Field: sf.Value})
				} else {
					subQueries = append(subQueries, p.AndOP(sf.SubFilters).(sq.Eq))
				}
			}
			query = append(query, subQueries)
		case "AND":
			for _, sf := range filter.SubFilters {
				if len(sf.SubFilters) == 0 {
					query = append(query, sq.Eq{sf.Field: sf.Value})
				} else {
					p.AndOP(sf.SubFilters)
				}
			}
		case "NOT":
			query = append(query, sq.NotEq{filter.Field: filter.Value})
		}
	}
	return query
}

func (p PgOperator) NotOP(filters []Filters) interface{} {
	// Todo: Unimplemented yet
	return nil
}

func ToPostgres(filters string) (sql string, args []interface{}, err error) {
	dsl, err := transformer(filters)
	if err != nil {
		return "", nil, err
	}
	pgOP := PgOperator{}
	switch dsl.Filter.Type {
	case "AND":
		sql, args, err = pgOP.AndOP(dsl.Filter.Filters).(sq.And).ToSql()
		if err != nil {
			return sql, args, err
		}
		return sql, args, err
	case "OR":
		sql, args, err = pgOP.OrOP(dsl.Filter.Filters).(sq.Or).ToSql()
		return sql, args, err
	case "NOT":
		_ = pgOP.NotOP(dsl.Filter.Filters)
		return sql, args, err
	}
	return
}
