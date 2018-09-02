package main

import (
	"github.com/riomhaire/lightauth2/frameworks/application/lightauth2/cmd"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// tracefile, err := os.Create("app.trace")
	// check(err)

	// pprof.StartCPUProfile(tracefile)
	//	trace.Start(tracefile)
	// Shutdown
	cmd.Execute()
}
