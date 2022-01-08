package summary

import (
	"fmt"
	"io"
	"sort"

	"go.seankhliao.com/mono/proto/finpb"
)

func Run() {
	var all *finpb.All
	var w io.Writer

	fmt.Fprintf(w, "Summary for %s in %s\n", all.Name, all.Currency)
}

func delta(all *finpb.All) months {
	ms := make(months, 0, len(all.Months))
	for _, mo := range all.Months {
		m := month{
			year:       int(mo.Year),
			month:      int(mo.Month),
			categories: make(map[finpb.Transaction_Category]int64),
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

func balances(w io.Writer, all *finpb.All) {
}

type month struct {
	year       int
	month      int
	categories map[finpb.Transaction_Category]int64
}
type months []month

func (m months) Len() int           { return len(m) }
func (m months) Less(i, j int) bool { return m[i].year*100+m[i].month < m[j].year*100+m[j].month }
func (m months) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
