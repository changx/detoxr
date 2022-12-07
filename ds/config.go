package ds

import (
	"sync"
)

type config struct {
	HoneypotIP string `json:"honeypot"`
	LocalNS    string `json:"fast_ns"`
	DohService string `json:"doh_service_url"`
	NSAddr     string `json:"ns_addr"`
	WebAddr    string `json:"web_addr"`
	locker     *sync.RWMutex
}

var serverCfg *config

func initServerConfig() {
	serverCfg = &config{
		locker: &sync.RWMutex{},
	}
}

func SetHoneypotIP(ip string) {
	serverCfg.locker.Lock()
	lg.Infof("Set HONEYPOT IP to %s", ip)
	serverCfg.HoneypotIP = ip
	serverCfg.locker.Unlock()
}

func GetHoneypotIP() string {
	serverCfg.locker.RLock()
	defer serverCfg.locker.RUnlock()
	return serverCfg.HoneypotIP
}

func SetDohService(u string) {
	serverCfg.locker.Lock()
	lg.Infof("Set DoH service to %s", u)
	serverCfg.DohService = u
	serverCfg.locker.Unlock()
}

func GetDohService() string {
	serverCfg.locker.RLock()
	defer serverCfg.locker.RUnlock()
	return serverCfg.DohService
}

func SetLocalNS(ip string) {
	serverCfg.locker.Lock()
	lg.Infof("Set local NS server to %s", ip)
	serverCfg.LocalNS = ip
	serverCfg.locker.Unlock()
}

func GetLocalNS() string {
	serverCfg.locker.RLock()
	defer serverCfg.locker.RUnlock()
	return serverCfg.LocalNS
}

func GetNSAddr() string {
	return serverCfg.NSAddr
}

func GetWebAddr() string {
	return serverCfg.WebAddr
}

func SaveToConfigFile(f string) {}

func LoadFromConfigFile(f string) {}
