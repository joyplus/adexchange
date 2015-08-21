package lib

type HashSet struct {
	set map[string]bool
}

func NewHashSet() *HashSet {
	return &HashSet{make(map[string]bool)}
}

func (set *HashSet) Add(i string) bool {
	_, found := set.set[i]
	set.set[i] = true
	return !found //False if it existed already
}

func (set *HashSet) Get(i string) bool {
	_, found := set.set[i]
	return found //true if it existed already
}

func (set *HashSet) Remove(i string) {
	delete(set.set, i)
}
