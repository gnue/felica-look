package rapica

import (
	"fmt"
	"strings"
	"time"

	"github.com/gnue/felica-look/felica"
)

/*
 #include "rapica_get.h"
*/
import "C"

var moduleName = "RapiCa"

func init() {
	felica.Register(moduleName, new(felica_module))
}

type felica_module struct {
}

// RapiCa/鹿児島市交通局
type RapiCa struct {
	Info    RapicaInfo      // 発行情報
	Attr    RapicaAttr      // 属性情報
	Hist    []*RapicaValue  // 利用履歴
	Charges []*RapicaCharge // 積増情報

	cardinfo felica.CardInfo // カード情報（生データ）
}

// RapiCa発行情報データ
type RapicaInfo struct {
	Date    time.Time // 発行日
	Company int       // 事業者
	Deposit int       // デポジット金額

	Raw [][]byte // Rawデータ
}

// RapiCa属性情報データ
type RapicaAttr struct {
	DateTime   time.Time // 直近処理日時
	Company    int       // 事業者
	TicketNo   int       // 整理券番号
	Busstop    int       // 停留所
	Busline    int       // 系統
	Busno      int       // 装置
	Kind       int       // 利用種別
	Amount     int       // 残額
	Premier    int       // プレミア
	Point      int       // ポイント
	No         int       // 取引連番
	OnBusstop  int       // 乗車停留所(整理券)番号
	OffBusstop int       // 降車停留所(整理券)番号
	Payment    int       // 利用金額
	Point2     int       // ポイント？

	Raw [][]byte // Rawデータ
}

// Rapica利用履歴データ
type RapicaValue struct {
	DateTime time.Time // 処理日時
	Company  int       // 事業者
	Busstop  int       // 停留所
	Busline  int       // 系統
	Busno    int       // 装置
	Kind     int       // 利用種別
	Amount   int       // 残額

	Payment  int // 利用料金（積増の場合はマイナス）
	OnValue  int // 対応する乗車データ
	OffValue int // 対応する降車データ

	Raw []byte // Rawデータ
}

// Rapica積増情報データ
type RapicaCharge struct {
	Date    time.Time // 積増日付
	Charge  int       // 積増金額
	Premier int       // プレミア
	Company int       // 事業者

	Raw []byte // Rawデータ
}

const (
	FELICA_POLLING_RAPICA   = uint16(C.FELICA_POLLING_RAPICA)   // RapiCa/鹿児島市交通局
	FELICA_SC_RAPICA_INFO   = uint16(C.FELICA_SC_RAPICA_INFO)   // RapiCa発行情報・サービスコード
	FELICA_SC_RAPICA_ATTR   = uint16(C.FELICA_SC_RAPICA_ATTR)   // RapiCa属性情報・サービスコード
	FELICA_SC_RAPICA_VALUE  = uint16(C.FELICA_SC_RAPICA_VALUE)  // RapiCa利用履歴データ・サービスコード
	FELICA_SC_RAPICA_CHARGE = uint16(C.FELICA_SC_RAPICA_CHARGE) // RapiCa積増情報データ・サービスコード
)

// *** felica_module メソッド

// 対応カードか？
func (module *felica_module) IsCard(cardinfo felica.CardInfo) bool {
	for syscode, _ := range cardinfo {
		if syscode == FELICA_POLLING_RAPICA {
			return true
		}
	}

	return false
}

// CardInfo を束縛した Engine を作成する
func (module *felica_module) Bind(cardinfo felica.CardInfo) felica.Engine {
	return &RapiCa{cardinfo: cardinfo}
}

// *** RapiCa メソッド

// カード名
func (rapica *RapiCa) Name() string {
	return moduleName
}

