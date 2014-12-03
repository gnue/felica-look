package nicepass

import (
	"fmt"
	"strings"
	"time"

	"github.com/gnue/felica-look/felica"
)

/*
 #include "nicepass_get.h"
*/
import "C"

var moduleName = "nice-pass"

func init() {
	felica.Register(moduleName, new(felica_module))
}

type felica_module struct {
}

// nice-pass/遠州鉄道
type Nicepass struct {
	Hist []*NicepassValue // 利用履歴

	cardinfo felica.CardInfo // カード情報（生データ）
	syscode  uint16          // システムコード
}

// nice-pass利用履歴データ
type NicepassValue struct {
	DateTime   time.Time // 処理日時
	Train      int       // 装置番号
	InStation  int       // 乗車駅
	OutStation int       // 降車駅
	Type       int       // 端末種
	Proc       int       // 処理
	Use        int       // 利用金額（支払いはマイナス）
	Balance    int       // 残額

	Raw []byte // Rawデータ
}

const (
	FELICA_POLLING_NICEPASS = uint16(C.FELICA_POLLING_NICEPASS) // nice-passシステムコード

	FELICA_SC_NICEPASS_VALUE = uint16(C.FELICA_SC_NICEPASS_VALUE) // nice-pass利用履歴データ・サービスコード
)

// システムコード・リスト
var SystemCodes = map[uint16]string{
	FELICA_POLLING_NICEPASS: "nice-pass",
}

// *** felica_module メソッド

// 対応カードか？
func (module *felica_module) IsCard(cardinfo felica.CardInfo) bool {
	syscode := find_syscode(cardinfo)
	return (syscode != 0)
}

// CardInfo を束縛した Engine を作成する
func (module *felica_module) Bind(cardinfo felica.CardInfo) felica.Engine {
	return &Nicepass{cardinfo: cardinfo, syscode: find_syscode(cardinfo)}
}

// *** Nicepass メソッド

// カード名
func (nicepass *Nicepass) Name() string {
	return SystemCodes[nicepass.syscode]
}

// カード情報を読込む
func (nicepass *Nicepass) Read() {
	if 0 < len(nicepass.Hist) {
		// 読込済みなら何もしない
		return
	}

	cardinfo := nicepass.cardinfo

	// システムデータの取得
	currsys := cardinfo[nicepass.syscode]

	// nice-pass利用履歴
	for _, raw := range currsys.Services[FELICA_SC_NICEPASS_VALUE] {
		history := (*C.nicepass_value_t)(felica.DataPtr(&raw))
		h_time := C.nicepass_value_datetime(history)
		if h_time == 0 {
			continue
		}

		value := NicepassValue{}
		value.DateTime = time.Unix(int64(h_time), 0)
		value.Train = int(C.nicepass_value_train(history))
		value.InStation = int(C.nicepass_value_in_station(history))
		value.OutStation = int(C.nicepass_value_out_station(history))
		value.Type = int(C.nicepass_value_type(history))
		value.Proc = int(C.nicepass_value_proc(history))
		value.Use = int(C.nicepass_value_use(history))
		value.Balance = int(C.nicepass_value_balance(history))
		value.Raw = raw

		nicepass.Hist = append(nicepass.Hist, &value)
	}
}

// カード情報を表示する
func (nicepass *Nicepass) ShowInfo(options *felica.Options) {
	// テーブルデータの読込み
	if nicepass_tables == nil {
		nicepass_tables, _ = felica.LoadYAML("nicepass.yml")
	}

	// データの読込み
	nicepass.Read()

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
		fmt.Printf("%s          日時         利用金額      残額     乗車駅    降車駅    装置番号     端末種    処理\n", indent_space)
		fmt.Printf("  %s\n", strings.Repeat("-", indent+106))
		for _, value := range nicepass.Hist {
			if options.Hex {
				fmt.Printf("   %16X   ", value.Raw)
			}
			fmt.Printf("   %s %8d円 %8d円    0x%05X   0x%05X     0x%04X    %v  %v\n",
				value.DateTime.Format("2006-01-02 15:04"),
				value.Use,
				value.Balance,
				value.InStation,
				value.OutStation,
				value.Train,
				t(value.TypeName(), 10),
				value.ProcName())
		}
	}

	fmt.Println("\n[利用履歴]\n")
	fmt.Printf("%s          日時          支払い       残額       乗車駅           降車駅            端末種    処理\n", indent_space)
	fmt.Printf("  %s\n", strings.Repeat("-", indent+110))
	for _, value := range nicepass.Hist {
		disp_payment := "---　"

		if 0 < value.Use {
			// チャージ
			disp_payment = fmt.Sprintf("(+%d円)", value.Use)
		} else if value.Use < 0 {
			disp_payment = fmt.Sprintf("%d円", -value.Use)
		}

		if options.Hex {
			fmt.Printf("   %16X   ", value.Raw)
		}
		fmt.Printf("   %s  %8s %8d円    %v  %v  %v  %v\n",
			value.DateTime.Format("2006-01-02 15:04"),
			disp_payment,
			value.Balance,
			t(value.InStationName(), 16),
			t(value.OutStationName(), 16),
			t(value.TypeName(), 10),
			value.ProcName())
	}
}

// *** NicepassValue メソッド

// 端末種
func (value *NicepassValue) TypeName() interface{} {
	return nicepass_disp_name("TYPE", value.Type, 4)
}

// 処理
func (value *NicepassValue) ProcName() interface{} {
	return nicepass_disp_name("PROC", value.Proc, 2)
}

// 乗車駅
func (value *NicepassValue) InStationName() interface{} {
	if value.InStation == 0 {
		return ""
	}
	return nicepass_disp_name("STATION", value.InStation, 6)
}

// 降車駅
func (value *NicepassValue) OutStationName() interface{} {
	if value.OutStation == 0 {
		return ""
	}
	return nicepass_disp_name("STATION", value.OutStation, 6)
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

// 表示幅を指定した文字数
func t(value interface{}, width int) string {
	return felica.DispString(fmt.Sprintf("%v", value), width)
}

// ***

// nice-passテーブル
var nicepass_tables map[interface{}]interface{}

// nice-passテーブルを検索して表示用の文字列を返す
func nicepass_disp_name(name string, value int, base int, opt_values ...int) interface{} {
	return felica.DispName(nicepass_tables, name, value, base, opt_values...)
}
