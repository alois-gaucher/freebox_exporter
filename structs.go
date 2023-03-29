package main

import "bufio"

type apiResponse struct {
	Success   bool   `json:"success"`
	ErrorCode string `json:"error_code,omitempty"`
}

type track struct {
	apiResponse
	Result struct {
		AppToken string `json:"app_token"`
		TrackID  int    `json:"track_id"`
	} `json:"result"`
}

type grant struct {
	apiResponse
	Result struct {
		Status    string `json:"status"`
		Challenge string `json:"challenge"`
	} `json:"result"`
}

type challenge struct {
	apiResponse
	Result struct {
		LoggedIN  bool   `json:"logged_in,omitempty"`
		Challenge string `json:"challenge"`
	} `json:"result"`
}

type session struct {
	AppID    string `json:"app_id"`
	Password string `json:"password"`
}

type sessionToken struct {
	apiResponse
	Msg    string `json:"msg,omitempty"`
	UID    string `json:"uid,omitempty"`
	Result struct {
		SessionToken string `json:"session_token,omitempty"`
		Challenge    string `json:"challenge"`
		Permissions  struct {
			Settings   bool `json:"settings,omitempty"`
			Contacts   bool `json:"contacts,omitempty"`
			Calls      bool `json:"calls,omitempty"`
			Explorer   bool `json:"explorer,omitempty"`
			Downloader bool `json:"downloader,omitempty"`
			Parental   bool `json:"parental,omitempty"`
			Pvr        bool `json:"pvr,omitempty"`
			Home       bool `json:"home,omitempty"`
			Camera     bool `json:"camera,omitempty"`
		} `json:"permissions,omitempty"`
	} `json:"result"`
}

type rrd struct {
	apiResponse
	UID    string `json:"uid,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Result struct {
		DateStart int                `json:"date_start,omitempty"`
		DateEnd   int                `json:"date_end,omitempty"`
		Data      []map[string]int64 `json:"data,omitempty"`
	} `json:"result"`
}

// https://dev.freebox.fr/sdk/os/connection/
type connectionXdsl struct {
	apiResponse
	Result struct {
		Status struct {
			Status     string `json:"status"`
			Modulation string `json:"modulation"`
			Protocol   string `json:"protocol"`
			Uptime     int    `json:"uptime"`
		} `json:"status"`
		Down struct {
			Attn       int    `json:"attn"`
			Attn10     int    `json:"attn_10"`
			Crc        int    `json:"crc"`
			Es         int    `json:"es"`
			Fec        int    `json:"fec"`
			Ginp       bool   `json:"ginp"`
			Hec        int    `json:"hec"`
			Maxrate    uint64 `json:"maxrate"`
			Nitro      bool   `json:"nitro"`
			Phyr       bool   `json:"phyr"`
			Rate       int    `json:"rate"`
			RtxC       int    `json:"rtx_c,omitempty"`
			RtxTx      int    `json:"rtx_tx,omitempty"`
			RtxUc      int    `json:"rtx_uc,omitempty"`
			Rxmt       int    `json:"rxmt"`
			RxmtCorr   int    `json:"rxmt_corr"`
			RxmtUncorr int    `json:"rxmt_uncorr"`
			Ses        int    `json:"ses"`
			Snr        int    `json:"snr"`
			Snr10      int    `json:"snr_10"`
		} `json:"down"`
		Up struct {
			Attn       int    `json:"attn"`
			Attn10     int    `json:"attn_10"`
			Crc        int    `json:"crc"`
			Es         int    `json:"es"`
			Fec        int    `json:"fec"`
			Ginp       bool   `json:"ginp"`
			Hec        int    `json:"hec"`
			Maxrate    uint64 `json:"maxrate"`
			Nitro      bool   `json:"nitro"`
			Phyr       bool   `json:"phyr"`
			Rate       uint64 `json:"rate"`
			RtxC       int    `json:"rtx_c,omitempty"`
			RtxTx      int    `json:"rtx_tx,omitempty"`
			RtxUc      int    `json:"rtx_uc,omitempty"`
			Rxmt       int    `json:"rxmt"`
			RxmtCorr   int    `json:"rxmt_corr"`
			RxmtUncorr int    `json:"rxmt_uncorr"`
			Ses        int    `json:"ses"`
			Snr        int    `json:"snr"`
			Snr10      int    `json:"snr_10"`
		} `json:"up"`
	} `json:"result"`
}

type connectionFtth struct {
	apiResponse
	Result struct {
		SfpPresent         bool   `json:"sfp_present,omitempty"`
		SfpAlimOk		   bool   `json:"sfp_alim_ok,omitempty"`
		SfpHasPowerReport  bool   `json:"sfp_has_power_report,omitempty"`
		SfpHasSignal 	   bool   `json:"sfp_has_signal,omitempty"`
		Link               bool   `json:"link,omitempty"`
		SfpModel 		   string `json:"sfp_model,omitempty"`
		SfpSerial 		   string `json:"sfp_serial,omitempty"`
		SfpVendor          string `json:"sfp_vendor,omitempty"`
		SfpPwrTx		   int64  `json:"sfp_pwr_tx,omitempty"`
		SfpPwrRx		   int64  `json:"sfp_pwr_rx,omitempty"`
	} `json:"result"`
}

type database struct {
	DB        string   `json:"db"`
	DateStart int      `json:"date_start,omitempty"`
	DateEnd   int      `json:"date_end,omitempty"`
	Precision int      `json:"precision,omitempty"`
	Fields    []string `json:"fields"`
}

// https://dev.freebox.fr/sdk/os/freeplug/
type freeplug struct {
	apiResponse
	Result []freeplugNetwork `json:"result"`
}

type freeplugNetwork struct {
	ID      string           `json:"id"`
	Members []freeplugMember `json:"members"`
}

type freeplugMember struct {
	ID            string `json:"id"`
	Local         bool   `json:"local"`
	NetRole       string `json:"net_role"`
	EthPortStatus string `json:"eth_port_status"`
	EthFullDuplex bool   `json:"eth_full_duplex"`
	HasNetwork    bool   `json:"has_network"`
	EthSpeed      int    `json:"eth_speed"`
	Inative       int    `json:"inactive"`
	NetID         string `json:"net_id"`
	RxRate        int64  `json:"rx_rate"`
	TxRate        int64  `json:"tx_rate"`
	Model         string `json:"model"`
}

// https://dev.freebox.fr/sdk/os/lan/
type l3c struct {
	Addr string `json:"addr,omitempty"`
}

type lanHost struct {
	Reachable   bool   `json:"reachable,omitempty"`
	PrimaryName string `json:"primary_name,omitempty"`
	Vendor_name string `json:"vendor_name,omitempty"`
	L3c         []l3c  `json:"l3connectivities,omitempty"`
	L2Ident     struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"l2ident,omitempty"`
}

