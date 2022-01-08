package summary

import (
	"flag"
	"fmt"
	"strings"

	"go.seankhliao.com/mono/cli/cmd/fin/internal/run"
	"go.seankhliao.com/mono/cli/cmd/fin/internal/store"
	"go.seankhliao.com/mono/proto/finpb"
)

type Options struct{}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	o.InitFlags(fs)
	return &o
}

func (o *Options) InitFlags(fs *flag.FlagSet) {}

func Run(o run.Options, args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	all, err := store.ReadFile(o.File)
	if err != nil {
		return fmt.Errorf("read %s: %w", o.File, err)
	}

	balance, income, expense := summarize(all)
	fmt.Println(expense)
	fmt.Println(income)
	fmt.Println(balance)
	return nil
}

type GroupSummary struct {
	Name       string
	Categories []string
	Months     []MonthSummary
}

type MonthSummary struct {
	Year, Month int
	Total       []int64
	Delta       []int64
}

func (s GroupSummary) String() string {
	b := &strings.Builder{}
	// Category header
	fmt.Fprintf(b, "%-7s", s.Name)
	for _, c := range s.Categories {
		fmt.Fprintf(b, "%20s", c)
	}
	fmt.Fprintf(b, "\n")
	// Month line
	for _, m := range s.Months {
		fmt.Fprintln(b, m)
	}
	fmt.Fprintf(b, "\n")
	return b.String()
}

func (s MonthSummary) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%d-%02d", s.Year, s.Month)
	for i := range s.Delta {
		fmt.Fprintf(b, "%20s", fmt.Sprintf("%8.2f %+8.2f", float64(s.Total[i])/100, float64(s.Delta[i])/100))
	}
	return b.String()
}

func summarize(all *finpb.All) (balance, income, expense GroupSummary) {
	balance = GroupSummary{
		Name:       "Balance",
		Categories: group(finpb.Transaction_CASH, finpb.Transaction_BITTREX),
		Months:     months(len(all.Months), int(finpb.Transaction_BITTREX-finpb.Transaction_CASH)+1),
	}
	income = GroupSummary{
		Name:       "Income",
		Categories: group(finpb.Transaction_SALARY, finpb.Transaction_IN_OTHER),
		Months:     months(len(all.Months), int(finpb.Transaction_IN_OTHER-finpb.Transaction_SALARY)+1),
	}
	expense = GroupSummary{
		Name:       "Expense",
		Categories: group(finpb.Transaction_FOOD, finpb.Transaction_EDUCATION),
		Months:     months(len(all.Months), int(finpb.Transaction_EDUCATION-finpb.Transaction_FOOD)+1),
	}

	total := make(map[finpb.Transaction_Category]int64)
	for i, m := range all.Months {
		delta := make(map[finpb.Transaction_Category]int64)
		for _, tr := range m.Transactions {
			delta[tr.Src] -= tr.Amount
			delta[tr.Dst] += tr.Amount
		}
		for c, a := range delta {
			total[c] += a
		}

		balance.Months[i].Year = int(m.Year)
		balance.Months[i].Month = int(m.Month)
		income.Months[i].Year = int(m.Year)
		income.Months[i].Month = int(m.Month)
		expense.Months[i].Year = int(m.Year)
		expense.Months[i].Month = int(m.Month)
		for c := range total {
			idx := categoryIdx(balance.Categories, c.String())
			if idx != -1 {
				balance.Months[i].Delta[idx] = delta[c]
				balance.Months[i].Total[idx] = total[c]
			}
			idx = categoryIdx(income.Categories, c.String())
			if idx != -1 {
				income.Months[i].Delta[idx] = delta[c] * -1
				income.Months[i].Total[idx] = total[c] * -1
			}
			idx = categoryIdx(expense.Categories, c.String())
			if idx != -1 {
				expense.Months[i].Delta[idx] = delta[c] * -1
				expense.Months[i].Total[idx] = total[c] * -1
			}
		}
	}

	return balance, income, expense
}

func group(first, last finpb.Transaction_Category) []string {
	var ss []string
	for i := first; i <= last; i++ {
		ss = append(ss, i.String())
	}
	return ss
}

func months(months, categories int) []MonthSummary {
	ms := make([]MonthSummary, months)
	for i := range ms {
		ms[i].Delta = make([]int64, categories)
		ms[i].Total = make([]int64, categories)
	}
	return ms
}

func categoryIdx(cs []string, c string) int {
	for i, s := range cs {
		if c == s {
			return i
		}
	}
	return -1
}
