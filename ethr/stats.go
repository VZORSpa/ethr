// -----------------------------------------------------------------------------
// Copyright (C) Microsoft. All rights reserved.
// Licensed under the MIT license.
// See LICENSE.txt file in the project root for full license information.
// -----------------------------------------------------------------------------
package ethr

import (
	"net"
)

type EthrNetStat struct {
	NetDevStats []EthrNetDevStat
	tcpStats    ethrTCPStat
}

type EthrNetDevStat struct {
	InterfaceName string
	HwAddr        net.HardwareAddr
	RxBytes       uint64
	TxBytes       uint64
	RxPkts        uint64
	TxPkts        uint64
	TxErrPkts     uint64
	RxErrPkts     uint64
	TxDrops       uint64
	RxDrops       uint64
	Flags         net.Flags
}

type ethrTCPStat struct {
	segRetrans uint64
}

func getNetDevStatDiff(curStats EthrNetDevStat, prevNetStats EthrNetStat, seconds uint64) EthrNetDevStat {
	for _, prevStats := range prevNetStats.NetDevStats {
		if prevStats.InterfaceName != curStats.InterfaceName {
			continue
		}

		if curStats.RxBytes >= prevStats.RxBytes {
			curStats.RxBytes -= prevStats.RxBytes
		} else {
			curStats.RxBytes += (^uint64(0) - prevStats.RxBytes)
		}

		if curStats.TxBytes >= prevStats.TxBytes {
			curStats.TxBytes -= prevStats.TxBytes
		} else {
			curStats.TxBytes += (^uint64(0) - prevStats.TxBytes)
		}

		if curStats.RxPkts >= prevStats.RxPkts {
			curStats.RxPkts -= prevStats.RxPkts
		} else {
			curStats.RxPkts += (^uint64(0) - prevStats.RxPkts)
		}

		if curStats.TxPkts >= prevStats.TxPkts {
			curStats.TxPkts -= prevStats.TxPkts
		} else {
			curStats.TxPkts += (^uint64(0) - prevStats.TxPkts)
		}

		break
	}
	curStats.RxBytes /= seconds
	curStats.TxBytes /= seconds
	curStats.RxPkts /= seconds
	curStats.TxPkts /= seconds
	return curStats
}
