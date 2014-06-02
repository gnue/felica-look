package suica

import (
	"fmt"
	"strings"
	"time"

	"github.com/gnue/felica-look/felica"
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
	syscode  uint16          // システムコード
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

	Raw []byte // Rawデータ
}

const (
	FELICA_POLLING_SUICA  = uint16(C.FELICA_POLLING_SUICA)  // Suicaシステムコード
	FELICA_POLLING_IRUCA  = uint16(C.FELICA_POLLING_IRUCA)  // IruCaシステムコード
	FELICA_POLLING_SAPICA = uint16(C.FELICA_POLLING_SAPICA) // SAPICAシステムコード
	FELICA_POLLING_PASPY  = uint16(C.FELICA_POLLING_PASPY)  // PASPYシステムコード

	FELICA_SC_SUICA_VALUE = uint16(C.FELICA_SC_SUICA_VALUE) // Suica利用履歴データ・サービスコード
)

// システムコード・リスト
var SystemCodes = map[uint16]string{
	FELICA_POLLING_SUICA:  "Suica",
	FELICA_POLLING_IRUCA:  "IruCa",
	FELICA_POLLING_SAPICA: "SAPICA",
	FELICA_POLLING_PASPY:  "PASPY",
}

// *** felica_module メソッド
// 対応カードか？
func (module *felica_module) IsCard(cardinfo felica.CardInfo) bool {
	syscode := find_syscode(cardinfo)
	return (syscode != 0)
}

// CardInfo を束縛した Engine を作成する
func (module *felica_module) Bind(cardinfo felica.CardInfo) felica.Engine {
	return &Suica{cardinfo: cardinfo, syscode: find_syscode(cardinfo)}
}

// *** Suica メソッド
// カード名
func (suica *Suica) Name() string {
	return SystemCodes[suica.syscode]
}

// カード情報を読込む
func (suica *Suica) Read() {
	if 0 < len(suica.Hist) {
		// 読込済みなら何もしない
		return
	}

	cardinfo := suica.cardinfo

	// システムデータの取得
	currsys := cardinfo[suica.syscode]

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
		value.Proc = int(history.proc)
		value.InStation = int(C.suica_value_in_station(history))
		value.OutStation = int(C.suica_value_out_station(history))
		value.Balance = int(C.suica_value_balance(history))
		value.No = int(C.suica_value_no(history))
		value.Region = int(history.region)
		value.Raw = raw

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

	// インデント
	indent := 0
	indent_space := ""

	if options.Hex {
		indent = 38
		indent_space = strings.Repeat(" ", indent)
	}

	// 表示
	if options.Extend || options.Hex {
		fmt.Println("\n[利用履歴（元データ）]\n")
		fmt.Printf("%s   利用年月日        残額    入場駅  出場駅 / リージョン (連番） 端末種    処理\n", indent_space)
		fmt.Printf("  %s\n", strings.Repeat("-", indent+106))
		for _, value := range suica.Hist {
			if options.Hex {
				fmt.Printf("   %16X   ", value.Raw)
			}
			fmt.Printf("   %s  %8d円    0x%04X  0x%04X /    0x%02X    (%4d)  %6v  %v\n",
				value.Date.Format("2006-01-02"),
				value.Balance,
				value.InStation,
				value.OutStation,
				value.Region,
				value.No,
				value.TypeName(),
				value.ProcName())
		}
	}

	fmt.Println("\n[利用履歴]\n")
	fmt.Printf("%s   利用年月日     支払い       残額     入場駅      出場駅   (連番） 端末種    処理\n", indent_space)
	fmt.Printf("  %s\n", strings.Repeat("-", indent+110))
	for _, value := range suica.Hist {
		disp_payment := "---　"

		if value.Payment < 0 {
			// チャージ
			disp_payment = fmt.Sprintf("(+%d円)", -value.Payment)
		} else if 0 < value.Payment {
			disp_payment = fmt.Sprintf("%d円", value.Payment)
		}

		if options.Hex {
			fmt.Printf("   %16X   ", value.Raw)
		}
		fmt.Printf("   %s  %8s %8d円  %10v  %10v  (%4d)  %6v  %v\n",
			value.Date.Format("2006-01-02"),
			disp_payment,
			value.Balance,
			value.InStationName(),
			value.OutStationName(),
			value.No,
			value.TypeName(),
			value.ProcName())
	}
}

// *** SuicaValue メソッド
// 処理
func (value *SuicaValue) ProcName() interface{} {
	return suica_disp_name("PROC", value.Proc, 2)
}

// 入場駅
func (value *SuicaValue) InStationName() interface{} {
	if value.InStation == 0 {
		return ""
	}
	return suica_disp_name("STATION", (value.Region<<16)+value.InStation, 6)
}

// 出場駅
func (value *SuicaValue) OutStationName() interface{} {
	if value.OutStation == 0 {
		return ""
	}
	return suica_disp_name("STATION", (value.Region<<16)+value.OutStation, 6)
}

// 端末種
func (value *SuicaValue) TypeName() interface{} {
	return suica_disp_name("TYPE", value.Type, 4)
}

// *** 関数
// システムコードを検索する
func find_syscode(cardinfo felica.CardInfo) uint16 {
	for syscode, _ := range cardinfo {
		if len(SystemCodes[syscode]) != 0 {
			return syscode
		}
	}

	return 0
}

// ***
// Suicaテーブル
var suica_tables map[interface{}]interface{}

// Suicaテーブルを検索して表示用の文字列を返す
func suica_disp_name(name string, value int, base int, opt_values ...int) interface{} {
	return felica.DispName(suica_tables, name, value, base, opt_values...)
}
