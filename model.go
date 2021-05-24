package goparkruncrawler

import (
	"fmt"
	"net/url"
	"time"
)

// Parkrun represents the parkrun event
type Parkrun struct {
	Parkrun     string
	ParkrunLink url.URL
}

// ParkrunRun represents a separate race
type ParkrunRun struct {
	Parkrun         *Parkrun
	Date            time.Time
	GenderPosition  int64
	OverallPosition int64
	Time            time.Duration
	AgeGrade        float64
}

func (p ParkrunRun) String() string {
	return fmt.Sprintf("%s [%s %s]", p.Parkrun.Parkrun, p.Date.Format("02/01/2006"), p.Time.String())
}
