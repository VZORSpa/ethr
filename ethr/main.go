package ethr

func GetNetworkInfo() (*EthrNetStat, error) {
	out := EthrNetStat{}

	err := getNetDevStats(&out)
	if err != nil {
		return nil, err
	}

	return &out, err

}
