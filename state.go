package radl

type State interface {
	IsNew(Show) bool
	Update(Show) error
}
