package summary

import (
	"fmt"
	"io"
	"sort"

	finv1 "go.seankhliao.com/mono/apis/fin/v1"
)

func Run() {
	var all *finv1.All
	var w io.Writer

	fmt.Fprintf(w, "Summary for %s in %s\n", all.Name, all.Currency)
}

func delta(all *finv1.All) months {
	ms := make(months, 0, len(all.Months))
	for _, mo := range all.Months {
		m := month{
			year:       int(mo.Year),
			month:      int(mo.Month),
			categories: make(map[finv1.Transaction_Category]int64),
		}
		for _, tr := range mo.Transactions {
			m.categories[tr.Src] -= tr.Amount
			m.categories[tr.Dst] += tr.Amount
		}
		ms = append(ms, m)
	}
	sort.Sort(ms)
	return ms
}

func balances(w io.Writer, all *finv1.All) {
}

type month struct {
	year       int
	month      int
	categories map[finv1.Transaction_Category]int64
}
type months []month

func (m months) Len() int           { return len(m) }
func (m months) Less(i, j int) bool { return m[i].year*100+m[i].month < m[j].year*100+m[j].month }
func (m months) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
