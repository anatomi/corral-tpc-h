package queries

import (
	"fmt"
	"github.com/anatomi/corral"
	"math/rand"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Q11 struct {
	Experiment
	Nation   string
	Fraction   float64
	//Sum float64
}

func (q *Q11) Name() string {
	return fmt.Sprintf("%s_tcph_q11_n%s_f%g", q.ShortName(), q.Nation, q.Fraction)
}

func (q *Q11) Check(driver *corral.Driver) error {
	panic("implement me")
}

func (q *Q11) Configure() []corral.Option {
	inputs := [][]string{
		inputTables(q, "supplier", "nation"),
		inputTables(q, "partsupp"),
		inputTables(q, "job1"),
		inputTables(q, "job2"),
		[]string{},
	}

	return []corral.Option{
		corral.WithMultiStageInputs(inputs),
		corral.WithSplitSize(25 * 1024 * 1024),
		corral.WithMapBinSize(100 * 1024 * 1024),
		corral.WithReduceBinSize(200 * 1024 * 1024),
	}
}

func (q *Q11) Validate(strings []string) (bool, error) {
	panic("implement me")
}

func (q *Q11) Serialize() map[string]string {
	m := make(map[string]string)
	m["fraction"] = fmt.Sprintf("%g", q.Fraction)
	m["nation"] = q.Nation
	return m
}

func (q *Q11) Read(m map[string]string) error {
	fraction, err := strconv.ParseFloat(m["fraction"], 64)
	if err != nil {
		return err
	}

	q.Fraction = float64(fraction)
	q.Nation = m["nation"]

	return nil
}

func (q *Q11) Default() {
	q.Fraction = 0.0001
	q.Nation = "GERMANY"
}

func (q *Q11) Randomize() {
	q.Fraction = rand.Float64() // To DO
	q.Nation = RandomNation()
}

/*
select
	ps_partkey,
	sum(ps_supplycost * ps_availqty) as value
from
	partsupp,
	supplier,
	nation
where
	ps_suppkey = s_suppkey
	and s_nationkey = n_nationkey
	and n_name = '[NATION]'
group by
	ps_partkey having
		sum(ps_supplycost * ps_availqty) > (
			select
				sum(ps_supplycost * ps_availqty) * [FRACTION]
			from
				partsupp,
				supplier,
				nation
			where
				ps_suppkey = s_suppkey
				and s_nationkey = n_nationkey
				and n_name = '[NATION]'
			)
order by
	value desc;
*/
func (q *Q11) Create() []*corral.Job {
	// supplierJoin
	supplierJoin := &Join{
		Query: q,
		left:  Supplier(),
		right: Nation(),
		on: [2]int{int(S_NATIONKEY), int(N_NATIONKEY)},
		filter: [2]Projection{
			func(table *GenericTable) []int {
				return []int{
					//0
					int(S_SUPPKEY),
				}
			},
			func(table *GenericTable) []int {
				if table.Get(int(N_NAME)) == q.Nation {
					//1
					return []int{int(N_NAME)}
				} else {
					return nil
				}
			},
		},
	}

	//partsuppJoin
	partsuppJoin := &Join{
		Query: q,
		left: &GenericTable{
			Name: "job0",
			numFields: 2,
		},
		right: Partsupp(),
		on: [2]int{int(S_SUPPKEY), int(PS_SUPPKEY)},
		filter: [2]Projection{
			nil,
			func(table *GenericTable) []int {
				return []int{
					//2,3,4
					int(PS_PARTKEY), int(PS_SUPPLYCOST), int(PS_AVAILQTY),
				}
			},
		},
	}

	sumAvaiableParts := &Q11TotalAvaiableParts{
		q: q,
	}

	selectValues := &Q11SelectValues{
		q: q,
	}

	return []*corral.Job{
		corral.NewJob(supplierJoin, supplierJoin),
		corral.NewJob(partsuppJoin, partsuppJoin),
		corral.NewJob(sumAvaiableParts, sumAvaiableParts),
		corral.NewJob(selectValues, selectValues),
		//TODO: SORT
		//corral.NewSort(),
	}
}

type Q11TotalAvaiableParts struct{
	q *Q11
}

func (q Q11TotalAvaiableParts) Map(key, value string, emitter corral.Emitter) {
	emitter.Emit("", value)
}

func (q Q11TotalAvaiableParts) Reduce(key string, values corral.ValueIterator, emitter corral.Emitter) {
	tab := GenericTable{
		numFields: 5,
	}

	vals := []string{}
	sum_suppliercost_avialability := 0.0
	for l := range values.Iter() {
		tab.Read(l)

		partkey := tab.Get(2)

		suppyliercost, err := strconv.ParseFloat(tab.Get(3), 64)
		if err != nil {
			log.Infof("Error: %v", err)
		}
		avaialbility, err := strconv.ParseFloat(tab.Get(4), 64)
		if err != nil {
			log.Infof("Error: %v", err)
		}
		vals = append(vals, fmt.Sprintf("%s|%f|%f", partkey, suppyliercost, avaialbility))
		sum_suppliercost_avialability += (suppyliercost * avaialbility)
	}

	for _, v := range vals {
		emitter.Emit("", fmt.Sprintf("%s|%f", v, sum_suppliercost_avialability))
	}

}

type Q11SelectValues struct{
	q *Q11
}

func (q Q11SelectValues) Map(key, value string, emitter corral.Emitter) {
	tab := GenericTable{
		numFields: 4,
	}
	tab.Read(value)
	partkey := tab.Get(0)
	supplycost, err := strconv.ParseFloat(tab.Get(1), 64)
	if err != nil {
		log.Infof("Error: %v", err)
	}
	avaialbility, err := strconv.ParseFloat(tab.Get(2), 64)
	if err != nil {
		log.Infof("Error: %v", err)
	}
	total_sum, err := strconv.ParseFloat(tab.Get(3), 64)
	emitter.Emit(partkey, fmt.Sprintf("%f|%f|%f", supplycost, avaialbility, total_sum))
}

func (q Q11SelectValues) Reduce(key string, values corral.ValueIterator, emitter corral.Emitter) {
	tab := GenericTable{
		numFields: 3,
	}
	sum := 0.0
	total_sum := 0.0
	for l := range values.Iter() {

		tab.Read(l)
		supplycost, err := strconv.ParseFloat(tab.Get(0), 64)
		if err != nil {
			log.Infof("Error: %v", err)
		}
		avaialbility, err := strconv.ParseFloat(tab.Get(1), 64)
		if err != nil {
			log.Infof("Error: %v", err)
		}
		sum += (supplycost*avaialbility)
		total_sum, err = strconv.ParseFloat(tab.Get(2), 64)
		if err != nil {
			log.Infof("Error: %v", err)
		}
	}

	if sum > (total_sum*q.q.Fraction) {
		err := emitter.Emit(key, fmt.Sprintf("|%.2f", sum))
		if err != nil {
			log.Infof("failed to emit %s,+%v", key, err)
		}
	}

	
}
