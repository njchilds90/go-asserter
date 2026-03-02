# go-asserter

Declarative assertion and predicate builder for Go.

## Features

- Zero runtime dependencies  
- Chainable assertions  
- Structured machine-readable error outcomes  
- Context-aware evaluation  
- Clean functional API surface

## Installation

```bash
go get github.com/njchilds90/go-asserter

## Example

import (
  "context"
  "fmt"
  "github.com/njchilds90/go-asserter"
)

data := map[string]any{
  "name": "alex",
  "age":  20,
}

result := asserter.Assert(context.Background(), data,
  asserter.Field("name").IsString().MinLength(3),
  asserter.Field("age").IsInteger().GreaterThan(18),
)

if !result.Success {
  fmt.Println("Validation errors:", result.Errors)
}
