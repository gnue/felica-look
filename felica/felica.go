package felica

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
	IsCard(cardinfo CardInfo) bool // 対応カードか？
	Bind(cardinfo CardInfo) Engine // CardInfo を束縛した Engine を作成する
}

type Engine interface {
	Name() string              // カード名
	ShowInfo(options *Options) // カード情報を表示する
}
