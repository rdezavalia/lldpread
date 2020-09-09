# lldpread

lldpread is an lldp client that is able to read an lldp frame and display information about your network neighboars.

The easiest way to use it is to copy the binary to the host where you want to check network neighboars and run it:

    # ./lldpread
    - Iface: eth1
      src MAC: b1:07:4f:84:ca:4d
      ChassisID: 000102030405
      PortID: gi1/0/30
      System Name: switch-test1

if you want, you can choose on which interfaces to run:

    # ./lldpread eth1 eth2 eth5

or change some of the default options:

     # ./lldpread --help
     Usage of /home/rdz/repo/lldpread/lldpread:
       -e	show errors
       -lo
         	allow loopback
       -nop
        	disable promiscuous interface mode
       -t int
        	timeout in seconds (default 60)

## Using tshark
One easy way to get the same information is using tshark:

    # tshark -c 1 -V -n -i lo -f "ether proto 0x88cc"
      Frame 1: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface lo, id 0
      Interface id: 0 (lo)
          Interface name: lo
      Encapsulation type: Ethernet (1)
      Arrival Time: Sep  9, 2020 17:08:20.697126010 -03
      [Time shift for this packet: 0.000000000 seconds]
      Epoch Time: 1599682100.697126010 seconds
      [Time delta from previous captured frame: 0.000000000 seconds]
      [Time delta from previous displayed frame: 0.000000000 seconds]
      [Time since reference or first frame: 0.000000000 seconds]
      Frame Number: 1
      Frame Length: 60 bytes (480 bits)
      Capture Length: 60 bytes (480 bits)
      [Frame is marked: False]
      [Frame is ignored: False]
      [Protocols in frame: eth:ethertype:lldp]
         Ethernet II, Src: b8:08:cf:84:ce:4d, Dst: 01:80:c2:00:00:0e
         Destination: 01:80:c2:00:00:0e
         Address: 01:80:c2:00:00:0e
         .... ..0. .... .... .... .... = LG bit: Globally unique address (factory default)
         .... ...1 .... .... .... .... = IG bit: Group address (multicast/broadcast)
         Source: b8:08:cf:84:ce:4d
         Address: b8:08:cf:84:ce:4d
         .... ..0. .... .... .... .... = LG bit: Globally unique address (factory default)
         .... ...0 .... .... .... .... = IG bit: Individual address (unicast)
         Type: 802.1 Link Layer Discovery Protocol (LLDP) (0x88cc)
         Padding: 310000000000000000
     Link Layer Discovery Protocol
     Chassis Subtype = MAC address, Id: 00:01:02:03:04:05
        0000 001. .... .... = TLV Type: Chassis Id (1)
        .... ...0 0000 0111 = TLV Length: 7
        Chassis Id Subtype: MAC address (4)
        Chassis Id: 00:01:02:03:04:05
     Port Subtype = Interface name, Id: gi1/0/30
        0000 010. .... .... = TLV Type: Port Id (2)
        .... ...0 0000 1001 = TLV Length: 9
        Port Id Subtype: Interface name (5)
        Port Id: gi1/0/30
     Time To Live = 120 sec
        0000 011. .... .... = TLV Type: Time to Live (3)
        .... ...0 0000 0010 = TLV Length: 2
        Seconds: 120
     System Name = switch-test1
        0000 101. .... .... = TLV Type: System Name (5)
        .... ...0 0000 1100 = TLV Length: 12
        System Name: switch-test1
     End of LLDPDU
        0000 000. .... .... = TLV Type: End of LLDPDU (0)
        .... ...0 0000 0000 = TLV Length: 0

     1 packet captured
