package shared

// Fork is an interface used to control and set fork state. It is implemented by each different algorithm
type Fork interface {
	GetID() int
	IsHeld() bool
	IsHeldBy(id int) bool
	SetHolder(id int)
	SetFree()
}
