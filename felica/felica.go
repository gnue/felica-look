package felica

import (
	"reflect"
	"unsafe"
)

// カード情報
type CardInfo map[uint16]*SystemInfo

// システム情報
type SystemInfo struct {
	IDm          string
	PMm          string
	ServiceCodes []uint16
	Services     ServiceInfo
}

// オプションフラグ
type Options struct {
	Extend bool // 拡張表示
	Hex    bool // データの16進表示もいっしょに表示する
}

// サービス情報
type ServiceInfo map[uint16]([][]byte)

type Module interface {
	Name() string                                 // カード名
	SystemCode() uint16                           // システムコード
	ShowInfo(cardinfo CardInfo, options *Options) // カード情報を表示する
}

// *** CardInfo のメソッド
// C言語で使うためにデータにアクセスするポインタを取得する
func (sysinfo *SystemInfo) SvcDataPtr(svccode uint16, index int) unsafe.Pointer {
	data := sysinfo.Services[svccode]
	raw := (*reflect.SliceHeader)(unsafe.Pointer(&data[index])).Data

	return unsafe.Pointer(raw)
}
