package felica

// カード情報
type CardInfo map[string]*SystemInfo

// システム情報
type SystemInfo struct {
	idm      string
	pmm      string
	svccodes []string
	services ServiceInfo
}

// サービス情報
type ServiceInfo map[string]([][]byte)

// SystemInfoメンバーの getter
func (sysinfo SystemInfo) IDm() string {
	return sysinfo.idm
}

func (sysinfo SystemInfo) PMm() string {
	return sysinfo.pmm
}

func (sysinfo SystemInfo) Services() ServiceInfo {
	return sysinfo.services
}

func (sysinfo SystemInfo) ServiceCodes() []string {
	return sysinfo.svccodes
}
