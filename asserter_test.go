package asserter

import (
	"context"
	"testing"
)

func TestAssertions(t *testing.T) {
	ctx := context.Background()

	data := map[string]any{
		"name": "alex",
		"age":  20,
	}

	res := Assert(ctx, data,
		Field("name").IsString().MinLength(3),
		Field("age").IsInteger().GreaterThan(18),
	)

	if !res.Success {
		t.Fatalf("expected success, got errors: %v", res.Errors)
	}

	res2 := Assert(ctx, data,
		Field("name").IsString().MinLength(10),
	)

	if res2.Success {
		t.Fatal("expected failure for name too short")
	}
}
