package main

//
// a MapReduce pseudo-application that sometimes crashes,
// and sometimes takes a long time,
// to test MapReduce's ability to recover.
//
// go build -buildmode=plugin crash.go
//

import crand "crypto/rand"
import "math/big"
import "strings"
import "os"
import "sort"
import "strconv"
import "time"
import "github.com/BWbwchen/MapReduce/worker"

func maybeCrash() {
	max := big.NewInt(1000)
	rr, _ := crand.Int(crand.Reader, max)
	if rr.Int64() < 330 {
		// crash!
		os.Exit(1)
	} else if rr.Int64() < 660 {
		// delay for a while.
		maxms := big.NewInt(10 * 1000)
		ms, _ := crand.Int(crand.Reader, maxms)
		time.Sleep(time.Duration(ms.Int64()) * time.Millisecond)
	}
}

func Map(filename string, contents string, ctx worker.MrContext) {
	maybeCrash()

	ctx.EmitIntermediate("a", filename)
	ctx.EmitIntermediate("b", strconv.Itoa(len(filename)))
	ctx.EmitIntermediate("c", strconv.Itoa(len(contents)))
	ctx.EmitIntermediate("d", "xyzzy")
}

func Reduce(key string, values []string, ctx worker.MrContext) {
	maybeCrash()

	// sort values to ensure deterministic output.
	vv := make([]string, len(values))
	copy(vv, values)
	sort.Strings(vv)

	val := strings.Join(vv, " ")
	ctx.Emit(key, val)
}