type lan struct {
	apiResponse
	Result []lanHost `json:"result"`
}

type idNameValue struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value int    `json:"value,omitempty"`
}

// https://dev.freebox.fr/sdk/os/system/
type system struct {
	apiResponse
	Result struct {
		Mac              string `json:"mac,omitempty"`
		FanRPM           int    `json:"fan_rpm,omitempty"`
		BoxFlavor        string `json:"box_flavor,omitempty"`
		TempCpub         int    `json:"temp_cpub,omitempty"`
		TempCpum         int    `json:"temp_cpum,omitempty"`
		DiskStatus       string `json:"disk_status,omitempty"`
		TempHDD          int    `json:"temp_hdd,omitempty"`
		BoardName        string `json:"board_name,omitempty"`
		TempSW           int    `json:"temp_sw,omitempty"`
		Uptime           string `json:"uptime,omitempty"`
		UptimeVal        int    `json:"uptime_val,omitempty"`
		UserMainStorage  string `json:"user_main_storage,omitempty"`
		BoxAuthenticated bool   `json:"box_authenticated,omitempty"`
		Serial           string `json:"serial,omitempty"`
		FirmwareVersion  string `json:"firmware_version,omitempty"`
	} `json:"result"`
}

type systemV6 struct {
	apiResponse
	Result struct {
		Mac       string        `json:"mac"`
		Sensors   []idNameValue `json:"sensors"`
		ModelInfo struct {
			NetOperator        string   `json:"net_operator"`
			SupportedLanguages []string `json:"supported_languages"`
			HasDsl             bool     `json:"has_dsl"`
			HasDect            bool     `json:"has_dect"`
			CustomerHddSlots   int      `json:"customer_hdd_slots"`
			WifiType           string   `json:"wifi_type"`
			HasHomeAutomation  bool     `json:"has_home_automation"`
			PrettyName         string   `json:"pretty_name"`
			Name               string   `json:"name"`
			HasLanSfp          bool     `json:"has_lan_sfp"`
			InternalHddSize    int      `json:"internal_hdd_size"`
			DefaultLanguage    string   `json:"default_language"`
			HasVm              bool     `json:"has_vm"`
			HasExpansions      bool     `json:"has_expansions"`
		} `json:"model_info"`
		Fans       []idNameValue `json:"fans"`
		Expansions []struct {
			Type      string `json:"type"`
			Present   bool   `json:"present"`
			Slot      int    `json:"slot"`
			ProbeDone bool   `json:"probe_done"`
			Supported bool   `json:"supported"`
			Bundle    string `json:"bundle"`
		} `json:"expansions"`
		BoardName        string `json:"board_name"`
		DiskStatus       string `json:"disk_status"`
		Uptime           string `json:"uptime"`
		UptimeVal        int    `json:"uptime_val"`
		UserMainStorage  string `json:"user_main_storage"`
		BoxAuthenticated bool   `json:"box_authenticated"`
		Serial           string `json:"serial"`
		FirmwareVersion  string `json:"firmware_version"`
	} `json:"result"`
}

