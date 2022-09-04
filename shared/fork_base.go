package shared

const UnOwned = 0 // Default value ensures Forks always initialize free

// ForkBase is a common element of all implementation Forks
type ForkBase struct {
	ID     int // The fork ID
	Holder int // Who holds the fork - UnOwned if free
}

// GetIS returns the Fork ID
func (f *ForkBase) GetID() int {
	return f.ID
}

// IsHeld indicates that the fork is held by someone
func (f *ForkBase) IsHeld() bool {
	return f.Holder != UnOwned
}

// IsHeldBy indicates that the fork is held by Philosopher id
func (f *ForkBase) IsHeldBy(id int) bool {
	// Philosophers index from 0, so increment to avoid the UnOwned value
	return f.Holder == id+1
}

// SetHolder sets the Fork owner
func (f *ForkBase) SetHolder(id int) {
	f.Holder = id + 1
}

// SetFree marks the fork as free
func (f *ForkBase) SetFree() {
	f.Holder = UnOwned
}
