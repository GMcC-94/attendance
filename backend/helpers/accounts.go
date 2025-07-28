package helpers

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/gmcc94/attendance-go/types"
)

func GroupedAccounts(entries []types.AccountEntry) map[int]map[int]map[string][]types.AccountEntry {
	grouped := make(map[int]map[int]map[string][]types.AccountEntry)
	for _, e := range entries {
		year, month, day := e.CreatedAt.Year(), int(e.CreatedAt.Month()), e.CreatedAt.Format("02 Jan 2006")
		if grouped[year] == nil {
			grouped[year] = make(map[int]map[string][]types.AccountEntry)
		}
		if grouped[year][month] == nil {
			grouped[year][month] = make(map[string][]types.AccountEntry)
		}
		grouped[year][month][day] = append(grouped[year][month][day], e)
	}
	return grouped
}

func BuildGroupedResponse(grouped map[int]map[int]map[string][]types.AccountEntry) []types.GroupedAccounts {
	var result []types.GroupedAccounts

	// --- Sort years ---
	var years []int
	for y := range grouped {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years))) // newest year first

	for _, year := range years {
		monthsMap := grouped[year]

		// --- Sort months ---
		var months []int
		for m := range monthsMap {
			months = append(months, m)
		}
		sort.Ints(months) // Jan -> Dec

		var monthList []types.GroupedMonthAccounts
		for _, month := range months {
			daysMap := monthsMap[month]

			// --- Sort days ---
			var dayKeys []time.Time
			dayKeyMap := make(map[time.Time]string)
			for d := range daysMap {
				parsed, _ := time.Parse("02 Jan 2006", d)
				dayKeys = append(dayKeys, parsed)
				dayKeyMap[parsed] = d
			}
			sort.Slice(dayKeys, func(i, j int) bool {
				return dayKeys[i].After(dayKeys[j]) // newest first
			})

			var dayList []types.GroupedDayAccounts
			for _, t := range dayKeys {
				dayStr := dayKeyMap[t]
				dayList = append(dayList, types.GroupedDayAccounts{
					Date:    dayStr,
					Records: daysMap[dayStr],
				})
			}

			// Convert month int to name
			monthName := time.Month(month).String()

			monthList = append(monthList, types.GroupedMonthAccounts{
				Month:   monthName,
				Entries: dayList,
			})
		}

		result = append(result, types.GroupedAccounts{
			Year:   year,
			Months: monthList,
		})
	}

	return result
}

func ValidateEntries(entries []types.AccountEntry) error {
	for _, e := range entries {
		if strings.TrimSpace(e.Description) == "" {
			return errors.New("description cannot be empty")
		}
		if e.Amount <= 0 {
			return errors.New("amount must be greater than 0")
		}
	}
	return nil
}
