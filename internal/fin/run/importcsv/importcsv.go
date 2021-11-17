package importcsv

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"go.seankhliao.com/mono/internal/fin/run"
	"go.seankhliao.com/mono/internal/fin/store"
	fin "go.seankhliao.com/mono/proto/seankhliao/fin/v1alpha1"
)

type Options struct {
	File string
}

func NewOptions(fs *flag.FlagSet) Options {
	var o Options
	o.InitFlags(fs)
	return o
}

func (o *Options) InitFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.File, "f", "nl.csv", "path to csv to import")
}

func Run(o run.Options, args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	co := NewOptions(fs)
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	return importCSV(o, co)
}

func importCSV(o run.Options, co Options) error {
	records, err := getRecords(co.File)
	if err != nil {
		return err
	}

	all, err := parseRecords(records)
	if err != nil {
		return err
	}

	all.Name = co.File

	return store.WriteFile(o.File, all)
}

func getRecords(name string) ([][]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	return records, nil
}

func parseRecords(records [][]string) (*fin.All, error) {
	all := &fin.All{}

	a := make(map[int]map[int][]*fin.Transaction)

	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		year, month, tr, err := parseRecord(record)
		if err != nil {
			return nil, fmt.Errorf("record %d: %w", i, err)
		}
		if _, ok := a[year]; !ok {
			a[year] = make(map[int][]*fin.Transaction)
		}
		a[year][month] = append(a[year][month], tr)
	}

	years := make([]int, 0, len(a))
	for y := range a {
		years = append(years, y)
	}
	sort.Ints(years)
	for _, y := range years {
		months := make([]int, 0, 12)
		for m := range a[y] {
			months = append(months, m)
		}
		sort.Ints(months)
		for _, m := range months {
			all.Months = append(all.Months, &fin.Month{
				Year:         int32(y),
				Month:        int32(m),
				Transactions: a[y][m],
			})
		}
	}
	return all, nil
}

func parseRecord(r []string) (year, month int, tr *fin.Transaction, err error) {
	var monthString string
	_, err = fmt.Sscanf(r[0], "%d %s", &year, &monthString)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("parse date: %w", err)
	}
	month, ok := monthMap[monthString]
	if !ok {
		return 0, 0, nil, fmt.Errorf("unknown month %s", monthString)
	}

	amount, err := strconv.ParseInt(strings.ReplaceAll(r[1], ".", ""), 10, 64)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("parse amount %s: %w", r[1], err)
	}

	src, err := parseCategory(r[3])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("unknown src %s: %w", r[3], err)
	}

	dst, err := parseCategory(r[4])
	if err != nil {
		return 0, 0, nil, fmt.Errorf("unknown dst %s: %w", r[4], err)
	}

	return year, month, &fin.Transaction{
		Amount: amount,
		Src:    src,
		Dst:    dst,
		Note:   r[2],
	}, nil
}

func parseCategory(c string) (fin.Transaction_Category, error) {
	c = strings.ToUpper(c)
	if cat, ok := fin.Transaction_Category_value[c]; ok {
		return fin.Transaction_Category(cat), nil
	}
	if c == "OTHER" {
		return fin.Transaction_IN_OTHER, nil
	}
	return 0, fmt.Errorf("unknown")
}

var monthMap = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}