// カード情報を読込む
func (rapica *RapiCa) Read() {
	if rapica.Info.Company != 0 {
		// 読込済みなら何もしない
		return
	}

	cardinfo := rapica.cardinfo

	// システムデータの取得
	currsys := cardinfo[FELICA_POLLING_RAPICA]

	// RapiCa発行情報
	raw_info := currsys.Services[FELICA_SC_RAPICA_INFO]
	rapica.Info.Raw = raw_info

	info := (*C.rapica_info_t)(felica.DataPtr(&raw_info[0]))
	i_time := C.rapica_info_date(info)

	rapica.Info.Company = int(C.rapica_info_company(info))
	rapica.Info.Deposit = int(C.rapica_info_deposit(info))
	rapica.Info.Date = time.Unix(int64(i_time), 0)

	// RapiCa属性情報
	raw_attr := currsys.Services[FELICA_SC_RAPICA_ATTR]
	rapica.Attr.Raw = raw_attr

	// RapiCa属性情報(1)
	attr1 := (*C.rapica_attr1_t)(felica.DataPtr(&raw_attr[0]))
	a_time := C.rapica_attr_time(attr1)

	rapica.Attr.DateTime = time.Unix(int64(a_time), 0)
	rapica.Attr.Company = int(C.rapica_attr_company(attr1))
	rapica.Attr.TicketNo = int(attr1.ticketno)
	rapica.Attr.Busstop = int(C.rapica_attr_busstop(attr1))
	rapica.Attr.Busline = int(C.rapica_attr_busline(attr1))
	rapica.Attr.Busno = int(C.rapica_attr_busno(attr1))

	// RapiCa属性情報(2)
	attr2 := (*C.rapica_attr2_t)(felica.DataPtr(&raw_attr[1]))
	rapica.Attr.Kind = int(C.rapica_attr_kind(attr2))
	rapica.Attr.Amount = int(C.rapica_attr_amount(attr2))
	rapica.Attr.Premier = int(C.rapica_attr_premier(attr2))
	rapica.Attr.Point = int(C.rapica_attr_point(attr2))
	rapica.Attr.No = int(C.rapica_attr_no(attr2))
	rapica.Attr.OnBusstop = int(attr2.on_busstop)
	rapica.Attr.OffBusstop = int(attr2.off_busstop)

	// RapiCa属性情報(3)
	attr3 := (*C.rapica_attr3_t)(felica.DataPtr(&raw_attr[2]))
	rapica.Attr.Payment = int(C.rapica_attr_payment(attr3))

	// RapiCa属性情報(4)
	attr4 := (*C.rapica_attr4_t)(felica.DataPtr(&raw_attr[3]))
	rapica.Attr.Point2 = int(C.rapica_attr_point2(attr4))

	// RapiCa利用履歴
	last_time := C.time_t(rapica.Attr.DateTime.Unix())

	for _, raw := range currsys.Services[FELICA_SC_RAPICA_VALUE] {
		history := (*C.rapica_value_t)(felica.DataPtr(&raw))
		h_time := C.rapica_value_datetime(history, last_time)
		if h_time == 0 {
			continue
		}

		value := RapicaValue{}
		value.DateTime = time.Unix(int64(h_time), 0)
		value.Company = int(history.company)
		value.Busstop = int(C.rapica_value_busstop(history))
		value.Busline = int(C.rapica_value_busline(history))
		value.Busno = int(C.rapica_value_busno(history))
		value.Kind = int(history.kind)
		value.Amount = int(C.rapica_value_amount(history))
		value.OnValue = -1
		value.OffValue = -1
		value.Raw = raw

		rapica.Hist = append(rapica.Hist, &value)
		last_time = h_time
	}

	// 乗車データと降車データの関連付けをする
	for i, value := range rapica.Hist[:len(rapica.Hist)-1] {
		pre_data := rapica.Hist[i+1]

		if value.Kind == C.RAPICA_KIND_GETOFF {
			// 降車
			for j, v := range rapica.Hist[i+1:] {
				if v.Kind == C.RAPICA_KIND_GETON {
					// 乗車を見つけた
					value.OnValue = i + 1 + j
					v.OffValue = i
					break
				}
			}
		}

		value.Payment = pre_data.Amount - value.Amount
	}

	// RapiCa積増情報
	for _, raw := range currsys.Services[FELICA_SC_RAPICA_CHARGE] {
		charge := (*C.rapica_charge_t)(felica.DataPtr(&raw))
		c_time := C.rapica_charge_date(charge)
		if c_time == 0 {
			continue
		}

		c := RapicaCharge{}
		c.Date = time.Unix(int64(c_time), 0)
		c.Charge = int(C.rapica_charge_charge(charge))
		c.Premier = int(C.rapica_charge_premier(charge))
		c.Company = int(C.rapica_charge_company(charge))
		c.Raw = raw

		rapica.Charges = append(rapica.Charges, &c)
	}
}

