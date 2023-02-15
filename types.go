package motxn

type Int int

func (i Int) Less(than Int) bool {
	return i < than
}
