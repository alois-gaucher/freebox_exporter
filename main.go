package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mafreebox string
	listen    string
	debug     bool
	fiber     bool
	v6        bool
)

func init() {
	flag.StringVar(&mafreebox, "endpoint", "http://mafreebox.freebox.fr/", "Endpoint for freebox API")
	flag.StringVar(&listen, "listen", ":10001", "Prometheus metrics port")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&fiber, "fiber", false, "Turn on if you're using a fiber Freebox")
	flag.BoolVar(&v6, "v6", false, "Use v6+ system API endpoint")
}

func main() {
	flag.Parse()

	if !strings.HasSuffix(mafreebox, "/") {
		mafreebox = mafreebox + "/"
	}

	endpoint := mafreebox + "api/v4/login/"
	myAuthInfo := &authInfo{
		myAPI: api{
			login:        endpoint,
			authz:        endpoint + "authorize/",
			loginSession: endpoint + "session/",
		},
		myStore: store{location: os.Getenv("HOME") + "/.freebox_token"},
		myApp: app{
			AppID:      "fr.freebox.exporter",
			AppName:    "prometheus-exporter",
			AppVersion: "0.4",
			DeviceName: "local",
		},
		myReader: bufio.NewReader(os.Stdin),
	}

	myPostRequest := &postRequest{
		method: "POST",
		url:    mafreebox + "api/v4/rrd/",
		header: "X-Fbx-App-Auth",
	}

	myConnectionXdslRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/connection/xdsl/",
		header: "X-Fbx-App-Auth",
	}

	myConnectionFtthRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/connection/ftth/",
		header: "X-Fbx-App-Auth",
	}

	myFreeplugRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/freeplug/",
		header: "X-Fbx-App-Auth",
	}

	myLanRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/lan/browser/pub/",
		header: "X-Fbx-App-Auth",
	}

	mySystemRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/system/",
		header: "X-Fbx-App-Auth",
	}

	mySystemV6Request := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v6/system/",
		header: "X-Fbx-App-Auth",
	}

	myWifiRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v2/wifi/ap/",
		header: "X-Fbx-App-Auth",
	}

	myVpnRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/vpn/connection/",
		header: "X-Fbx-App-Auth",
	}

	mySwitchStatusRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v8/switch/status/",
		header: "X-Fbx-App-Auth",
	}

	var mySessionToken string

	go func() {
		for {
			// There is no DSL metric on fiber Freebox
			// If you use a fiber Freebox, use -fiber flag to turn off this metric
			if !fiber {
				// connectionXdsl metrics
				connectionXdslStats, err := getConnectionXdsl(myAuthInfo, myConnectionXdslRequest, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with connectionXdsl metrics: %v", err)
				}

				if connectionXdslStats.Success {
					status := connectionXdslStats.Result.Status
					result := connectionXdslStats.Result
					down := result.Down
					up := result.Up

					connectionXdslStatusUptimeGauges.
						WithLabelValues(status.Status, status.Protocol, status.Modulation).
						Set(float64(status.Uptime))

					connectionXdslDownAttnGauge.Set(float64(down.Attn10) / 10)
					connectionXdslUpAttnGauge.Set(float64(up.Attn10) / 10)

					// XXX: sometimes the Freebox is reporting zero as SNR which
					// does not make sense so we don't log these
					if down.Snr10 > 0 {
						connectionXdslDownSnrGauge.Set(float64(down.Snr10) / 10)
					}
					if up.Snr10 > 0 {
						connectionXdslUpSnrGauge.Set(float64(up.Snr10) / 10)
					}

					connectionXdslNitroGauges.WithLabelValues("down").
						Set(bool2float(down.Nitro))
					connectionXdslNitroGauges.WithLabelValues("up").
						Set(bool2float(up.Nitro))

					connectionXdslGinpGauges.WithLabelValues("down", "enabled").
						Set(bool2float(down.Ginp))
					connectionXdslGinpGauges.WithLabelValues("up", "enabled").
						Set(bool2float(up.Ginp))

					logFields(result, connectionXdslGinpGauges,
						[]string{"rtx_tx", "rtx_c", "rtx_uc"})

					logFields(result, connectionXdslErrorGauges,
						[]string{"crc", "es", "fec", "hec", "ses"})
				}

				// dsl metrics
				getDslResult, err := getDsl(myAuthInfo, myPostRequest, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with DSL metrics: %v", err)
				}

				if len(getDslResult) > 0 {
					rateUpGauge.Set(float64(getDslResult[0]))
					rateDownGauge.Set(float64(getDslResult[1]))
					snrUpGauge.Set(float64(getDslResult[2]))
					snrDownGauge.Set(float64(getDslResult[3]))
				}
			} else {
				// connectionFtth metrics
				connectionFtthStats, err := getConnectionFtth(myAuthInfo, myConnectionFtthRequest, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with connectionFtth metrics: %v", err)
				}

				if connectionFtthStats.Success {
					SfpHasPowerReport := connectionFtthStats.Result.SfpHasPowerReport
					SfpHasSignal := connectionFtthStats.Result.SfpHasSignal
					SfpModel := connectionFtthStats.Result.SfpModel
					SfpVendor := connectionFtthStats.Result.SfpVendor
					SfpPwrRx := connectionFtthStats.Result.SfpPwrRx
					SfpPwrTx := connectionFtthStats.Result.SfpPwrTx
					Link := connectionFtthStats.Result.Link
					SfpAlimOk := connectionFtthStats.Result.SfpAlimOk
					SfpSerial := connectionFtthStats.Result.SfpSerial
					SfpPresent := connectionFtthStats.Result.SfpPresent

					connectionFtthStatusUptimeGauges.With(prometheus.Labels{"sfp_model": SfpModel, "sfp_vendor": SfpVendor, "sfp_serial": SfpSerial, "sfp_present": strconv.FormatBool(SfpPresent), "sfp_has_power_report": strconv.FormatBool(SfpHasPowerReport), "sfp_has_signal": strconv.FormatBool(SfpHasSignal), "sfp_alim_ok": strconv.FormatBool(SfpAlimOk), "link": strconv.FormatBool(Link)}).Set(float64(1))
					connectionFtthRxPwrGauge.Set(float64(SfpPwrRx) / 100)
					connectionFtthTxPwrGauge.Set(float64(SfpPwrTx) / 100)
				}

				// getFtthResult, err := getFtth(myAuthInfo, myPostRequest, &mySessionToken)
				// if err != nil {
				// 	log.Printf("An error occured with FTTH metrics: %v", err)
				// }
			}

			// freeplug metrics
			freeplugStats, err := getFreeplug(myAuthInfo, myFreeplugRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with freeplug metrics: %v", err)
			}

			for _, freeplugNetwork := range freeplugStats.Result {
				for _, freeplugMember := range freeplugNetwork.Members {
					if freeplugMember.HasNetwork {
						freeplugHasNetworkGauge.WithLabelValues(freeplugMember.ID).Set(float64(1))
					} else {
						freeplugHasNetworkGauge.WithLabelValues(freeplugMember.ID).Set(float64(0))
					}

					Mb := 1e6
					rxRate := float64(freeplugMember.RxRate) * Mb
					txRate := float64(freeplugMember.TxRate) * Mb

					if rxRate >= 0 { // -1 if not unavailable
						freeplugRxRateGauge.WithLabelValues(freeplugMember.ID).Set(rxRate)
					}

					if txRate >= 0 { // -1 if not unavailable
						freeplugTxRateGauge.WithLabelValues(freeplugMember.ID).Set(txRate)
					}
				}
			}

			// net metrics
			getNetResult, err := getNet(myAuthInfo, myPostRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with NET metrics: %v", err)
			}

			if len(getNetResult) > 0 {
				bwUpGauge.Set(float64(getNetResult[0]))
				bwDownGauge.Set(float64(getNetResult[1]))
				netRateUpGauge.Set(float64(getNetResult[2]))
				netRateDownGauge.Set(float64(getNetResult[3]))
				vpnRateUpGauge.Set(float64(getNetResult[4]))
				vpnRateDownGauge.Set(float64(getNetResult[5]))
			}

			// lan metrics
			lanAvailable, err := getLan(myAuthInfo, myLanRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with LAN metrics: %v", err)
			}
			for _, v := range lanAvailable {
				var Ip string
				if len(v.L3c) > 0 {
					Ip = v.L3c[0].Addr
				} else {
					Ip = ""
				}
				if v.Reachable {
					lanReachableGauges.With(prometheus.Labels{"name": v.PrimaryName, "vendor": v.Vendor_name, "mac": v.L2Ident.ID, "ip": Ip}).Set(float64(1))
				} else {
					lanReachableGauges.With(prometheus.Labels{"name": v.PrimaryName, "vendor": v.Vendor_name, "mac": v.L2Ident.ID, "ip": Ip}).Set(float64(0))
				}
			}

			// system metrics
			if v6 {
				systemStats, err := getSystemV6(myAuthInfo, mySystemV6Request, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with System metrics: %v", err)
				}

				for _, sensor := range systemStats.Result.Sensors {
					systemTempGauges.WithLabelValues(sensor.Name).Set(float64(sensor.Value))
				}
				for _, fan := range systemStats.Result.Fans {
					systemFanGauges.WithLabelValues(fan.Name).Set(float64(fan.Value))
				}

				systemUptimeGauges.
					WithLabelValues(systemStats.Result.FirmwareVersion).
					Set(float64(systemStats.Result.UptimeVal))
			} else {
				systemStats, err := getSystem(myAuthInfo, mySystemRequest, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with System metrics: %v", err)
				}

				systemTempGauges.WithLabelValues("Température CPU B").Set(float64(systemStats.Result.TempCpub))
				systemTempGauges.WithLabelValues("Température CPU M").Set(float64(systemStats.Result.TempCpum))
				systemTempGauges.WithLabelValues("Température Switch").Set(float64(systemStats.Result.TempSW))
				systemTempGauges.WithLabelValues("Disque dur").Set(float64(systemStats.Result.TempHDD))
				systemFanGauges.WithLabelValues("Ventilateur 1").Set(float64(systemStats.Result.FanRPM))

				systemUptimeGauges.
					WithLabelValues(systemStats.Result.FirmwareVersion).
					Set(float64(systemStats.Result.UptimeVal))
			}

			// wifi metrics
			wifiStats, err := getWifi(myAuthInfo, myWifiRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with Wifi metrics: %v", err)
			}
			for _, accessPoint := range wifiStats.Result {
				myWifiStationRequest := &postRequest{
					method: "GET",
					url:    mafreebox + "api/v2/wifi/ap/" + strconv.Itoa(accessPoint.ID) + "/stations",
					header: "X-Fbx-App-Auth",
				}
				wifiStationsStats, err := getWifiStations(myAuthInfo, myWifiStationRequest, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with Wifi station metrics: %v", err)
				}
				for _, station := range wifiStationsStats.Result {
					wifiSignalGauges.With(prometheus.Labels{"access_point": accessPoint.Name, "mac": station.MAC, "hostname": station.Hostname, "state": station.State}).Set(float64(station.Signal))
					wifiInactiveGauges.With(prometheus.Labels{"access_point": accessPoint.Name, "mac": station.MAC, "hostname": station.Hostname, "state": station.State}).Set(float64(station.Inactive))
					wifiConnectionDurationGauges.With(prometheus.Labels{"access_point": accessPoint.Name, "mac": station.MAC, "hostname": station.Hostname, "state": station.State}).Set(float64(station.ConnectionDuration))
					wifiRXBytesGauges.With(prometheus.Labels{"access_point": accessPoint.Name, "mac": station.MAC, "hostname": station.Hostname, "state": station.State}).Set(float64(station.RXBytes))
					wifiTXBytesGauges.With(prometheus.Labels{"access_point": accessPoint.Name, "mac": station.MAC, "hostname": station.Hostname, "state": station.State}).Set(float64(station.TXBytes))
					wifiRXRateGauges.With(prometheus.Labels{"access_point": accessPoint.Name, "mac": station.MAC, "hostname": station.Hostname, "state": station.State}).Set(float64(station.RXRate))
					wifiTXRateGauges.With(prometheus.Labels{"access_point": accessPoint.Name, "mac": station.MAC, "hostname": station.Hostname, "state": station.State}).Set(float64(station.TXRate))
				}
			}

			// VPN Server Connections List
			getVpnServerResult, err := getVpnServer(myAuthInfo, myVpnRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with VPN station metrics: %v", err)
			}
			for _, connection := range getVpnServerResult.Result {
				vpnServerConnectionsList.With(prometheus.Labels{"user": connection.User, "vpn": connection.Vpn, "src_ip": connection.SrcIP, "local_ip": connection.LocalIP, "name": "rx_bytes"}).Set(float64(connection.RxBytes))
				vpnServerConnectionsList.With(prometheus.Labels{"user": connection.User, "vpn": connection.Vpn, "src_ip": connection.SrcIP, "local_ip": connection.LocalIP, "name": "tx_bytes"}).Set(float64(connection.TxBytes))
			}

			// Switch status
			switchStats, err := getSwitchStatus(myAuthInfo, mySwitchStatusRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with switch metrics: %v", err)
			}
			for _, port := range switchStats.Result {
				if port.Link == "up" {
					mySwitchPortRequest := &postRequest{
						method: "GET",
						url:    mafreebox + "api/v8/switch/port/" + strconv.Itoa(port.ID) + "/stats",
						header: "X-Fbx-App-Auth",
					}
					switchPortStats, err := getSwitchPort(myAuthInfo, mySwitchPortRequest, &mySessionToken)
					if err != nil {
						log.Printf("An error occured with switch port metrics: %v", err)
					}

					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "broadcast", "error": "0"}).Set(float64(switchPortStats.Result.RxBroadcastPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "multicast", "error": "0"}).Set(float64(switchPortStats.Result.RxMulticastPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "unicast", "error": "0"}).Set(float64(switchPortStats.Result.RxUnicastPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "broadcast", "error": "0"}).Set(float64(switchPortStats.Result.TxBroadcastPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "multicast", "error": "0"}).Set(float64(switchPortStats.Result.TxMulticastPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "unicast", "error": "0"}).Set(float64(switchPortStats.Result.TxUnicastPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "err", "error": "1"}).Set(float64(switchPortStats.Result.RxErrPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "fcs", "error": "1"}).Set(float64(switchPortStats.Result.RxFcsPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "fragment", "error": "1"}).Set(float64(switchPortStats.Result.RxFragmentsPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "jabber", "error": "1"}).Set(float64(switchPortStats.Result.RxJabberPackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "oversize", "error": "1"}).Set(float64(switchPortStats.Result.RxOversizePackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "undersize", "error": "1"}).Set(float64(switchPortStats.Result.RxUndersizePackets))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "collision", "error": "1"}).Set(float64(switchPortStats.Result.TxCollisions))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "deferred", "error": "1"}).Set(float64(switchPortStats.Result.TxDeferred))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "excessive", "error": "1"}).Set(float64(switchPortStats.Result.TxExcessive))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "fcs", "error": "1"}).Set(float64(switchPortStats.Result.TxFcs))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "late", "error": "1"}).Set(float64(switchPortStats.Result.TxLate))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "multiple", "error": "1"}).Set(float64(switchPortStats.Result.TxMultiple))
					switchPortPacketsGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "single", "error": "1"}).Set(float64(switchPortStats.Result.TxSingle))

					switchPortPacketsTotalGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx"}).Set(float64(switchPortStats.Result.RxGoodPackets))
					switchPortPacketsTotalGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx"}).Set(float64(switchPortStats.Result.TxPackets))

					switchPortBytesGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "bad"}).Set(float64(switchPortStats.Result.RxBadBytes))
					switchPortBytesGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx", "type": "good"}).Set(float64(switchPortStats.Result.RxGoodBytes))
					switchPortBytesGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx", "type": "total"}).Set(float64(switchPortStats.Result.TxBytes))

					switchPortPauseGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx"}).Set(float64(switchPortStats.Result.RxPause))
					switchPortPauseGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx"}).Set(float64(switchPortStats.Result.TxPause))

					switchPortPacketsRateGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx"}).Set(float64(switchPortStats.Result.RxPacketsRate))
					switchPortPacketsRateGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx"}).Set(float64(switchPortStats.Result.TxPacketsRate))

					switchPortBytesRateGauges.With(prometheus.Labels{"name": port.Name, "direction": "rx"}).Set(float64(switchPortStats.Result.RxBytesRate))
					switchPortBytesRateGauges.With(prometheus.Labels{"name": port.Name, "direction": "tx"}).Set(float64(switchPortStats.Result.TxBytesRate))
				}
			}

			time.Sleep(10 * time.Second)
		}
	}()

	log.Println("freebox_exporter started on port", listen)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}

func logFields(result interface{}, gauge *prometheus.GaugeVec, fields []string) error {
	resultReflect := reflect.ValueOf(result)

	for _, direction := range []string{"down", "up"} {
		for _, field := range fields {
			value := reflect.Indirect(resultReflect).
				FieldByName(strcase.ToCamel(direction)).
				FieldByName(strcase.ToCamel(field))

			if value.IsZero() {
				continue
			}

			gauge.WithLabelValues(direction, field).
				Set(float64(value.Int()))
		}
	}

	return nil
}

func bool2float(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
