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
	IDm          string
	PMm          string
	ServiceCodes []string
	Services     ServiceInfo
}

// サービス情報
type ServiceInfo map[string]([][]byte)

type Module interface {
	Name() string                            // カード名
	SystemCode() uint64                      // システムコード
	ShowInfo(cardinfo CardInfo, extend bool) // カード情報を表示する
}

// *** CardInfo のメソッド
// システムコードから SystemInfo を取得する
func (cardinfo CardInfo) sysinfo(syscode uint64) *SystemInfo {
	return cardinfo[fmt.Sprintf("%04X", syscode)]
}

// サービスコードからデータを取得する
func (sysinfo SystemInfo) svcdata(svccode uint64) [][]byte {
	return sysinfo.Services[fmt.Sprintf("%04X", svccode)]
}

// C言語で使うためにデータにアクセスするポインタを取得する
func (sysinfo *SystemInfo) svcdata_ptr(svccode uint64, index int) unsafe.Pointer {
	data := sysinfo.svcdata(svccode)
	raw := (*reflect.SliceHeader)(unsafe.Pointer(&data[index])).Data

	return unsafe.Pointer(raw)
}
