package validx

type Error struct {
	Field   string
	Code    string
	Message string
}

type Errors struct {
	list []Error
}

func (e *Errors) Add(field, code, msg string) {
	e.list = append(e.list, Error{
		Field:   field,
		Code:    code,
		Message: msg,
	})
}

func (e *Errors) Has() bool {
	return len(e.list) > 0
}

func (e *Errors) All() []Error {
	return e.list
}
