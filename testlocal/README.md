# Local Testing

You may use the contents of the current directory to experiment and test the library yourself.

## Getting Started

1) On your machine, clone the project directory or download the zip.

```
git clone https://github.com/globocom/go-m3u8.git
```

2) Access the project root directory from your terminal, and then enter the `testlocal` folder.

``` 
cd go-m3u8/testlocal/
```

3) To test the library, write your code to the `main.go` file.

For example, we can parse the m3u8 manifest `multivariant.m3u8` and print the resulting playlist, like below.

```go
package main

import (
	"os"

	go_m3u8 "github.com/globocom/go-m3u8"
)

func main() {
	file, _ := os.Open("multivariant.m3u8")
	p, err := go_m3u8.ParsePlaylist(file)

	if err != nil {
		panic(err)
	}

	p.Print()
}
```

4) Then, use the following Makefile commands to run it:

```sh
make run      # Run the main.go code
make output   # Run the main.go code and writes output to local output.txt file
```

