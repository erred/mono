package summary

import (
	"flag"
	"fmt"
	"strings"

	"go.seankhliao.com/mono/cli/cmd/fin/internal/run"
	"go.seankhliao.com/mono/cli/cmd/fin/internal/store"
	fin "go.seankhliao.com/mono/proto/seankhliao/fin/v1alpha1"
)

type Options struct {
	Year     int
	Month    int
	Delta    bool
	Holdings bool
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	o.InitFlags(fs)
	return &o
}

func (o *Options) InitFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.Year, "y", 0, "year, 0 means all")
	fs.IntVar(&o.Month, "m", 0, "month, 0 means all")
	fs.BoolVar(&o.Delta, "delta", false, "print deltas per month")
	fs.BoolVar(&o.Holdings, "holdings", true, "print final value per month")
}

func Run(o run.Options, args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	so := NewOptions(fs)
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	all, err := store.ReadFile(o.File)
	if err != nil {
		return fmt.Errorf("read %s: %w", o.File, err)
	}

	var months []*fin.Month
	for i, m := range all.Months {
		if (so.Year == 0 || so.Year == int(m.Year)) && (so.Month == 0 || so.Month == int(m.Month)) {
			months = append(months, all.Months[i])
		}
	}

	mss := make([]MonthSummary, 0, len(months))
	for _, m := range months {
		ms := summarize(m)
		mss = append(mss, ms)
	}

	fmt.Println(so.Delta, so.Holdings)
	switch {
	case so.Delta:
		for _, ms := range mss {
			fmt.Println(ms)
		}
	case so.Holdings:
		var gs GroupSummary
		for _, ms := range mss {
			gs.Name = fmt.Sprintf("%d-%02d", ms.Year, ms.Month)
			gs.Total -= ms.Holdings.Total
			gs.Categories = ms.Holdings.Categories
			if len(gs.Deltas) == 0 {
				gs.Deltas = make([]int64, len(ms.Holdings.Deltas))
			}
			for j := range ms.Holdings.Deltas {
				gs.Deltas[j] -= ms.Holdings.Deltas[j]
			}
			fmt.Println(gs)
		}
	}

	return nil
}

type MonthSummary struct {
	Year, Month int
	Holdings    GroupSummary
	Income      GroupSummary
	Exppenses   GroupSummary

	Delta map[fin.Transaction_Category]int64
}

type GroupSummary struct {
	Name  string
	Total int64

	Categories []string
	Deltas     []int64
}

func (s GroupSummary) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%s\t%+.2f\n", s.Name, float64(s.Total)/100)
	for i, c := range s.Categories {
		if i != 0 {
			b.WriteRune('\t')
		}
		fmt.Fprintf(b, "%-12s", c)
	}
	b.WriteRune('\n')
	for i, d := range s.Deltas {
		if i != 0 {
			b.WriteRune('\t')
		}
		fmt.Fprintf(b, "%+-12.2f", float64(d)/100)
	}
	b.WriteRune('\n')
	return b.String()
}

func (s MonthSummary) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%d-%02d\n", s.Year, s.Month)
	b.WriteString(s.Holdings.String())
	b.WriteString(s.Income.String())
	b.WriteString(s.Exppenses.String())
	b.WriteRune('\n')
	return b.String()
}

func summarize(m *fin.Month) MonthSummary {
	delta := make(map[fin.Transaction_Category]int64)
	for _, tr := range m.Transactions {
		delta[tr.Src] += tr.Amount
		delta[tr.Dst] -= tr.Amount
	}

	return MonthSummary{
		Year:      int(m.Year),
		Month:     int(m.Month),
		Holdings:  group("Holdings", fin.Transaction_CASH, fin.Transaction_BITTREX, delta),
		Income:    group("Income", fin.Transaction_SALARY, fin.Transaction_IN_OTHER, delta),
		Exppenses: group("Expenses", fin.Transaction_FOOD, fin.Transaction_EDUCATION, delta),
		Delta:     delta,
	}
}

func group(name string, first, last fin.Transaction_Category, delta map[fin.Transaction_Category]int64) GroupSummary {
	s := GroupSummary{
		Name: name,
	}
	for i := first; i <= last; i++ {
		s.Categories = append(s.Categories, i.String())
		s.Deltas = append(s.Deltas, delta[i])
		s.Total += delta[i]
	}

	return s
}
