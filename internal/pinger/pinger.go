package pinger

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	network string
	l_addr string
	recv_proto int
	icmp_type icmp.Type
)

type IcmpReply struct {
	Peer net.Addr
	Type icmp.Type
	Checksum int
	Body icmp.MessageBody
	IcmpProto int
	Duration time.Duration
}

func lookupAddress(s string) ([]net.IP, error) {
	result, err := net.LookupIP(s)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type ip_version int
func determineAddressFamily(s string) (ip_version) {
	addr := net.ParseIP(s)

	switch {
		case net.IP(addr).To4() != nil:
			return 4
		case net.IP(addr).To16() != nil:
			return 6
		default:
			return -1
	}
}

func SendICMPEcho(addr string, ttl int) (*IcmpReply, error) {
	var (
		network string
		l_addr string
		recv_proto int
		icmp_type icmp.Type
		remote_addr net.IPAddr
	)

	switch {
		case determineAddressFamily(addr) == 4:
			network, l_addr, recv_proto, icmp_type = "ip4:icmp", "0.0.0.0", 1, ipv4.ICMPTypeEcho
			remote_addr = net.IPAddr{IP: net.ParseIP(addr)}
		case determineAddressFamily(addr) == 6:
			network, l_addr, recv_proto, icmp_type = "ip6:ipv6-icmp", "::", 58, ipv6.ICMPTypeEchoRequest
			remote_addr = net.IPAddr{IP: net.ParseIP(addr)}
		default:
			// do we even wanna do lookups here?
			// might be better outside of this function
			if result, err := net.ResolveIPAddr("ip4", addr); err != nil {
				return &IcmpReply{}, errors.New(fmt.Sprintf("lookup error %v: %v", addr, err))
			} else {
				remote_addr = *result
			}
			network, l_addr, recv_proto, icmp_type = "ip4:icmp", "0.0.0.0", 1, ipv4.ICMPTypeEcho
	}

	conn, err := icmp.ListenPacket(network, l_addr)
	if err != nil {
		return &IcmpReply{}, errors.New(fmt.Sprintf("icmp.ListenPacket: %+v", err))
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(time.Millisecond*500))
	if err != nil {
		return &IcmpReply{}, errors.New(fmt.Sprintf("Setting connection deadline: %+v", err))
	}

	switch icmp_type {
		case ipv4.ICMPTypeEcho:
			err = conn.IPv4PacketConn().SetTTL(ttl)
		case ipv6.ICMPTypeEchoRequest:
			err = conn.IPv6PacketConn().SetHopLimit(ttl)
	}

	if err != nil {
		return &IcmpReply{}, errors.New(fmt.Sprintf("Setting TTL: %+v", err))
	}

	message := icmp.Message{
		Type: icmp_type , Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("mping request body"),
		},
	}

	msg, err := message.Marshal(nil)
	if err != nil {
		return &IcmpReply{}, errors.New(fmt.Sprintf("Marshalling Message: %+v", err))
	}
	
	start := time.Now()
	//if _, err := conn.WriteTo(msg, &net.UDPAddr{IP: net.ParseIP(addr), Zone: ""}); err != nil {
	if _, err := conn.WriteTo(msg, &remote_addr); err != nil {
		return &IcmpReply{}, errors.New(fmt.Sprintf("Writing to conn: %+v", err))
	}

	read_buff := make([]byte, 1500)
	n_bytes, peer, err := conn.ReadFrom(read_buff)
	if err != nil {
		return &IcmpReply{}, errors.New(fmt.Sprintf("Reading from conn: Err:%+v", err))
	}

	duration := time.Since(start)

	// recv_proto:
	// icmpv4 proto number -> 0x01 -> 1
	// icmpv6 proto number -> 0x3a -> 58
	// e.q: rm, err := icmp.ParseMessage(58, rb[:n])
	rm, err := icmp.ParseMessage(recv_proto, read_buff[:n_bytes])
	if err != nil {
		return &IcmpReply{}, errors.New(fmt.Sprintf("parsing message: %+v", err))
	}

	return &IcmpReply{
		Peer: peer, Type: rm.Type, Checksum: rm.Checksum,
		Body: rm.Body, IcmpProto: recv_proto, Duration: duration,
	}, nil
}
