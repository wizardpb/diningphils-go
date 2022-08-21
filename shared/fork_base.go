package shared

const UnOwned = -1

type ForkBase struct {
	ID    int
	Owner int
}

func (f *ForkBase) IsOwned() bool {
	return f.Owner != UnOwned
}

func (f *ForkBase) IsOwnedBy(id int) bool {
	return f.Owner == id
}

func (f *ForkBase) SetOwner(id int) {
	f.Owner = id
}

func (f *ForkBase) SetUnowned() {
	f.Owner = UnOwned
}
