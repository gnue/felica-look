package edy

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/gnue/felica-look/felica"
)

/*
 #include "edy_get.h"
*/
import "C"

var moduleName = "Edy"

func init() {
	felica.Register(moduleName, new(felica_module))
}

type felica_module struct {
}

// Edy
type Edy struct {
	Info EdyInfo     // Edyカード情報
	Last EdyLast     // Edy残額情報
	Hist []*EdyValue // 利用履歴

	cardinfo felica.CardInfo // カード情報（生データ）
}

// Edyカード情報
type EdyInfo struct {
	EdyNo []byte // Edy番号

	Raw [][]byte // Rawデータ
}

// Edy残額情報（最終利用状況）
type EdyLast struct {
	Rest int // 残額
	Use  int // 直近使用金額（チャージのときは更新されない場合がある）
	No   int // 取引通番

	Raw [][]byte // Rawデータ
}

// Edy利用履歴データ
type EdyValue struct {
	DateTime time.Time // 処理日時
	Type     int       // 端末種
	No       int       // 連番
	Use      int       // 入金／出金
	Rest     int       // 残額金

	Charge  int // チャージ
	Payment int // 支払い

	Raw []byte // Rawデータ
}

const (
	FELICA_POLLING_EDY  = uint16(C.FELICA_POLLING_EDY)  // Edyシステムコード
	FELICA_SC_EDY_INFO  = uint16(C.FELICA_SC_EDY_INFO)  // Edyカード情報・サービスコード
	FELICA_SC_EDY_LAST  = uint16(C.FELICA_SC_EDY_LAST)  // Edy残額情報・サービスコード
	FELICA_SC_EDY_VALUE = uint16(C.FELICA_SC_EDY_VALUE) // Edy利用履歴データ・サービスコード
)

// *** felica_module メソッド

// 対応カードか？
func (module *felica_module) IsCard(cardinfo felica.CardInfo) bool {
	for syscode, currsys := range cardinfo {
		if syscode == FELICA_POLLING_EDY {
			if currsys.Services[FELICA_SC_EDY_INFO] != nil {
				return true
			}
			break
		}
	}

	return false
}

// CardInfo を束縛した Engine を作成する
func (module *felica_module) Bind(cardinfo felica.CardInfo) felica.Engine {
	return &Edy{cardinfo: cardinfo}
}

// *** Edy メソッド

// カード名
func (edy *Edy) Name() string {
	return moduleName
}

// カード情報を読込む
func (edy *Edy) Read() {
	if 0 < len(edy.Hist) {
		// 読込済みなら何もしない
		return
	}

	cardinfo := edy.cardinfo

	// システムデータの取得
	currsys := cardinfo[FELICA_POLLING_EDY]

	// Edyカード情報
	raw_info := currsys.Services[FELICA_SC_EDY_INFO]
	edy.Info.Raw = raw_info

	info := (*C.edy_info0_t)(felica.DataPtr(&raw_info[0]))

	edyno := unsafe.Pointer(&info.edyno[0])
	edy.Info.EdyNo = C.GoBytes(edyno, C.int(unsafe.Sizeof(edyno)))

	// Edy残額情報（最終利用状況）
	raw_last := currsys.Services[FELICA_SC_EDY_LAST]
	edy.Last.Raw = raw_last

	last := (*C.edy_last_t)(felica.DataPtr(&raw_last[0]))
	edy.Last.Rest = int(C.edy_last_rest(last))
	edy.Last.Use = int(C.edy_last_use(last))
	edy.Last.No = int(C.edy_last_no(last))

	// Edy利用履歴
	for _, raw := range currsys.Services[FELICA_SC_EDY_VALUE] {
		history := (*C.edy_value_t)(felica.DataPtr(&raw))
		h_time := C.edy_value_datetime(history)
		if h_time == 0 {
			continue
		}

		value := EdyValue{}
		value.DateTime = time.Unix(int64(h_time), 0)
		value.Type = int(C.edy_value_type(history))
		value.No = int(C.edy_value_no(history))
		value.Use = int(C.edy_value_use(history))
		value.Rest = int(C.edy_value_rest(history))
		value.Raw = raw

		edy.Hist = append(edy.Hist, &value)
	}

	for _, value := range edy.Hist {
		switch value.Type {
		case 0x02, 0x04: // 入金（チャージ）, 入金（Edyギフト）
			value.Charge = value.Use
		case 0x20: // 出金
			value.Payment = value.Use
		}
	}
}

