package pinger

import (
	"testing"
	"fmt"
	"net"
)

type input struct {
	addr string
	exp ip_version
}

var addresses = []input{
	input{"1.1.1.1", 4},
	input{"192.168.1.1", 4},
	input{"8.8.8.8", 4},
	input{"127.0.0.1", 4},
	input{"::1", 6},
	input{"2a00:1450:400a:802::200e", 6},
}

func TestDetermineAddressFamily(t *testing.T) {
	for _, test := range addresses {
		result:= determineAddressFamily(test.addr)
		if result == -1 {
			t.Errorf("input wrongly determined '%v' not %v", result, test.exp)
		}

		if result != test.exp {
			t.Errorf("input determined as '%v' not %v", result, test.exp)
		}
	}
}

var ping_addresses = []input{
	input{"2a00:1450:400a:802::200e", 6},
	input{"185.178.192.107", 4},
	input{"nevious.ch", 4},
	input{"google.com", 6},
}

func TestPingDestination(t *testing.T) {
	for _, test := range ping_addresses {
		result, err := Ping(test.addr)
		if err != nil {
			t.Errorf("threw error, unexpected: %+v", err)
		}

		fmt.Printf("Pinging result: %+v\n", result)
	}
}

type string_input struct {
	addr string
	exp net.IP
}

var lookup_addresses = []string_input {
	string_input{"example.com", net.ParseIP("2606:2800:21f:cb07:6820:80da:af6b:8b2c")},
	string_input{"nevious.ch", net.ParseIP("75.2.60.5")},
	string_input{"nevious_.ch", nil},
}

func TestLookupAddress(t *testing.T) {
	for _, test := range lookup_addresses {
		result, err := lookupAddress(test.addr)
		if err != nil && test.exp != nil {
			t.Errorf("%v threw error, unexpected: %+v\n", test.addr, err)
		} else if err != nil && test.exp == nil {
			t.Logf("%v expected: %v\n", test.addr, err)
			return
		}

		for _, element := range result {
			if element.String() == test.exp.String() {
				t.Logf("%v: %v Found in %v", test.addr, test.exp, result)
			}
		}

	}
}

