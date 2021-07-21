package queries

import (
	"strconv"
	"strings"
	"time"
)

func inputTables(q QueryExperiment, tables ...string) []string {
	inputs := make([]string, 0)

	inputs = append(inputs, strings.Join([]string{q.GetPrefix(), "lineitem", "*"}, "/"))

	return inputs
}

func SQLDate(val string) interface{} {
	s_year, _ := strconv.Atoi(val[0:4])
	s_month, _ := strconv.Atoi(val[6:7])
	s_day, _ := strconv.Atoi(val[9:10])

	//Shipdate, _ := time.Parse("2006-01-02", line[L_SHIPDATE])
	return time.Date(s_year, time.Month(s_month), s_day, 0, 0, 0, 0, time.Local)
}

func Integer(val string) interface{} {
	r, _ := strconv.ParseInt(val, 10, 64)
	return r
}

func Float(val string) interface{} {
	r, _ := strconv.ParseFloat(val, 64)
	return r
}
