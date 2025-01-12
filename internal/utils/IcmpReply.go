package utils

import (
	"net"
	"time"

	"golang.org/x/net/icmp"
)

type IcmpReply struct {
	Peer net.Addr
	Type icmp.Type
	Checksum int
	Body icmp.MessageBody
	IcmpProto int
	Duration time.Duration
}
