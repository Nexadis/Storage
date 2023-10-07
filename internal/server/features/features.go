package features

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/spf13/viper"
)

var privateCIDRs []*net.IPNet

type Feature string

const (
	Use_user Feature = "use_user"
)

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

func fromPrivateIP(flag Feature, r *http.Request) (bool, error) {
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

var enabledFunctions map[Feature]Enabled

func init() {
	enabledFunctions = map[Feature]Enabled{}
	enabledFunctions[Use_user] = fromPrivateIP
}

type Enabled func(flag Feature, r *http.Request) (bool, error)

func FeatureEnabled(flag Feature, r *http.Request) bool {
	if viper.IsSet(string(flag)) {
		return viper.GetBool(string(flag))
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
