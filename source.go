package radl

type Source interface {
	Title() string
	Id() string
	GetShow(id string) (Show, error)
}
