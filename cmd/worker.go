package main

import (
	mp "github.com/BWbwchen/MapReduce"
)

func main() {
	mp.StartWorker(mp.ParseArg())
}
