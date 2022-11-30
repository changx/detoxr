package ds

func SetHoneypotIP(ip string) {
	resolverConfig.honeypotIP = ip
}

func SetDohService(u string) {
	resolverConfig.dohService = u
}

func SetLocalNS(ip string) {
	resolverConfig.localNS = ip
}