// カード情報を表示する
func (rapica *RapiCa) ShowInfo(options *felica.Options) {
	// テーブルデータの読込み
	if rapica_tables == nil {
		rapica_tables, _ = felica.LoadYAML("rapica.yml")
	}

	// データの読込み
	rapica.Read()

	// インデント
	indent := 0
	indent_space := ""

	if options.Hex {
		indent = 38
		indent_space = strings.Repeat(" ", indent)
	}

	// 表示
	attr := rapica.Attr

	fmt.Println("\n[発行情報]")

	if options.Hex {
		fmt.Println()
		for _, v := range rapica.Info.Raw {
			fmt.Printf("   %16X\n", v)
		}
	}

	fmt.Printf(
		`
  事業者: %v
  発行日: %s
  デポジット金額: %d円
`, rapica.Info.CompanyName(), rapica.Info.Date.Format("2006-01-02"), rapica.Info.Deposit)

	fmt.Println("\n[属性情報]")

	if options.Hex {
		fmt.Println()
		for _, v := range rapica.Attr.Raw {
			fmt.Printf("   %16X\n", v)
		}
	}

	fmt.Printf(`
  直近処理日時:	%s
  事業者:	%v
  整理券番号:	%d
  停留所:	0x%06X
  系統:		0x%04X
  装置・車号？:	%d
  利用種別:	%v
  残額:		%d円
  プレミア:	%d円
  ポイント:	%dpt
  取引連番:	%d
  乗車停留所(整理券)番号: %d
  降車停留所(整理券)番号: %d
  利用金額:	%d円
  ポイント？:	%dpt
`, attr.DateTime.Format("2006-01-02 15:04"),
		attr.CompanyName(),
		attr.TicketNo, attr.Busstop, attr.Busline, attr.Busno,
		attr.KindName(),
		attr.Amount, attr.Premier, attr.Point, attr.No, attr.OnBusstop, attr.OffBusstop,
		attr.Payment, attr.Point2)

	if options.Extend || options.Hex {
		fmt.Println("\n[利用履歴（元データ）]\n")
		fmt.Printf("%s       日時     利用種別      残額             事業者                     系統              /      停留所          (装置)\n", indent_space)
		fmt.Printf("  %s\n", strings.Repeat("-", indent+120))
		for _, value := range rapica.Hist {
			if options.Hex {
				fmt.Printf("   %16X   ", value.Raw)
			}
			fmt.Printf("   %s    %v  %8d円  %v %v / %v (%d)\n",
				value.DateTime.Format("01/02 15:04"),
				value.KindName(),
				value.Amount,
				t(value.CompanyName(), 24),
				t(value.BuslineName(), 30),
				t(value.BusstopName(), 20),
				value.Busno)
		}
	}

	fmt.Println("\n[利用履歴]\n")
	fmt.Printf("%s          日時       利用種別      利用料金        残額             事業者                     系統              /              停留所                (装置)\n", indent_space)
	fmt.Printf("  %s\n", strings.Repeat("-", indent+156))
	for _, value := range rapica.Hist {
		disp_payment := "---　"
		disp_busstop := value.BusstopName()

		if 0 <= value.OffValue && value.Payment == 0 {
			// 対応する降車データがあり利用金額が 0 ならば表示しない
			continue
		}

		if value.Payment < 0 {
			// 積増
			disp_payment = fmt.Sprintf("(+%d円)", -value.Payment)
		} else if 0 < value.Payment {
			disp_payment = fmt.Sprintf("%d円", value.Payment)
		}

		if 0 <= value.OnValue {
			OnValue := rapica.Hist[value.OnValue]
			disp_busstop = fmt.Sprintf("%v -> %v", OnValue.BusstopName(), disp_busstop)
		}

		if options.Hex {
			fmt.Printf("   %16X   ", value.Raw)
		}
		fmt.Printf("   %s    %v %14s %9d円  %v %v / %v (%d)\n",
			value.DateTime.Format("2006-01-02 15:04"),
			value.KindName(),
			disp_payment,
			value.Amount,
			t(value.CompanyName(), 24),
			t(value.BuslineName(), 30),
			t(disp_busstop, 34),
			value.Busno)
	}

	fmt.Println("\n[積増情報]\n")
	fmt.Printf("%s      日時       チャージ   プレミア    事業者\n", indent_space)
	fmt.Printf("  %s\n", strings.Repeat("-", indent+48))
	for _, charge := range rapica.Charges {
		if options.Hex {
			fmt.Printf("   %16X   ", charge.Raw)
		}
		fmt.Printf("   %s %8d円 %8d円    %v\n",
			charge.Date.Format("2006-01-02"), charge.Charge, charge.Premier, charge.CompanyName())
	}
}

// *** RapicaInfo メソッド

// 事業者名
func (info *RapicaInfo) CompanyName() interface{} {
	return rapica_disp_name("ATTR_COMPANY", info.Company, 4)
}

// *** RapicaAttr メソッド

// 事業者名
func (attr *RapicaAttr) CompanyName() interface{} {
	return rapica_disp_name("ATTR_COMPANY", attr.Company, 4)
}

// 利用種別
func (attr *RapicaAttr) KindName() interface{} {
	return rapica_disp_name("ATTR_KIND", attr.Kind&0xff0000, 6, attr.Kind)
}

// *** RapicaValue メソッド

// 利用種別
func (value *RapicaValue) KindName() interface{} {
	return rapica_disp_name("HIST_KIND", value.Kind, 2)
}

// 事業者名
func (value *RapicaValue) CompanyName() interface{} {
	return rapica_disp_name("HIST_COMPANY", value.Company>>4, 2, value.Company)
}

// 停留所
func (value *RapicaValue) BusstopName() interface{} {
	return rapica_disp_name("BUSSTOP", value.Busstop, 6)
}

// 系統名
func (value *RapicaValue) BuslineName() interface{} {
	return rapica_disp_name("BUSLINE", (value.Busstop&0xff0000)+value.Busline, 4, value.Busline)
}

// *** RapicaCharge メソッド

func (charge *RapicaCharge) CompanyName() interface{} {
	return rapica_disp_name("ATTR_COMPANY", charge.Company, 4)
}

// ***

// RapiCaテーブル
var rapica_tables map[interface{}]interface{}

// RapiCaテーブルを検索して表示用の文字列を返す
func rapica_disp_name(name string, value int, base int, opt_values ...int) interface{} {
	return felica.DispName(rapica_tables, name, value, base, opt_values...)
}

// 表示幅を指定した文字数
func t(value interface{}, width int) string {
	return felica.DispString(fmt.Sprintf("%v", value), width)
}
