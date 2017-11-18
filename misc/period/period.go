package period

import (
	"time"
	"errors"
)

type Period struct {
	startIncl,
	endExcl time.Time
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

type Relationship int

const (
	Unknown Relationship = iota
	//given this period	(d, i) 	a b c [---------) j k l
	//the following are the examples of relationships of given period to the above
	DisparateAndLower    //		[-) c d e f g h i j k l
	AdjacentAndLower     //		[-----) e f g h i j k l
	OverlappingLowerEnd  //		a b [---) f g h i j k l
	Contained            //		a b c d e [---) i j k l
	Same                 //		a b c [---------) j k l
	OverlappingUpperEnd  //		a b c d e f [-------) l
	AdjacentAndHigher    //		a b c d e f g h [-----)
	DisparateAndHigher   //		a b c d e f g h i [---)
	Containing           //		[---------------------)
)

func CreatePeriod(startIncl, endExcl time.Time) (Period, error) {
	if !endExcl.After(startIncl) {
		return Period{}, errors.New("endExcl has to be after startIncl (equality also not allowed)")
	}
	return Period{
		startIncl: startIncl,
		endExcl:   endExcl,
	}, nil
}
