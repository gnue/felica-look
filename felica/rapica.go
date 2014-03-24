package felica

import (
	"fmt"
	"time"
)

/*
 #include "rapica/rapica.c"
*/
import "C"

// RapiCa/鹿児島市交通局
type RapiCa struct {
	Info    RapicaInfo      // 発行情報
	Attr    RapicaAttr      // 属性情報
	Hist    []*RapicaValue  // 利用履歴
	Charges []*RapicaCharge // 積増情報
}

// RapiCa発行情報データ
type RapicaInfo struct {
	Date    time.Time // 発行日
	Company int       // 事業者
	Deposit int       // デポジット金額
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
}

// Rapica積増情報データ
type RapicaCharge struct {
	Date    time.Time // 積増日付
	Charge  int       // 積増金額
	Premier int       // プレミア
	Company int       // 事業者
}

// *** RapiCa メソッド
// カード名
func (rapica *RapiCa) Name() string {
	return "RapiCa"
}

// システムコード
func (rapica *RapiCa) SystemCode() uint64 {
	return C.FELICA_POLLING_RAPICA
}

// カード情報を読込む
func (rapica *RapiCa) Read(cardinfo CardInfo) {
	if rapica.Info.Company != 0 {
		// 読込済みなら何もしない
		return
	}

	// システムデータの取得
	currsys := cardinfo.SysInfo(rapica.SystemCode())

	// RapiCa発行情報
	info := (*C.rapica_info_t)(currsys.SvcDataPtr(C.FELICA_SC_RAPICA_INFO, 0))
	i_time := C.rapica_info_date(info)

	rapica.Info.Company = int(C.rapica_info_company(info))
	rapica.Info.Deposit = int(C.rapica_info_deposit(info))
	rapica.Info.Date = time.Unix(int64(i_time), 0)

	// RapiCa属性情報(1)
	attr1 := (*C.rapica_attr1_t)(currsys.SvcDataPtr(C.FELICA_SC_RAPICA_ATTR, 0))
	a_time := C.rapica_attr_time(attr1)

	rapica.Attr.DateTime = time.Unix(int64(a_time), 0)
	rapica.Attr.Company = int(C.rapica_attr_company(attr1))
	rapica.Attr.TicketNo = int(attr1.ticketno)
	rapica.Attr.Busstop = int(C.rapica_attr_busstop(attr1))
	rapica.Attr.Busline = int(C.rapica_attr_busline(attr1))
	rapica.Attr.Busno = int(C.rapica_attr_busno(attr1))

	// RapiCa属性情報(2)
	attr2 := (*C.rapica_attr2_t)(currsys.SvcDataPtr(C.FELICA_SC_RAPICA_ATTR, 1))
	rapica.Attr.Kind = int(C.rapica_attr_kind(attr2))
	rapica.Attr.Amount = int(C.rapica_attr_amount(attr2))
	rapica.Attr.Premier = int(C.rapica_attr_premier(attr2))
	rapica.Attr.Point = int(C.rapica_attr_point(attr2))
	rapica.Attr.No = int(C.rapica_attr_no(attr2))
	rapica.Attr.OnBusstop = int(attr2.on_busstop)
	rapica.Attr.OffBusstop = int(attr2.off_busstop)

	// RapiCa属性情報(3)
	attr3 := (*C.rapica_attr3_t)(currsys.SvcDataPtr(C.FELICA_SC_RAPICA_ATTR, 2))
	rapica.Attr.Payment = int(C.rapica_attr_payment(attr3))

	// RapiCa属性情報(4)
	attr4 := (*C.rapica_attr4_t)(currsys.SvcDataPtr(C.FELICA_SC_RAPICA_ATTR, 3))
	rapica.Attr.Point2 = int(C.rapica_attr_point2(attr4))

	// RapiCa利用履歴
	last_time := C.time_t(rapica.Attr.DateTime.Unix())

	for i, _ := range currsys.SvcData(C.FELICA_SC_RAPICA_VALUE) {
		history := (*C.rapica_value_t)(currsys.SvcDataPtr(C.FELICA_SC_RAPICA_VALUE, i))
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
	for i, _ := range currsys.SvcData(C.FELICA_SC_RAPICA_CHARGE) {
		charge := (*C.rapica_charge_t)(currsys.SvcDataPtr(C.FELICA_SC_RAPICA_CHARGE, i))
		c_time := C.rapica_charge_date(charge)
		if c_time == 0 {
			continue
		}

		raw := RapicaCharge{}
		raw.Date = time.Unix(int64(c_time), 0)
		raw.Charge = int(C.rapica_charge_charge(charge))
		raw.Premier = int(C.rapica_charge_premier(charge))
		raw.Company = int(C.rapica_charge_company(charge))

		rapica.Charges = append(rapica.Charges, &raw)
	}
}

// カード情報を表示する
func (rapica *RapiCa) ShowInfo(cardinfo CardInfo, extend bool) {
	// テーブルデータの読込み
	if rapica_tables == nil {
		rapica_tables, _ = load_yaml("rapica.yml")
	}

	// データの読込み
	rapica.Read(cardinfo)

	// 表示
	attr := rapica.Attr

	fmt.Printf(`[発行情報]
  事業者: %v
  発行日: %s
  デポジット金額: %d円
`, rapica.Info.CompanyName(), rapica.Info.Date.Format("2006-01-02"), rapica.Info.Deposit)

	fmt.Println()
	fmt.Printf(`[属性情報]
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

	if extend {
		fmt.Println()
		fmt.Println("[利用履歴（元データ）]")
		fmt.Println("      日時      利用種別     残額         事業者 系統 / 停留所 (装置)")
		fmt.Println("  -------------------------------------------------------------------------------------")
		for _, value := range rapica.Hist {
			fmt.Printf("   %s    %v  %8d円    %v %v / %v (%d)\n",
				value.DateTime.Format("01/02 15:04"),
				value.KindName(),
				value.Amount,
				value.CompanyName(),
				value.BuslineName(),
				value.BusstopName(),
				value.Busno)
		}
	}

	fmt.Println()
	fmt.Println("[利用履歴]")
	fmt.Println("          日時       利用種別      利用料金        残額         事業者 系統 / 停留所 (装置)")
	fmt.Println("  ----------------------------------------------------------------------------------------------------------------------")
	for _, value := range rapica.Hist {
		disp_payment := "---"
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

		fmt.Printf("   %s    %v %14s\t%5d円    %v %v / %v (%d)\n",
			value.DateTime.Format("2006-01-02 15:04"), value.KindName(), disp_payment, value.Amount,
			value.CompanyName(), value.BuslineName(), disp_busstop, value.Busno)
	}

	fmt.Println()
	fmt.Println("[積増情報]")
	fmt.Println("      日時       チャージ   プレミア    事業者")
	fmt.Println("  ------------------------------------------------")
	for _, raw := range rapica.Charges {
		fmt.Printf("   %s %8d円 %8d円    %v\n",
			raw.Date.Format("2006-01-02"), raw.Charge, raw.Premier, raw.CompanyName())
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
	return disp_name(rapica_tables, name, value, base, opt_values...)
}
