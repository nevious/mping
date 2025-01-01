package pinger

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime"
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

func Ping(addr string) (string, error) {
	switch runtime.GOOS {
		case "darwin", "ios":
		case "linux":
			slog.Debug("you may need to adjust the net.ipv4.ping_group_range kernel state")
		default:
			return "", errors.New("Unsupported Operating system")
	}

	switch {
		case determineAddressFamily(addr) == 4:
			network, l_addr, recv_proto, icmp_type = "udp4", "0.0.0.0", 1, ipv4.ICMPTypeEcho
		case determineAddressFamily(addr) == 6:
			network, l_addr, recv_proto, icmp_type = "udp6", "::", 58, ipv6.ICMPTypeEchoRequest
		default:
			// most likely a hostname.. let's first try to look it up and then figure out the rest
			addr_slice, err := lookupAddress(addr)
			if err != nil {
				return "", errors.New(fmt.Sprintf("lookup error %v: %v", addr, err))
			}
			new_addr := addr_slice[0].String()

			if determineAddressFamily(new_addr) == 4 {
				network, l_addr, recv_proto, icmp_type = "udp4", "0.0.0.0", 1, ipv4.ICMPTypeEcho
			} else if determineAddressFamily(new_addr) == 6 {
				network, l_addr, recv_proto, icmp_type = "udp6", "::", 58, ipv6.ICMPTypeEchoRequest
			}

			// Keeping it here just in case...
			slog.Debug(
				fmt.Sprintf("default mapping performend on %v", addr),
				slog.String("network", network),
				slog.String("l_addr", l_addr),
				slog.Int("recv_proto", recv_proto),
				slog.String("icmp_type", fmt.Sprintf("%d", icmp_type)),
				slog.String("name lookup", fmt.Sprintf("%v", addr_slice)),
				slog.String("new addr", new_addr),
			)

			addr = new_addr
	}

	conn, err := icmp.ListenPacket(network, l_addr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("icmp.ListenPacket: %+v", err))
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(time.Millisecond*500))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Setting connection deadline: %+v", err))
	}

	message := icmp.Message{
		Type: icmp_type , Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("ICMP echo request"),
		},
	}

	wb, err := message.Marshal(nil)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Marshalling Message: %+v", err))
	}
	
	if _, err := conn.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(addr), Zone: ""}); err != nil {
		return "", errors.New(fmt.Sprintf("Writing to conn: %+v", err))
	}

	rb := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Reading from conn: %+v", err))
	}

	// recv_proto:
	// icmpv4 proto number -> 0x01 -> 1
	// icmpv6 proto number -> 0x3a -> 58
	// e.q: rm, err := icmp.ParseMessage(58, rb[:n])
	rm, err := icmp.ParseMessage(recv_proto, rb[:n])
	if err != nil {
		return "", errors.New(fmt.Sprintf("parsing message: %+v", err))
	}

	switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			return fmt.Sprintf("Reply %v - Type: %d - Checksum: %d - Message: %v ", peer, rm.Type, rm.Checksum, rm), nil
		case ipv6.ICMPTypeEchoReply:
			return fmt.Sprintf("Reply %v - Type: %d - Checksum: %d - Message: %v ", peer, rm.Type, rm.Checksum, rm), nil
		default:
			return fmt.Sprintf("Non-Reply: %v - Type: %d - Checksum: %d - Raw: %v", peer, rm.Type, rm.Checksum, rm), nil
	}
}

