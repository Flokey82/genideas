package simsettlers

// Opinions is a map of people to opinions.
type Opinions map[*Person][2]int

// Value returns the opinion value of the person.
func (o Opinions) Value(p *Person) int {
	return o[p][0]
}

// Counter returns the number of times the person has been evaluated.
func (o Opinions) Counter(p *Person) int {
	return o[p][1]
}

// IncrementBy increments the opinion of the person by the given value.
func (o Opinions) IncrementBy(p *Person, value int) {
	o[p] = [2]int{
		min(max(o[p][0]+value, 0), 255),
		min(o[p][1]+1, 255),
	}
}

// Change changes the opinion of the person by the given value (as running average).
func (o Opinions) Change(p *Person, value int) {
	o[p] = [2]int{min(max((o[p][0]*o[p][1]+value)/(o[p][1]+1), 0), 255), min(o[p][1]+1, 255)}
}
