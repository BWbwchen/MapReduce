# MapReduce

[<img alt="github" src="https://img.shields.io/badge/github-BWbwchen%2FMapReduce-blue?style=for-the-badge&logo=appveyor" height="20">](https://github.com/BWbwchen/MapReduce)
[<img alt="build status" src="https://img.shields.io/github/workflow/status/BWbwchen/MapReduce/Go/master?style=for-the-badge" height="20">](https://github.com/BWbwchen/MapReduce/actions?query=branch%3Amaster)
[![Go Reference](https://pkg.go.dev/badge/github.com/BWbwchen/MapReduce.svg)](https://pkg.go.dev/github.com/BWbwchen/MapReduce)

This is an easy-to-use Map Reduce Go framework inspired by [2021 6.824 lab1](http://nil.csail.mit.edu/6.824/2021/labs/lab-mr.html).

![mapReduce](https://github.com/kiarash8112/MapReduce/assets/133909368/03c6b149-213c-4906-91b7-a05c4f083c9e)



## Feature
- Multiple workers goroutine in a program on a single machine.
- Multiple workers process in separate program on a single machine.
- Fault tolerance.
- Easy to parallel your code with just Map and Reduce function. 


## Library Usage - Your own map and reduce function
Here's a simply example for word count program.
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

### Usage - 1 program with master, worker goroutine

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

Run with :
```
# Compile plugin
go build -race -buildmode=plugin -o wc.so wc.go

# Word count
go run -race main.go -i 'input/files' -p 'wc.so' -r 1 -w 8
```

Output file name is `mr-out-0.txt`

More example can be found in the [`mrapps/`](mrapps/) folder, and we will add more example in the future.

## Usage - Master program, and worker program (Isolate master and workers)

master.go
```golang
package main

import (
	mp "github.com/BWbwchen/MapReduce"
)

func main() {
	mp.StartMaster(mp.ParseArg())
}
```

worker.go
```golang
package main

import (
	mp "github.com/BWbwchen/MapReduce"
)

func main() {
	mp.StartWorker(mp.ParseArg())
}
```

Run with :
```
# Compile plugin
go build -race -buildmode=plugin -o wc.so wc.go

# Word count
go run -race cmd/master.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 8 &
sleep 1
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 1 &
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 2 &
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 3 &
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 4 &
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 5 &
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 6 &
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 7 &
go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 1 -w 8 
```



## Help
```
MapReudce is an easy-to-use Map Reduce Go parallel-computing framework inspired by 2021 6.824 lab1.
It supports multiple workers threads on a single machine and multiple processes on a single machine right now.

Usage:
  mapreduce [flags]

Flags:
  -h, --help            help for mapreduce
  -m, --inRAM           Whether write the intermediate file in RAM (default true)
  -i, --input strings   Input files
  -p, --plugin string   Plugin .so file
      --port int        Port number (default 10000)
  -r, --reduce int      Number of Reducers (default 1)
  -w, --worker int      Number of Workers(for master node)
                        ID of worker(for worker node) (default 4)
```

## Contributions
Pull requests are always welcome!


<sup>
Made by Bo-Wei Chen. All code is
licensed under the <a href="LICENSE">MIT License</a>.
</sup>
