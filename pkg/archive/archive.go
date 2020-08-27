package archive

import "gocms/pkg/page"

// YearItem type
type YearItem struct {
	year  int16
	count int16
}

// MonthItem type
type MonthItem struct {
	month int16
	count int16
}

// DayItem type
type DayItem struct {
	day   int16
	posts []page.Info
}

// Generate the complete archive tree and save it in the store
func Generate() {

}

// GetYears - get a list of years containing posts
func GetYears() []YearItem {
	return []YearItem{YearItem{2012, 8}, YearItem{2020, 6}}
}

// GetMonths - get a list of months containing posts
func GetMonths(year int) []MonthItem {
	return []MonthItem{MonthItem{4, 2}, MonthItem{6, 34}}
}

// GetMonth - get a list of days containing posts, and the post meta data
func GetMonth(month int) []DayItem {
	return []DayItem{
		DayItem{
			day: 1,
			posts: []page.Info{
				page.Get(1, false),
			},
		},
	}
}
