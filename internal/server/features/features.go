package features

import (
	"errors"
	"log"
	"net"
	"net/http"
)

var privateCIDRs []*net.IPNet

func init() {
	for _, cidr := range []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateCIDRs = append(privateCIDRs, block)
	}
}

func fromPrivateIP(flag string, r *http.Request) (bool, error) {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false, err
	}
	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return false, errors.New("can't parse ip")
	}
	if ip.IsLoopback() {
		return true, nil
	}
	return ip.IsPrivate(), nil
}

var enabledFunctions map[string]Enabled

func init() {
	enabledFunctions = map[string]Enabled{}
	enabledFunctions["new-storage"] = fromPrivateIP
}

type Enabled func(flag string, r *http.Request) (bool, error)

func FeatureEnabled(flag string, r *http.Request) bool {
	if viper.IsSet(flag) {
		return viper.GetBool(flag)
	}

	enabledFunc, exists := enabledFunctions[flag]
	if !exists {
		return false
	}
	res, err := enabledFunc(flag, r)
	if err != nil {
		log.Println(err)
		return false
	}
	return res
}
