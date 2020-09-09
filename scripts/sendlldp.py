from scapy.all import *

load_contrib("lldp")

l='00:01:02:03:04:05'
c=LLDPDUChassisID.SUBTYPE_MAC_ADDRESS
p=LLDPDUPortID.SUBTYPE_INTERFACE_NAME
frm = Ether(dst=LLDP_NEAREST_BRIDGE_MAC)/  \
    LLDPDUChassisID(subtype=c, id=l)/  \
    LLDPDUPortID(subtype=p, id="gi1/0/30")/  \
    LLDPDUTimeToLive(ttl=120)/ \
    LLDPDUSystemName(system_name="switch-test1")/ \
    LLDPDUEndOfLLDPDU()

frm.show()

sendp(frm,iface="lo")
