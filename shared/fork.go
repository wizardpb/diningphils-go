package shared

type Fork interface {
	GetID() int
	IsHeld() bool
	IsHeldBy(id int) bool
	SetHolder(id int)
	SetFree()
}
