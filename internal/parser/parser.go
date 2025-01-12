package parser

import (
	"fmt"
	"os"
	"flag"
	"strings"
)

var (
	help bool
	version bool
	addrFlag hosts
)

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
	flag.Var(&addrFlag, "addresses", "Whitespace separated list of hosts to ping")
	flag.Var(&addrFlag, "a", "Whitespace separated list of hosts to ping")
	flag.Parse()

	if help != false {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if version != false {
		fmt.Println("No version yet...")
		os.Exit(0)
	}

	return addrFlag
}