// カード情報を表示する
func (edy *Edy) ShowInfo(options *felica.Options) {
	// テーブルデータの読込み
	if edy_tables == nil {
		edy_tables, _ = felica.LoadYAML("edy.yml")
	}

	// データの読込み
	edy.Read()

	// インデント
	indent := 0
	indent_space := ""

	if options.Hex {
		indent = 38
		indent_space = strings.Repeat(" ", indent)
	}

	// 表示
	fmt.Println("\n[Edyカード情報]")
	if options.Hex {
		fmt.Println()
		for _, v := range edy.Info.Raw {
			fmt.Printf("   %16X\n", v)
		}
	}
	fmt.Printf("\n  Edy番号: %v\n", edy.Info.EdyNoDisp())

	fmt.Println("\n[Edy残額情報]")
	if options.Hex {
		fmt.Println()
		for _, v := range edy.Last.Raw {
			fmt.Printf("   %16X\n", v)
		}
	}
	fmt.Printf(`
  残額: %14d円
  直近使用金額: %6d円
  取引通番: %10d
`, edy.Last.Rest, edy.Last.Use, edy.Last.No)

	if options.Extend || options.Hex {
		fmt.Println("\n[利用履歴（元データ）]\n")
		fmt.Printf("%s      利用年月日         支払い        残額  (連番)  タイプ\n", indent_space)
		fmt.Printf("  %s\n", strings.Repeat("-", indent+85))
		for _, value := range edy.Hist {
			if options.Hex {
				fmt.Printf("   %16X   ", value.Raw)
			}
			fmt.Printf("   %s  %8d円  %8d円  (%4d)  %v\n",
				value.DateTime.Format("2006-01-02 15:04"),
				value.Use,
				value.Rest,
				value.No,
				value.TypeName())
		}
	}

	fmt.Println("\n[利用履歴]\n")
	fmt.Printf("%s      利用年月日        チャージ      支払い        残額  (連番)  タイプ\n", indent_space)
	fmt.Printf("  %s\n", strings.Repeat("-", indent+98))
	for _, value := range edy.Hist {
		if options.Hex {
			fmt.Printf("   %16X   ", value.Raw)
		}
		fmt.Printf("   %s  %10v  %10v  %8d円  (%4d)  %v\n",
			value.DateTime.Format("2006-01-02 15:04"),
			disp_money(value.Charge),
			disp_money(value.Payment),
			value.Rest,
			value.No,
			value.TypeName())
	}
}

// *** EdyInfo メソッド

// Edy番号
func (info *EdyInfo) EdyNoDisp() string {
	edyno := info.EdyNo
	return fmt.Sprintf("%0X-%0X-%0X-%X", edyno[:2], edyno[2:4], edyno[4:6], edyno[6:])
}

// タイプ
func (value *EdyValue) TypeName() interface{} {
	return edy_disp_name("TYPE", value.Type, 2)
}

// *** 表示用関数

// 金額（0円なら空文字列）
func disp_money(money int) string {
	if money == 0 {
		return ""
	}

	return fmt.Sprintf("%d円", money)
}

// ***

// Edyテーブル
var edy_tables map[interface{}]interface{}

// Edyテーブルを検索して表示用の文字列を返す
func edy_disp_name(name string, value int, base int, opt_values ...int) interface{} {
	return felica.DispName(edy_tables, name, value, base, opt_values...)
}
