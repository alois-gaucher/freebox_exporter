package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// XXX: see https://dev.freebox.fr/sdk/os/ for API documentation
	// XXX: see https://prometheus.io/docs/practices/naming/ for metric names

	// connectionXdsl
	connectionXdslStatusUptimeGauges = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_status_uptime_seconds_total",
	},
		[]string{
			"status",
			"protocol",
			"modulation",
		},
	)

	connectionXdslDownAttnGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_down_attn_decibels",
	})
	connectionXdslUpAttnGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_up_attn_decibels",
	})
	connectionXdslDownSnrGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_down_snr_decibels",
	})
	connectionXdslUpSnrGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_up_snr_decibels",
	})

	connectionXdslErrorGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_connection_xdsl_errors_total",
			Help: "Error counts",
		},
		[]string{
			"direction", // up|down
			"name",      // crc|es|fec|hec
		},
	)

	connectionXdslGinpGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_connection_xdsl_ginp",
		},
		[]string{
			"direction", // up|down
			"name",      // enabled|rtx_(tx|c|uc)
		},
	)

	connectionXdslNitroGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_connection_xdsl_nitro",
		},
		[]string{
			"direction", // up|down
		},
	)

	// connectionFtth
	connectionFtthStatusUptimeGauges = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_connection_ftth_status_uptime_seconds_total",
	},
		[]string{
			"sfp_model",
			"sfp_vendor",
			"sfp_serial",
			"sfp_has_power_report",
			"sfp_has_signal",
			"link",
			"sfp_alim_ok",
			"sfp_present",
		},
	)

	connectionFtthRxPwrGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_ftth_sfp_rx_pwr_decibels",
	})

	connectionFtthTxPwrGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_ftth_sfp_tx_pwr_decibels",
	})

	// RRD dsl [unstable]
	rateUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_up_bytes",
		Help: "Available upload bandwidth (in byte/s)",
	})
	rateDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_down_bytes",
		Help: "Available download bandwidth (in byte/s)",
	})
	snrUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_up_decibel",
		Help: "Upload signal/noise ratio (in 1/10 dB)",
	})
	snrDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_down_decibel",
		Help: "Download signal/noise ratio (in 1/10 dB)",
	})

	// freeplug
	freeplugRxRateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_freeplug_rx_rate_bits",
		Help: "rx rate (from the freeplugs to the \"cco\" freeplug) (in bits/s) -1 if not available",
	},
		[]string{
			"id",
		},
	)
	freeplugTxRateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_freeplug_tx_rate_bits",
		Help: "tx rate (from the \"cco\" freeplug to the freeplugs) (in bits/s) -1 if not available",
	},
		[]string{
			"id",
		},
	)
	freeplugHasNetworkGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_freeplug_has_network",
		Help: "is connected to the network",
	},
		[]string{
			"id",
		},
	)

	// RRD Net [unstable]
	bwUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_up_bytes",
		Help: "Upload available bandwidth (in byte/s)",
	})
	bwDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_down_bytes",
		Help: "Download available bandwidth (in byte/s)",
	})
	netRateUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_up_bytes",
		Help: "Upload rate (in byte/s)",
	})
	netRateDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_down_bytes",
		Help: "Download rate (in byte/s)",
	})
	vpnRateUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_up_bytes",
		Help: "Vpn client upload rate (in byte/s)",
	})
	vpnRateDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_down_bytes",
		Help: "Vpn client download rate (in byte/s)",
	})

	// Lan
	lanReachableGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_lan_reachable",
			Help: "Hosts reachable on LAN",
		},
		[]string{
			"name", // hostname
			"vendor",
			"mac",
			"ip",
		},
	)

	systemTempGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_temp_celsius",
			Help: "Temperature sensors reported by system (in °C)",
		},
		[]string{
			"name",
		},
	)

	systemFanGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_fan_rpm",
			Help: "Fan speed reported by system (in RPM)",
		},
		[]string{
			"name",
		},
	)

	systemUptimeGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_uptime_seconds_total",
		},
		[]string{
			"firmware_version",
		},
	)

	// wifi
	wifiLabels = []string{
		"access_point",
		"mac",
		"hostname",
		"state",
	}

	wifiSignalGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_wifi_signal_attenuation_db",
			Help: "Wifi signal attenuation in decibel",
		},
		wifiLabels,
	)

	wifiInactiveGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_wifi_inactive_duration_seconds",
			Help: "Wifi inactive duration in seconds",
		},
		wifiLabels,
	)

	wifiConnectionDurationGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_wifi_connection_duration_seconds",
			Help: "Wifi connection duration in seconds",
		},
		wifiLabels,
	)

	wifiRXBytesGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_wifi_rx_bytes",
			Help: "Wifi received data (from station to Freebox) in bytes",
		},
		wifiLabels,
	)

	wifiTXBytesGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_wifi_tx_bytes",
			Help: "Wifi transmitted data (from Freebox to station) in bytes",
		},
		wifiLabels,
	)

	wifiRXRateGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_wifi_rx_rate",
			Help: "Wifi reception data rate (from station to Freebox) in bytes/seconds",
		},
		wifiLabels,
	)

	wifiTXRateGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_wifi_tx_rate",
			Help: "Wifi transmission data rate (from Freebox to station) in bytes/seconds",
		},
		wifiLabels,
	)

	// vpn server connections list [unstable]
	vpnServerConnectionsList = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpn_server_connections_list",
			Help: "VPN server connections list",
		},
		[]string{
			"user",
			"vpn",
			"src_ip",
			"local_ip",
			"name", // rx_bytes|tx_bytes
		},
	)

	switchPortPacketsGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_switch_port_packets",
		},
		[]string{
			"name",
			"direction",
			"type",
			"error",
		},
	)

	switchPortPacketsTotalGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_switch_port_packets_total",
		},
		[]string{
			"name",
			"direction",
		},
	)

	switchPortBytesGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_switch_port_bytes",
		},
		[]string{
			"name",
			"direction",
			"type",
		},
	)

	switchPortBytesRateGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_switch_port_bytes_rate",
		},
		[]string{
			"name",
			"direction",
		},
	)

	switchPortPacketsRateGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_switch_port_packets_rate",
		},
		[]string{
			"name",
			"direction",
		},
	)

	switchPortPauseGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_switch_port_pause",
		},
		[]string{
			"name",
			"direction",
		},
	)
)
