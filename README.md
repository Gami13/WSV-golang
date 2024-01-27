# WSV-golang

This package is a simple implementation of a parser for White Space Separated Values (WSV) in Go as described in [the WSV Guide](https://dev.stenway.com/WSV/Index.html).

## Installation

```bash
go get github.com/Gami13/wsv-golang
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/Gami13/wsv-golang"
)

func main() {
	result, err := wsv.ParseDocument("a b c\n1 2 3")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
```

In this example result is a `[][]string` with the following content:

```go
[["a" "b" "c"] ["1" "2" "3"]]
```

you can also use this to read from a file:

```go
	file, err := os.ReadFile("test_input.wsv")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	result, err := wsv.ParseDocument(string(file))
	if err != nil {
		t.Errorf("Error: %v", err)
	}
```

## Future plans

- Add serialization
