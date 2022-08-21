package shared

type Fork interface {
	IsOwned() bool
	IsOwnedBy(id int) bool
	SetOwner(id int)
	SetUnowned()
}
