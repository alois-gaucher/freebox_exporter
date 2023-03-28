package main

func getConnectionXdsl(authInf *authInfo, pr *postRequest, xSessionToken *string) (connectionXdsl, error) {
	connectionXdslResp := connectionXdsl{}
	err := getApiData(authInf, pr, xSessionToken, &connectionXdslResp, nil)
	if err != nil {
		return connectionXdsl{}, err
	}
	return connectionXdslResp, nil
}

func getDsl(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	return getRrdData(authInf, pr, xSessionToken, "dsl", []string{"rate_up", "rate_down", "snr_up", "snr_down"})
}

func getTemp(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	return getRrdData(authInf, pr, xSessionToken, "temp", []string{"cpum", "cpub", "sw", "hdd", "fan_speed"})
}

func getNet(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	return getRrdData(authInf, pr, xSessionToken, "net", []string{"bw_up", "bw_down", "rate_up", "rate_down", "vpn_rate_up", "vpn_rate_down"})
}

func getSwitch(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	return getRrdData(authInf, pr, xSessionToken, "switch", []string{"rx_1", "tx_1", "rx_2", "tx_2", "rx_3", "tx_3", "rx_4", "tx_4"})
}

func getLan(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]lanHost, error) {
	lanResp := lan{}
	err := getApiData(authInf, pr, xSessionToken, &lanResp, nil)
	if err != nil {
		return []lanHost{}, err
	}
	return lanResp.Result, nil
}

func getFreeplug(authInf *authInfo, pr *postRequest, xSessionToken *string) (freeplug, error) {
	if _, err := setFreeboxToken(authInf, xSessionToken); err != nil {
		return freeplug{}, err
	}

	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return freeplug{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return freeplug{}, err
	}
	if resp.StatusCode == 404 {
		return freeplug{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return freeplug{}, err
	}

	freeplugResp := freeplug{}
	err := getApiData(authInf, pr, xSessionToken, &freeplugResp, nil)
	if err != nil {
		return freeplug{}, err
	}
	return freeplugResp, nil
}

func getSystem(authInf *authInfo, pr *postRequest, xSessionToken *string) (system, error) {
	systemResp := system{}
	err := getApiData(authInf, pr, xSessionToken, &systemResp, nil)
	if err != nil {
		return system{}, err
	}
	return systemResp, nil
}

func getSystemV6(authInf *authInfo, pr *postRequest, xSessionToken *string) (systemV6, error) {
	systemResp := systemV6{}
	err := getApiData(authInf, pr, xSessionToken, &systemResp, nil)
	if err != nil {
		return systemV6{}, err
	}
	return systemResp, nil
}

func getWifi(authInf *authInfo, pr *postRequest, xSessionToken *string) (wifi, error) {
	wifiResp := wifi{}
	err := getApiData(authInf, pr, xSessionToken, &wifiResp, nil)
	if err != nil {
		return wifi{}, err
	}
	return wifiResp, nil
}

func getWifiStations(authInf *authInfo, pr *postRequest, xSessionToken *string) (wifiStations, error) {
	wifiStationResp := wifiStations{}
	err := getApiData(authInf, pr, xSessionToken, &wifiStationResp, nil)
	if err != nil {
		return wifiStations{}, err
	}
	return wifiStationResp, nil
}

func getVpnServer(authInf *authInfo, pr *postRequest, xSessionToken *string) (vpnServer, error) {
	vpnServerResp := vpnServer{}
	err := getApiData(authInf, pr, xSessionToken, &vpnServerResp, nil)
	if err != nil {
		return vpnServer{}, err
	}
	return vpnServerResp, nil
}

func getSwitchStatus(authInf *authInfo, pr *postRequest, xSessionToken *string) (switchStatus, error) {
	switchStatusResp := switchStatus{}
	err := getApiData(authInf, pr, xSessionToken, &switchStatusResp, nil)
	if err != nil {
		return switchStatus{}, err
	}
	return switchStatusResp, nil
}

func getSwitchPort(authInf *authInfo, pr *postRequest, xSessionToken *string) (switchPort, error) {
	switchPortResp := switchPort{}
	err := getApiData(authInf, pr, xSessionToken, &switchPortResp, nil)
	if err != nil {
		return switchPort{}, err
	}
	return switchPortResp, nil
}
