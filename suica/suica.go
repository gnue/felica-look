package suica

import (
	"../felica"
	"fmt"
	"time"
)

/*
 #include "suica_get.h"
*/
import "C"

type felica_module string

var Module felica_module = "Suica"

// Suica
type Suica struct {
	Hist []*SuicaValue // 利用履歴

	cardinfo felica.CardInfo // カード情報（生データ）
}

// Suica利用履歴データ
type SuicaValue struct {
	Date       time.Time // 処理日時
	Type       int       // 端末種
	Proc       int       // 処理
	InStation  int       // 入場駅（線区コード、駅順コード）
	OutStation int       // 出場駅（線区コード、駅順コード）
	Balance    int       // 残額
	No         int       // 連番
	Region     int       // リージョン

	Payment int // 利用料金（積増の場合はマイナス）
}

const (
	FELICA_POLLING_SUICA  = uint16(C.FELICA_POLLING_SUICA)  // Suicaシステムコード
	FELICA_SC_SUICA_VALUE = uint16(C.FELICA_SC_SUICA_VALUE) // Suica利用履歴データ・サービスコード
)

// *** felica_module メソッド
// 対応カードか？
func (module *felica_module) IsCard(cardinfo felica.CardInfo) bool {
	for syscode, _ := range cardinfo {
		if syscode == FELICA_POLLING_SUICA {
			return true
		}
	}

	return false
}

// CardInfo を束縛した Engine を作成する
func (module *felica_module) Bind(cardinfo felica.CardInfo) felica.Engine {
	return &Suica{cardinfo: cardinfo}
}

// *** Suica メソッド
// カード名
func (suica *Suica) Name() string {
	return string(Module)
}

// カード情報を読込む
func (suica *Suica) Read() {
	if 0 < len(suica.Hist) {
		// 読込済みなら何もしない
		return
	}

	cardinfo := suica.cardinfo

	// システムデータの取得
	currsys := cardinfo[FELICA_POLLING_SUICA]

	// Suica利用履歴
	for _, raw := range currsys.Services[FELICA_SC_SUICA_VALUE] {
		history := (*C.suica_value_t)(felica.DataPtr(&raw))
		h_time := C.suica_value_date(history)
		if h_time == 0 {
			continue
		}

		value := SuicaValue{}
		value.Date = time.Unix(int64(h_time), 0)
		value.Type = int(C.suica_value_type(history))
		value.Proc = int(C.suica_value_proc(history))
		value.InStation = int(C.suica_value_in_station(history))
		value.OutStation = int(C.suica_value_out_station(history))
		value.Balance = int(C.suica_value_balance(history))
		value.No = int(C.suica_value_no(history))
		value.Region = int(C.suica_value_region(history))

		suica.Hist = append(suica.Hist, &value)
	}

	// 利用料金の計算
	for i, value := range suica.Hist[:len(suica.Hist)-1] {
		pre_data := suica.Hist[i+1]
		value.Payment = pre_data.Balance - value.Balance
	}
}

// カード情報を表示する
func (suica *Suica) ShowInfo(options *felica.Options) {
	// テーブルデータの読込み
	if suica_tables == nil {
		suica_tables, _ = felica.LoadYAML("suica.yml")
	}

	// データの読込み
	suica.Read()

	// 表示
	if options.Extend {
		fmt.Println()
		fmt.Println("[利用履歴（元データ）]")
		fmt.Println("   利用年月日        残額    入場駅  出場駅 / リージョン (連番） 端末種  処理")
		fmt.Println("  -------------------------------------------------------------------------------------")
		for _, value := range suica.Hist {
			fmt.Printf("   %s  %8d円    0x%04X  0x%04X /    0x%02X    (%4d)  0x%04X  0x%02X\n",
				value.Date.Format("2006-01-02"),
				value.Balance,
				value.InStation,
				value.OutStation,
				value.Region,
				value.No,
				value.Type,
				value.Proc)
		}
	}

	fmt.Println()
	fmt.Println("[利用履歴]")
	fmt.Println("   利用年月日    支払い       残額    入場駅  出場駅 / リージョン (連番） 端末種  処理")
	fmt.Println("  ----------------------------------------------------------------------------------------------------------------------")
	for _, value := range suica.Hist {
		fmt.Printf("   %s  %6d円 %8d円    0x%04X  0x%04X /    0x%02X    (%4d)  0x%04X  0x%02X\n",
			value.Date.Format("2006-01-02"),
			value.Payment,
			value.Balance,
			value.InStation,
			value.OutStation,
			value.Region,
			value.No,
			value.Type,
			value.Proc)
	}
}

// ***
// Suicaテーブル
var suica_tables map[interface{}]interface{}

// Suicaテーブルを検索して表示用の文字列を返す
func suica_disp_name(name string, value int, base int, opt_values ...int) interface{} {
	return felica.DispName(suica_tables, name, value, base, opt_values...)
}
