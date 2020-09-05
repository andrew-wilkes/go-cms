package archive

import (
	"gocms/pkg/page"
	"time"
)

// GetYears - get a map of years containing posts with the post count
func GetYears() map[int]int {
	years := make(map[int]int)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.Status == page.Published && p.Template == "post" {
			years[p.PubDate.Year()]++
		}
	}
	return years
}

// GetMonths - get a map of months containing posts with the post count
func GetMonths(year int) map[time.Month]int {
	months := make(map[time.Month]int)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.PubDate.Year() == year && p.Status == page.Published && p.Template == "post" {
			months[p.PubDate.Month()]++
		}
	}
	return months
}

// GetDays - get a map of days containing posts, and the post meta data for a given year, month
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
