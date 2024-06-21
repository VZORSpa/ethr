//go:build linux
// +build linux

// -----------------------------------------------------------------------------
// Copyright (C) Microsoft. All rights reserved.
// Licensed under the MIT license.
// See LICENSE.txt file in the project root for full license information.
// -----------------------------------------------------------------------------
package ethr

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type ethrNetDevInfo struct {
	bytes      uint64
	packets    uint64
	drop       uint64
	errs       uint64
	fifo       uint64
	frame      uint64
	compressed uint64
	multicast  uint64
}

type osStats struct {
}

func getNetDevStats(stats *ethrNetStat) {
	ifs, err := net.Interfaces()
	if err != nil {
		return
	}

	netStatsFile, err := os.Open("/proc/net/dev")
	if err != nil {
		return
	}
	defer netStatsFile.Close()

	reader := bufio.NewReader(netStatsFile)

	// Pass the header
	// Inter-|   Receive                                             |  Transmit
	//  face |bytes packets errs drop fifo frame compressed multicast|bytes packets errs drop fifo colls carrier compressed
	reader.ReadString('\n')
	reader.ReadString('\n')

	var line string
	for err == nil {
		line, err = reader.ReadString('\n')
		if line == "" {
			continue
		}
		netDevStat := buildNetDevStat(line)
		netDevStat.flags = getFlags(netDevStat.interfaceName, ifs)
		netDevStat.hwAddr = getHWAddr(netDevStat.interfaceName, ifs)
		if isIfUp(netDevStat.interfaceName, ifs) {
			stats.netDevStats = append(stats.netDevStats, buildNetDevStat(line))
		}
	}
}

func getTCPStats(stats *ethrNetStat) {
	snmpStatsFile, err := os.Open("/proc/net/snmp")
	if err != nil {
		return
	}
	defer snmpStatsFile.Close()

	reader := bufio.NewReader(snmpStatsFile)

	var line string
	for err == nil {
		// Tcp: RtoAlgorithm RtoMin RtoMax MaxConn ActiveOpens PassiveOpens AttemptFails EstabResets
		//      CurrEstab InSegs OutSegs RetransSegs InErrs OutRsts InCsumErrors
		line, err = reader.ReadString('\n')
		if line == "" || !strings.HasPrefix(line, "Tcp") {
			continue
		}
		// Skip the first line starting with Tcp
		line, err = reader.ReadString('\n')
		if !strings.HasPrefix(line, "Tcp") {
			break
		}
		fields := strings.Fields(line)
		stats.tcpStats.segRetrans = toUInt64(fields[12])
	}
}

func blockWindowResize() {
}

func buildNetDevStat(line string) ethrNetDevStat {
	fields := strings.Fields(line)
	if len(fields) < 17 {
		return ethrNetDevStat{}
	}
	interfaceName := strings.TrimSuffix(fields[0], ":")
	rxInfo := toNetDevInfo(fields[1:9])
	txInfo := toNetDevInfo(fields[9:17])
	return ethrNetDevStat{
		interfaceName: interfaceName,
		rxBytes:       rxInfo.bytes,
		txBytes:       txInfo.bytes,
		rxPkts:        rxInfo.packets,
		txPkts:        txInfo.packets,
		txErrPkts:     txInfo.errs,
		rxErrPkts:     rxInfo.errs,
		rxDrops:       rxInfo.drop,
		txDrops:       txInfo.drop,
	}
}

func toNetDevInfo(fields []string) ethrNetDevInfo {
	return ethrNetDevInfo{
		bytes:      toUInt64(fields[0]),
		packets:    toUInt64(fields[1]),
		errs:       toUInt64(fields[2]),
		drop:       toUInt64(fields[3]),
		fifo:       toUInt64(fields[4]),
		frame:      toUInt64(fields[5]),
		compressed: toUInt64(fields[6]),
		multicast:  toUInt64(fields[7]),
	}
}

func toUInt64(str string) uint64 {
	res, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return res
}

func isIfUp(ifName string, ifs []net.Interface) bool {
	for _, ifi := range ifs {
		if ifi.Name == ifName {
			if (ifi.Flags & net.FlagUp) != 0 {
				return true
			}
			return false
		}
	}
	return false
}
func getFlags(ifName string, ifs []net.Interface) *net.Flags {
	for _, ifi := range ifs {
		if ifi.Name == ifName {
			return &ifi.Flags
		}
	}
	return 0
}
func getHWAddr(ifName string, ifs []net.Interface) net.HardwareAddr {
	for _, ifi := range ifs {
		if ifi.Name == ifName {

			return ifi.HardwareAddr
		}
	}
	return nil
}

func setSockOptInt(fd uintptr, level, opt, val int) (err error) {
	err = syscall.SetsockoptInt(int(fd), level, opt, val)
	if err != nil {
	}
	return
}

func IcmpNewConn(address string) (net.PacketConn, error) {
	dialedConn, err := net.Dial(Icmp(), address)
	if err != nil {
		return nil, err
	}
	localAddr := dialedConn.LocalAddr()
	dialedConn.Close()
	conn, err := net.ListenPacket(Icmp(), localAddr.String())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func VerifyPermissionForTest(testID EthrTestID) {
	if testID.Protocol == ICMP || (testID.Protocol == TCP &&
		(testID.Type == TraceRoute || testID.Type == MyTraceRoute)) {
		if !IsAdmin() {

		}
	}
}

func IsAdmin() bool {
	return os.Geteuid() == 0
}

func SetTClass(fd uintptr, tos int) {
	setSockOptInt(fd, syscall.IPPROTO_IPV6, syscall.IPV6_TCLASS, tos)
}
