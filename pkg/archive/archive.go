package archive

import (
	"gocms/pkg/page"
	"sort"
	"time"
)

// Record type
type Record struct {
	Month time.Month
	Count int
	Posts []page.Info
}

// GetYears - get a slice of years containing posts with the post count
func GetYears() map[int]Record {
	years := make(map[int]Record)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.Status == page.Published && p.Template == "post" {
			year := p.PubDate.Year()
			r, present := years[year]
			if present {
				r.Count++
			} else {
				r = Record{Count: 1}
			}
			years[year] = r
		}
	}
	return years
}

// GetMonths - get a map of months containing posts with the post count
func GetMonths(year int) map[int]Record {
	months := make(map[int]Record)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.PubDate.Year() == year && p.Status == page.Published && p.Template == "post" {
			month := p.PubDate.Month()
			r, present := months[int(month)]
			if present {
				r.Count++
			} else {
				r = Record{Count: 1}
			}
			r.Month = month
			months[int(month)] = r
		}
	}
	return months
}

// GetDays - get a map of days containing posts, and the post meta data for a given year, month
func GetDays(year int, month time.Month) map[int]Record {
	days := make(map[int]Record)
	pages := page.GetAllPages()
	for _, p := range pages {
		if p.PubDate.Year() == year && p.PubDate.Month() == month && p.Status == page.Published && p.Template == "post" {
			day := p.PubDate.Day()
			r, present := days[day]
			if present {
				r.Count++
			} else {
				r = Record{Count: 1}
			}
			r.Posts = append(r.Posts, p)
			days[day] = r
		}
	}
	return days
}

// GetKeysInOrder sorts map keys in order since ranging over a map is done in random order
func GetKeysInOrder(items map[int]Record) []int {
	// Make slice to store keys
	keys := make([]int, len(items))
	// Add the map keys to the slice
	i := 0
	for k := range items {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}
