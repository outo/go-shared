package pops

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
