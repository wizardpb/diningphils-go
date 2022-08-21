package philstate

// Enum is an integer state symbol
type Enum int

// Philosopher states
const (
	Inactive Enum = iota
	Thinking
	Hungry
	Eating
	Stopped
)

var vals = map[Enum]string{
	Inactive: "Inactive",
	Thinking: "Thinking",
	Hungry:   "Hungry",
	Eating:   "Eating",
	Stopped:  "Stopped",
}

func (e Enum) String() string {
	s, ok := vals[e]
	if !ok {
		return "Unknown"
	}
	return s
}
