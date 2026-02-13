package validx

type Builder[T any] struct {
	schema *Schema[T]
}

func New[T any]() *Builder[T] {
	return &Builder[T]{
		schema: &Schema[T]{},
	}
}

func (b *Builder[T]) Schema() *Schema[T] {
	return b.schema
}
