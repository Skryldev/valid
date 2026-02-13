package validx

type Schema[T any] struct {
	fields []schemaField[T]
}
