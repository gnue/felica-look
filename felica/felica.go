package felica

import (
	"fmt"
	"reflect"
	"unsafe"
)

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

// *** CardInfo のメソッド
// システムコードから SystemInfo を取得する
func (cardinfo CardInfo) sysinfo(syscode uint64) *SystemInfo {
	return cardinfo[fmt.Sprintf("%04X", syscode)]
}

// *** SystemInfo のメソッド
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

// サービスコードからデータを取得する
func (sysinfo SystemInfo) svcdata(svccode uint64) [][]byte {
	return sysinfo.services[fmt.Sprintf("%04X", svccode)]
}

// C言語で使うためにデータにアクセスするポインタを取得する
func (sysinfo *SystemInfo) svcdata_ptr(svccode uint64, index int) unsafe.Pointer {
	data := sysinfo.svcdata(svccode)
	raw := (*reflect.SliceHeader)(unsafe.Pointer(&data[index])).Data

	return unsafe.Pointer(raw)
}
