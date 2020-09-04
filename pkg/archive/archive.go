package archive

import (
	"gocms/pkg/page"
	"time"
)

// Generate the complete archive tree and save it in the store
func Generate() {

}

// GetYears - get a list of years containing posts
func GetYears() map[int][]page.Info {
	years := make(map[int][]page.Info)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.Status == page.Published && p.Template == "post" {
			years[p.PubDate.Year()] = append(years[p.PubDate.Year()], p)
		}
	}
	return years
}

// GetMonths - get a list of months containing posts for a given year
func GetMonths(year int) map[time.Month][]page.Info {
	months := make(map[time.Month][]page.Info)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.PubDate.Year() == year && p.Status == page.Published && p.Template == "post" {
			months[p.PubDate.Month()] = append(months[p.PubDate.Month()], p)
		}
	}
	return months
}

// GetDays - get a list of days containing posts, and the post meta data for a given year, month
func GetDays(year int, month time.Month) map[int][]page.Info {
	days := make(map[int][]page.Info)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.PubDate.Year() == year && p.PubDate.Month() == month && p.Status == page.Published && p.Template == "post" {
			days[p.PubDate.Day()] = append(days[p.PubDate.Day()], p)
		}
	}
	return days
}
