package ds

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
)

type config struct {
	HoneypotAddr string        `json:"honeypot"`
	LocalNS      []string      `json:"fast_ns"`
	DohService   string        `json:"doh_service_url"`
	NSAddr       string        `json:"ns_listen"`
	WebAddr      string        `json:"web_listen"`
	locker       *sync.RWMutex `json:"-"`
}

var serverCfg *config

func InitServerConfig() {
	serverCfg = &config{
		HoneypotAddr: "202.12.27.33:53",
		LocalNS: []string{
			"114.114.114.114:53",
			"223.5.5.5:53",
		},
		DohService: "https://dns.google/resolve",
		NSAddr:     ":1053",
		WebAddr:    ":3000",
		locker:     &sync.RWMutex{},
	}
	cwd, _ := os.Getwd()
	configJsonFilePath := path.Join(cwd, "config.json")
	fmt.Printf("config.json: %s\n", configJsonFilePath)
	if _, err := os.Stat(configJsonFilePath); err == nil {
		cfgData, err := os.ReadFile(configJsonFilePath)
		if err == nil {
			json.Unmarshal(cfgData, serverCfg)
		}
	}
	fmt.Printf("server config: %+v\n", serverCfg)
}

func WriteConfigJson() {
	defer recover()
	cwd, _ := os.Getwd()
	configJsonFilePath := path.Join(cwd, "config.json")
	writer, err := os.OpenFile(configJsonFilePath, os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err == nil {
		json.NewEncoder(writer).Encode(serverCfg)
	}
}

func SetHoneypotIP(ip string) {
	serverCfg.locker.Lock()
	lg.Infof("Set HONEYPOT to %s", ip)
	serverCfg.HoneypotAddr = ip
	serverCfg.locker.Unlock()
}

func GetHoneypotIP() string {
	serverCfg.locker.RLock()
	defer serverCfg.locker.RUnlock()
	return serverCfg.HoneypotAddr
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

func SetLocalNS(iplist []string) {
	serverCfg.locker.Lock()
	lg.Infof("Set local NS server to %v", iplist)
	serverCfg.LocalNS = iplist
	serverCfg.locker.Unlock()
}

func GetLocalNS() []string {
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
