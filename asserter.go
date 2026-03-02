package asserter

import (
	"context"
	"fmt"
	"reflect"
)

type ErrorCode string

type FieldError struct {
	Field   string    `json:"field"`
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (f FieldError) Error() string {
	return fmt.Sprintf("%s: %s", f.Field, f.Message)
}

type Result struct {
	Success bool         `json:"success"`
	Errors  []FieldError `json:"errors"`
}

type predicate func(ctx context.Context, data any) *FieldError

type Assertion struct {
	field     string
	predicates []predicate
}

func Field(name string) *Assertion {
	return &Assertion{
		field: name,
	}
}

func (a *Assertion) add(pred predicate) *Assertion {
	a.predicates = append(a.predicates, pred)
	return a
}

func (a *Assertion) IsString() *Assertion {
	return a.add(func(ctx context.Context, data any) *FieldError {
		v := reflect.ValueOf(data)
		if v.Kind() != reflect.String {
			return &FieldError{
				Field:   a.field,
				Code:    "not_string",
				Message: "must be a string",
			}
		}
		return nil
	})
}

func (a *Assertion) IsInteger() *Assertion {
	return a.add(func(ctx context.Context, data any) *FieldError {
		v := reflect.ValueOf(data)
		kind := v.Kind()
		if kind != reflect.Int && kind != reflect.Int32 && kind != reflect.Int64 {
			return &FieldError{
				Field:   a.field,
				Code:    "not_integer",
				Message: "must be an integer",
			}
		}
		return nil
	})
}

func (a *Assertion) MinLength(min int) *Assertion {
	return a.add(func(ctx context.Context, data any) *FieldError {
		str, ok := data.(string)
		if !ok || len(str) < min {
			return &FieldError{
				Field:   a.field,
				Code:    "min_length",
				Message: fmt.Sprintf("must be at least %d characters", min),
			}
		}
		return nil
	})
}

func (a *Assertion) GreaterThan(min int) *Assertion {
	return a.add(func(ctx context.Context, data any) *FieldError {
		val, ok := data.(int)
		if !ok || val <= min {
			return &FieldError{
				Field:   a.field,
				Code:    "greater_than",
				Message: fmt.Sprintf("must be greater than %d", min),
			}
		}
		return nil
	})
}

func combine(preds ...predicate) predicate {
	return func(ctx context.Context, data any) *FieldError {
		for _, p := range preds {
			if fe := p(ctx, data); fe != nil {
				return fe
			}
		}
		return nil
	}
}

func Assert(ctx context.Context, data map[string]any, assertions ...*Assertion) Result {
	result := Result{ Success: true }
	for _, a := range assertions {
		select {
		case <-ctx.Done():
			result.Success = false
			result.Errors = append(result.Errors, FieldError{
				Field:   a.field,
				Code:    "context_cancelled",
				Message: "assertion cancelled",
			})
			return result
		default:
			value := data[a.field]
			if err := combine(a.predicates...)(ctx, value); err != nil {
				result.Success = false
				result.Errors = append(result.Errors, *err)
			}
		}
	}
	return result
}
