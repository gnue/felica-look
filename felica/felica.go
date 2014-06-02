package felica

import (
	"sort"
)

// カード情報
type CardInfo map[uint16]*SystemInfo

// システム情報
type SystemInfo struct {
	IDm      string
	PMm      string
	Services ServiceInfo
}

// サービス情報
type ServiceInfo map[uint16]([][]byte)

// オプションフラグ
type Options struct {
	Extend bool // 拡張表示
	Hex    bool // データの16進表示もいっしょに表示する
}

type Module interface {
	IsCard(cardinfo CardInfo) bool // 対応カードか？
	Bind(cardinfo CardInfo) Engine // CardInfo を束縛した Engine を作成する
}

type Engine interface {
	Name() string              // カード名
	ShowInfo(options *Options) // カード情報を表示する
}

// 出力フォーマット
const (
	OUTPUT_NORMAL = iota
	OUTPUT_JSON
	OUTPUT_LTSV
)

// *** ソート用
type ByUint16 []uint16

func (a ByUint16) Len() int           { return len(a) }
func (a ByUint16) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUint16) Less(i, j int) bool { return a[i] < a[j] }

// *** SystemInfoメソッド
func (sysinfo *SystemInfo) ServiceCodes() []uint16 {
	codes := make([]uint16, 0, len(sysinfo.Services))

	for svccode, _ := range sysinfo.Services {
		codes = append(codes, svccode)
	}
	sort.Sort(ByUint16(codes))

	return codes
}
