package validx

import (
	"sync"
)

func (s *Schema[T]) Validate(obj *T) *Errors {

	errs := &Errors{}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for _, f := range s.fields {
		wg.Add(1)

		go func(field schemaField[T]) {
			defer wg.Done()
			field.validate(obj, errs, &mutex)
		}(f)
	}

	wg.Wait()

	if errs.Has() {
		return errs
	}

	return nil
}
