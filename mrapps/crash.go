package main

//
// a MapReduce pseudo-application that sometimes crashes,
// and sometimes takes a long time,
// to test MapReduce's ability to recover.
//
// go build -buildmode=plugin crash.go
//

import (
	crand "crypto/rand"
	// "fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/BWbwchen/MapReduce/worker"
)

// import "time"

func maybeCrash() {
	max := big.NewInt(1000)
	rr, _ := crand.Int(crand.Reader, max)
	if rr.Int64() < 250 {
		// crash!
		// os.Exit(1)
		// } else if rr.Int64() < 660 {
		// delay for a while.
		// maxms := big.NewInt(10 * 1000)
		// ms, _ := crand.Int(crand.Reader, maxms)
		os.Exit(0)
		// time.Sleep(100 * time.Second)
	}
}

func Map(filename string, contents string, ctx worker.MrContext) {
	// maybeCrash()
	if filename == "/home/bwbwchen/work/mapreduce/txt/bw.txt" {
		maybeCrash()
	}

	ctx.EmitIntermediate("a", filename)
	ctx.EmitIntermediate("b", strconv.Itoa(len(filename)))
	ctx.EmitIntermediate("c", strconv.Itoa(len(contents)))
	ctx.EmitIntermediate("d", "xyzzy")
}

func Reduce(key string, values []string, ctx worker.MrContext) {
	// maybeCrash()

	// sort values to ensure deterministic output.
	vv := make([]string, len(values))
	copy(vv, values)
	sort.Strings(vv)

	val := strings.Join(vv, " ")
	ctx.Emit(key, val)
}
