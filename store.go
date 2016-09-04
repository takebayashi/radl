package radl

type Store interface {
	Save(ShowPayload) error
}
