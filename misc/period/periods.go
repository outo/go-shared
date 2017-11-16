package period

import (
	"time"
	"sort"
)

type Periods struct {
	ps []Period
}

func CreatePeriods(periods []Period) Periods {
	this := Periods{}
	this.ps = append([]Period{}, periods...)
	return this
}

func appendTimedBoundaries(inputTimedBoundaries []timedBoundary, periods []Period, impact int) (outputTimedBoundaries []timedBoundary) {
	outputTimedBoundaries = append(inputTimedBoundaries)
	for _, p := range periods {
		outputTimedBoundaries = append(outputTimedBoundaries,
			timedBoundary{time: p.startIncl, isStart: true, impact: impact},
			timedBoundary{time: p.endExcl, isStart: false, impact: impact})
	}
	return
}

func stitchAdjacentPeriods(periods []Period) (stitched []Period, err error) {

	for pos := 0; pos < len(periods); pos++ {
		lowerBoundary := periods[pos].startIncl
		upperBoundary := periods[pos].endExcl
		for scan := pos + 1; scan < len(periods); scan++ {
			periodToMerge := periods[scan]
			if periodToMerge.startIncl.Equal(upperBoundary) {
				upperBoundary = periodToMerge.endExcl
				pos ++
				continue
			}
			break
		}
		newPeriod, err := CreatePeriod(lowerBoundary, upperBoundary)
		if err != nil {
			return nil, err
		}
		stitched = append(stitched, newPeriod)
	}
	return
}


func apply(t, o []Period, tImpact, oImpact int) (result Periods, err error) {
	tbs := make([]timedBoundary, len(o) + len(t))
	tbs = appendTimedBoundaries(tbs, o, oImpact)
	tbs = appendTimedBoundaries(tbs, t, tImpact)
	sort.Sort(timedBoundaries(tbs))

	var newPeriods []Period
	var currentPeriodStart time.Time
	previousImpact, currentImpact := 0, 0
	for _, tb := range tbs {
		if tb.isStart {
			currentImpact += tb.impact
		} else {
			currentImpact -= tb.impact
		}

		if currentImpact > 0 && previousImpact == 0 {
			currentPeriodStart = tb.time
		}

		if currentImpact == 0 && previousImpact > 0 && !currentPeriodStart.Equal(tb.time){
			period, err := CreatePeriod(currentPeriodStart, tb.time)
			if err != nil {
				return Periods{}, err
			}
			newPeriods = append(newPeriods, period)
		}
		previousImpact = currentImpact
	}

	ps, err := stitchAdjacentPeriods(newPeriods)
	stitchedNewPeriods := CreatePeriods(ps)
	return stitchedNewPeriods, err

}

func (periods Periods) Subtract(o Periods) (result Periods, err error) {
	return apply(periods.ps, o.ps, 1, -1)
}

func (periods Periods) Union(o Periods) (result Periods, err error) {
	return apply(periods.ps, o.ps, 1, 1)
}

