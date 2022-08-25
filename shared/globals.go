package shared

const (
	NPhils   int = 5
	ThinkMin     = 5
	ThinkMax     = 15
	EatMin       = 5
	EatMax       = 15
)

// Philosophers are numbers 0 to NPhils-1, as are forks.
//
// The fork to the left of Philosopher[i] is Fork[i]; the fork to the right is Fork[i+1 mod NPhils], since
// the wrap around the table
var (
	PhilNames    [NPhils]string = [NPhils]string{"Hannah Arendt", "Judith Butler", "Patricia Churchland", "Simone de Beauvoir", "Themistoclea"}
	Philosophers [NPhils]Philosopher
	Forks        [NPhils]Fork
)
