package parser

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

var (
	help bool
	version bool
	addrFlag hosts
)

func buildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprint(os.Stderr, "Unable to get version number!")
		os.Exit(1)
	}

	if info.Main.Version != "" {
		fmt.Printf("Version: %s\n", info.Main.Version)
	} else {
		fmt.Println("unknown version")
	}

	os.Exit(0)
}

type hosts []string

func (h *hosts) String() string {
	return fmt.Sprint(*h)
}

func (h *hosts) Set(value string) error {
	for _, element := range strings.Split(value, ",") {
		*h = append(*h, element)
	}

	return nil
}

func Parse() []string {
	flag.BoolVar(&help, "help", false, "Show help message and exit")
	flag.BoolVar(&version, "version", false, "Show version number and exit" )
	flag.BoolVar(&version, "v", false, "Show version number and exit" )
	flag.Var(&addrFlag, "addresses", "Whitespace separated list of hosts to ping")
	flag.Var(&addrFlag, "a", "Whitespace separated list of hosts to ping")
	flag.Parse()

	if help != false {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if version != false {
		buildInfo()
		os.Exit(0)
	}

	return addrFlag
}
