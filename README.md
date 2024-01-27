# WSV-golang

This package is a simple implementation of a parser for White Space Separated Values (WSV) in Go as described in [the WSV Guide](https://dev.stenway.com/WSV/Index.html).

## Installation

```bash
go get github.com/gami13/wsv-golang
```

## Usage

### Parsing

```go
package main

import (
	"github.com/gami13/wsv-golang"
)

func main() {
	result, err := wsv.ParseDocument("a b c\n1 2 3")
	if err != nil {
		//handle error
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
	//handle error
}
result, err := wsv.ParseDocument(string(file))
if err != nil {
	//handle error
}
```

### Serialization

```go
package main

import (
	"github.com/gami13/wsv-golang"
)

func main() {
	//SERIALIZE DOCUMENT
	result := wsv.Serialize(input)

	//SERIALIZE ROW
	result := wsv.SerializeRow(input)
}
```

In this example result is a `string` , in the case of `Serialize`, the input is of type `[][]string` and in the case of `SerializeRow` the input is of type `[]string`.

## Future plans

- ~~Add serialization~~