// https://dev.freebox.fr/sdk/os/wifi/
type wifiAccessPoint struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

type wifi struct {
	apiResponse
	Result []wifiAccessPoint `json:"result,omitempty"`
}

type wifiStation struct {
	Hostname           string `json:"hostname,omitempty"`
	MAC                string `json:"mac,omitempty"`
	State              string `json:"state,omitempty"`
	Inactive           int    `json:"inactive,omitempty"`
	RXBytes            int64  `json:"rx_bytes,omitempty"`
	TXBytes            int64  `json:"tx_bytes,omitempty"`
	ConnectionDuration int    `json:"conn_duration,omitempty"`
	TXRate             int64  `json:"tx_rate,omitempty"`
	RXRate             int64  `json:"rx_rate,omitempty"`
	Signal             int    `json:"signal,omitempty"`
}

type wifiStations struct {
	apiResponse
	Result []wifiStation `json:"result,omitempty"`
}

type app struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

type api struct {
	authz        string
	login        string
	loginSession string
}

type store struct {
	location string
}

type authInfo struct {
	myApp    app
	myAPI    api
	myStore  store
	myReader *bufio.Reader
}

type postRequest struct {
	method, url, header string
}

// https://dev.freebox.fr/sdk/os/vpn/
type vpnServer struct {
	apiResponse
	Result  []struct {
		Success 	  bool   `json:"success"`
		RxBytes       int64  `json:"rx_bytes,omitempty"`
		Authenticated bool   `json:"authenticated,omitempty"`
		TxBytes       int64  `json:"tx_bytes,omitempty"`
		User          string `json:"user,omitempty"`
		ID            string `json:"id,omitempty"`
		Vpn           string `json:"vpn,omitempty"`
		SrcIP         string `json:"src_ip,omitempty"`
		AuthTime      int32  `json:"auth_time,omitempty"`
		LocalIP       string `json:"local_ip,omitempty"`
	} `json:"result,omitempty"`
}

type switchStatus struct {
	apiResponse
	Result []struct {
		ID      int    `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Duplex  string `json:"duplex,omitempty"`
		Link    string `json:"link,omitempty"`
		Mode    string `json:"mode,omitempty"`
		Speed   string `json:"speed,omitempty"`
		RrdId   string `json:"rrd_id,omitempty"`
		MacList []struct {
			Mac      string `json:"mac,omitempty"`
			Hostname string `json:"hostname"`
		} `json:"mac_list,omitempty"`
	} `json:"result,omitempty"`
}

type switchPort struct {
	apiResponse
	Result struct {
		RxBadBytes         int `json:"rx_bad_bytes,omitempty"`
		RxBroadcastPackets int `json:"rx_broadcast_packets,omitempty"`
		RxBytesRate        int `json:"rx_bytes_rate,omitempty"`
		RxErrPackets       int `json:"rx_err_packets,omitempty"`
		RxFcsPackets       int `json:"rx_fcs_packets,omitempty"`
		RxFragmentsPackets int `json:"rx_fragments_packets,omitempty"`
		RxGoodBytes        int `json:"rx_good_bytes,omitempty"`
		RxGoodPackets      int `json:"rx_good_packets,omitempty"`
		RxJabberPackets    int `json:"rx_jabber_packets,omitempty"`
		RxMulticastPackets int `json:"rx_multicast_packets,omitempty"`
		RxOversizePackets  int `json:"rx_oversize_packets,omitempty"`
		RxPacketsRate      int `json:"rx_packets_rate,omitempty"`
		RxPause            int `json:"rx_pause,omitempty"`
		RxUndersizePackets int `json:"rx_undersize_packets,omitempty"`
		RxUnicastPackets   int `json:"rx_unicast_packets,omitempty"`
		TxBroadcastPackets int `json:"tx_broadcast_packets,omitempty"`
		TxBytes            int `json:"tx_bytes,omitempty"`
		TxBytesRate        int `json:"tx_bytes_rate,omitempty"`
		TxCollisions       int `json:"tx_collisions,omitempty"`
		TxDeferred         int `json:"tx_deferred,omitempty"`
		TxExcessive        int `json:"tx_excessive,omitempty"`
		TxFcs              int `json:"tx_fcs,omitempty"`
		TxLate             int `json:"tx_late,omitempty"`
		TxMulticastPackets int `json:"tx_multicast_packets,omitempty"`
		TxMultiple         int `json:"tx_multiple,omitempty"`
		TxPackets          int `json:"tx_packets,omitempty"`
		TxPacketsRate      int `json:"tx_packets_rate,omitempty"`
		TxPause            int `json:"tx_pause,omitempty"`
		TxSingle           int `json:"tx_single,omitempty"`
		TxUnicastPackets   int `json:"tx_unicast_packets,omitempty"`
	} `json:"result,omitempty"`
}
