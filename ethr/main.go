package ethr

func GetNetworkInfo() (*ethrNetStat, error) {
	out := ethrNetStat{}

	err := getNetDevStats(&out)
	if err != nil {
		return nil, err
	}

	return &out, err

}
