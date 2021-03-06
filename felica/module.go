package felica

type Module interface {
	IsCard(cardinfo CardInfo) bool // 対応カードか？
	Bind(cardinfo CardInfo) Engine // CardInfo を束縛した Engine を作成する
}

// 登録モジュール
var Modules = make(map[string]Module)

// モジュールの登録
func Register(name string, module Module) {
	Modules[name] = module
}

// カードに対応するモジュールを探す
func Find(cardinfo CardInfo) Module {
	for _, m := range Modules {
		if m.IsCard(cardinfo) {
			return m
		}
	}

	return nil
}
