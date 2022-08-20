package entry

import "time"

type Entry struct {
	URL   string
	Title string
	Date  time.Time
}
