package pfsense

import (
	"net/netip"

	"github.com/jpmchia/ip2location-pfsense/util"
)

func DetermineIp(LogEntry LogEntry) (string, string) {
	var ip, direction string

	src, err1 := netip.ParseAddr(LogEntry.Srcip)
	dst, err2 := netip.ParseAddr(LogEntry.Dstip)

	if err1 != nil && err2 != nil {
		util.HandleError(err1, "[pfsense] DetermineIp:  Failed to parse IP address: %s %s", err1, err2)
		return "", ""
	}

	// If both the source and destination IP addresses are private IP addresses
	if src.IsPrivate() && dst.IsPrivate() {
		util.LogDebug("[pfsense] DetermineIp:  Both IP addresses are private")
		return "", ""
	}

	// If the source IP address is not a private IP address and the destination IP address is
	if !src.IsPrivate() && dst.IsPrivate() {
		if LogEntry.Direction == "in" {
			ip = LogEntry.Srcip
			direction = "in"
		} else {
			ip = LogEntry.Srcip
			direction = "[out]"
		}
		util.LogDebug("[pfsense] DetermineIp:  Source address is not private, destination address is private")

		return ip, direction
	}

	// If the source IP address is a private IP address and the destination IP address is not
	if src.IsPrivate() && !dst.IsPrivate() {
		if LogEntry.Direction == "out" {
			ip = LogEntry.Dstip
			direction = "out"
		} else {
			ip = LogEntry.Dstip
			direction = "[in]"
		}
		util.LogDebug("[pfsense] DetermineIp:  Source address is private, desitination address is not private")
		return ip, direction
	}

	if !src.IsPrivate() && !dst.IsPrivate() {
		if LogEntry.Direction == "out" {
			ip = LogEntry.Dstip
			direction = "out"
		} else {
			ip = LogEntry.Srcip
			direction = "in"
		}
		util.LogDebug("[pfsense] DetermineIp:  Both addresses are not private")
		return ip, direction
	}
	ip = ""
	direction = ""

	util.LogDebug("[pfsense] DetermineIp:  IP addresses unmatched")
	return ip, direction
}
