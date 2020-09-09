package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/lldp"
	"github.com/mdlayher/raw"
)

// TLVNames TLV ID to Name convertion
var TLVNames map[int]string = map[int]string{
	4: "Port Desc",
	5: "System Name",
	6: "System Desc",
	8: "Mngm Addr",
}

func scanLLDP(iface *net.Interface, timeout time.Duration, promiscuous bool) (string, error) {
	var output string
	c, err := raw.ListenPacket(iface, uint16(lldp.EtherType), &raw.Config{})
	if err != nil {
		return "", err
	}
	defer c.Close()

	if promiscuous {
		c.SetPromiscuous(true)
	}

	err = c.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return "", err
	}

	b := make([]byte, iface.MTU)
	var f ethernet.Frame
	var l lldp.Frame

	n, addr, err := c.ReadFrom(b)
	if promiscuous {
		c.SetPromiscuous(false)
	}
	if err != nil {
		return "", err
	}

	if err := (&f).UnmarshalBinary(b[:n]); err != nil {
		return "", err
	}

	if err := (&l).UnmarshalBinary(f.Payload); err != nil {
		return "", err
	}

	output = fmt.Sprintf("- Iface: %s\n", iface.Name)
	output += fmt.Sprintf("  src MAC: %s\n", addr)
	output += fmt.Sprintf("  ChassisID: %x\n", l.ChassisID.ID)
	output += fmt.Sprintf("  PortID: %s\n", l.PortID.ID)

	for _, tlv := range l.Optional {
		if name, ok := TLVNames[int(tlv.Type)]; ok {
			output += fmt.Sprintf("  %s: %s\n", name, tlv.Value)
		}
	}
	return output, nil
}

func main() {
	ifs := make([]net.Interface, 0)
	var wg sync.WaitGroup

	tFlag := flag.Int("t", 60, "timeout in seconds")
	eFlag := flag.Bool("e", false, "show errors")
	loFlag := flag.Bool("lo", false, "allow loopback")
	promFlag := flag.Bool("nop", false, "disable promiscuous interface mode")
	flag.Parse()
	ifNames := flag.Args()

	if len(ifNames) == 0 {
		ifNames = append(ifNames, "all")
	}

	for _, ifName := range ifNames {
		if ifName == "all" {
			ifs, _ = net.Interfaces()
			break
		}

		iface, err := net.InterfaceByName(ifName)
		if err != nil {
			log.Fatalf("unkown interface: %v", err)
		}

		ifs = append(ifs, *iface)
	}

	timeout := time.Duration(*tFlag) * time.Second

	for _, i := range ifs {
		if i.Name == "lo" && !*loFlag {
			continue
		}
		wg.Add(1)
		go func(i net.Interface, t time.Duration) {
			defer wg.Done()
			output, err := scanLLDP(&i, t, !*promFlag)
			if err != nil && *eFlag {
				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					fmt.Printf("- %s: - timeout (%d)\n", i.Name, t)
					return
				}
				fmt.Printf("- %s: - %v\n", i.Name, err)
				return
			}
			fmt.Print(output)

		}(i, timeout)
	}
	wg.Wait()
}
