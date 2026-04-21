package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/foliv/introduction-to-parallel-programing-collections/src/go/ch2"
)

func usage(prog string) {
	fmt.Fprintf(os.Stderr, "usage: %s <bin_count> <min_meas> <max_meas> <data_count>\n", prog)
	os.Exit(2)
}

func main() {
	if len(os.Args) != 5 {
		usage(os.Args[0])
	}

	binCount, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage(os.Args[0])
	}
	minMeas, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		usage(os.Args[0])
	}
	maxMeas, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		usage(os.Args[0])
	}
	dataCount, err := strconv.Atoi(os.Args[4])
	if err != nil {
		usage(os.Args[0])
	}

	data, err := ch2.GenData(minMeas, maxMeas, dataCount, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	h, err := ch2.BuildHistogramParallel(data, binCount, minMeas, maxMeas, runtime.NumCPU())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	h.Print(os.Stdout)
}
