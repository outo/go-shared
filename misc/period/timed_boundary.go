package period

import "time"

type timedBoundary struct {
	time    time.Time
	isStart bool
	impact  int
}

type timedBoundaries []timedBoundary

func (a timedBoundaries) Len() int           { return len(a) }
func (a timedBoundaries) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a timedBoundaries) Less(i, j int) bool { return a[i].time.Before(a[j].time) }
