package pops

import (
	"time"
	"errors"
)

type Period struct {
	startIncl,
	endExcl time.Time
}

func NewPeriod(startIncl, endExcl time.Time) (Period, error) {
	if !endExcl.After(startIncl) {
		return Period{}, errors.New("endExcl has to be after startIncl (equality also not allowed)")
	}
	return Period{
		startIncl: startIncl,
		endExcl:   endExcl,
	}, nil
}

func (period Period) Check(startIncl time.Time, endExcl time.Time) (periodRelationship Relationship) {
	lower := period.startIncl
	upper := period.endExcl

	if !startIncl.Before(lower) && startIncl.Before(upper) && endExcl.After(upper) {
		periodRelationship = OverlappingUpperEnd
	} else if startIncl.After(upper) {
		periodRelationship = DisparateAndHigher
	} else if startIncl.Equal(upper) {
		periodRelationship = AdjacentAndHigher
	} else if startIncl.After(lower) && endExcl.Before(upper) {
		periodRelationship = Contained
	} else if startIncl.Before(lower) && endExcl.After(upper) {
		periodRelationship = Containing
	} else if startIncl.Equal(lower) && endExcl.Equal(upper) {
		periodRelationship = Same
	} else if startIncl.Before(lower) && endExcl.After(lower) {
		periodRelationship = OverlappingLowerEnd
	} else if startIncl.Before(lower) && endExcl.Before(lower) {
		periodRelationship = DisparateAndLower
	} else if startIncl.Before(lower) && endExcl.Equal(lower) {
		periodRelationship = AdjacentAndLower
	}
	return
}

func (period Period) Get() (startIncl, endExcl time.Time) {
	return period.startIncl, period.endExcl
}

func (period Period) GetStartIncl() time.Time {
	return period.startIncl
}

func (period Period) GetEndExcl() time.Time {
	return period.endExcl
}

