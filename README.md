# MapReduce

[<img alt="github" src="https://img.shields.io/badge/github-BWbwchen%2FMapReduce-blue?style=for-the-badge&logo=appveyor" height="20">](https://github.com/BWbwchen/MapReduce)
[<img alt="build status" src="https://img.shields.io/github/workflow/status/BWbwchen/MapReduce/Go/master?style=for-the-badge" height="20">](https://github.com/BWbwchen/MapReduce/actions?query=branch%3Amaster)



This is an easy-to-use Map Reduce Go framework inspired by [2021 6.824 lab1](http://nil.csail.mit.edu/6.824/2021/labs/lab-mr.html).

## Feature
- Multiple workers on single machine right now.
- Easy to parallel your code with just Map and Reduce function. 

## Usage

Here's a simply example for word count program.

main.go
```golang
package main

import (
	mp "github.com/BWbwchen/MapReduce"
)

func main() {
	mp.StartSingleMachineJob(mp.ParseArg())
}
```

And write down your own Map and Reduce function.

wc.go
```golang
package main
import (
	"strconv"
	"strings"
	"unicode"

	"github.com/BWbwchen/MapReduce/worker"
)
func Map(filename string, contents string, ctx worker.MrContext) {
	// function to detect word separators.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	// split contents into an array of words.
	words := strings.FieldsFunc(contents, ff)

	for _, w := range words {
		ctx.EmitIntermediate(w, "1")
	}
}
func Reduce(key string, values []string, ctx worker.MrContext) {
	// return the number of occurrences of this word.
	ctx.Emit(key, strconv.Itoa(len(values)))
}
```

Run with :
```
go build -race -buildmode=plugin -o wc.so wc.go
go run -race main.go -i 'txt/*' -p 'wc.so' -r 1 -w 8
```

Output file name is `mr-out-0.txt`

More example can be found in the [`mrapps/`](mrapps/) folder, and we will add more example in the future.


<sup>
Made by Bo-Wei Chen. All code is
licensed under the <a href="LICENSE">MIT License</a>.
</sup>
