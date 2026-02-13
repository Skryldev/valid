package validx

import "sync"

type schemaField[T any] interface {
	validate(obj *T, errs *Errors, mutex *sync.Mutex)
}

type fieldRule[T any, F any] struct {
	name   string
	getter func(*T) F
	rules  []Rule[F]
}

type FieldBuilder[T any, F any] struct {
	schema *Schema[T]
	field  *fieldRule[T, F]
}

func (f *fieldRule[T, F]) validate(obj *T, errs *Errors, mutex *sync.Mutex) {
	val := f.getter(obj)

	for _, rule := range f.rules {
		if rule == nil {
			mutex.Lock()
			errs.Add(f.name, "invalid_rule", "rule is not configured")
			mutex.Unlock()
			continue
		}

		ok, ruleErr := rule(val)
		if ok {
			continue
		}

		code := "invalid"
		msg := "validation failed"
		if ruleErr != nil {
			if ruleErr.Code != "" {
				code = ruleErr.Code
			}
			if ruleErr.Message != "" {
				msg = ruleErr.Message
			}
		}

		mutex.Lock()
		errs.Add(f.name, code, msg)
		mutex.Unlock()
	}
}

// ✅ Generic function (کاملاً قانونی)
func Field[T any, F any](
	s *Schema[T],
	name string,
	getter func(*T) F,
) *FieldBuilder[T, F] {

	fr := &fieldRule[T, F]{
		name:   name,
		getter: getter,
	}

	s.fields = append(s.fields, fr)

	return &FieldBuilder[T, F]{
		schema: s,
		field:  fr,
	}
}

func (f *FieldBuilder[T, F]) Use(rule Rule[F]) *FieldBuilder[T, F] {
	f.field.rules = append(f.field.rules, rule)
	return f
}
