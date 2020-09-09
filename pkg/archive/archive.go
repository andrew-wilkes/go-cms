package archive

import (
	"gocms/pkg/page"
	"sort"
	"time"
)

// Record type that allows us to use a common index value for the various kinds of archive lists, be they in a map or a slice.
// Then it is easy to sort the order of the lists.
// Ranging over a map in Go is in a random order by design, and before, I used a mix of index types (int and time.Month).
type Record struct {
	Count int         // This represents the number of records in the current map
	Month time.Month  // This is only used when we want to record the related time.Month value
	Posts []page.Info // This is only used when we want to record a list of relevant pages
}

// GetYears - get a map of years containing posts with the post count
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

// GetKeysInOrder returns a sorted slice of map keys, since ranging over a map is done in a random order.
// We need this []int to index into the map in order to render ordered lists.
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
