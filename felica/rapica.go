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
	info    RapicaInfo      // 発行情報
	attr    RapicaAttr      // 属性情報
	hist    []*RapicaValue  // 利用履歴
	charges []*RapicaCharge // 積増情報
}

// RapiCa発行情報データ
type RapicaInfo struct {
	date    time.Time // 発行日
	company int       // 事業者
	deposit int       // デポジット金額
}

// RapiCa属性情報データ
type RapicaAttr struct {
	datetime      time.Time // 直近処理日時
	company       int       // 事業者
	ticketno      int       // 整理券番号
	busstop       int       // 停留所
	busline       int       // 系統
	busno         int       // 装置
	kind          int       // 利用種別
	amount        int       // 残額
	premier       int       // プレミア
	point         int       // ポイント
	no            int       // 取引連番
	start_busstop int       // 乗車停留所(整理券)番号
	end_busstop   int       // 降車停留所(整理券)番号
	payment       int       // 利用金額
	point2        int       // ポイント？
}

// Rapica利用履歴データ
type RapicaValue struct {
	datetime time.Time // 処理日時
	company  int       // 事業者
	busstop  int       // 停留所
	busline  int       // 系統
	busno    int       // 装置
	kind     int       // 利用種別
	amount   int       // 残額

	payment  int // 利用料金（積増の場合はマイナス）
	st_value int // 対応する乗車データ
	ed_value int // 対応する降車データ
}

// Rapica積増情報データ
type RapicaCharge struct {
	date    time.Time // 積増日付
	charge  int       // 積増金額
	premier int       // プレミア
	company int       // 事業者
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
func (rapica *RapiCa) Read(cardinfo *CardInfo) {
	if rapica.info.company != 0 {
		// 読込済みなら何もしない
		return
	}

	// システムデータの取得
	currsys := cardinfo.sysinfo(rapica.SystemCode())

	// RapiCa発行情報
	info := (*C.rapica_info_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_INFO, 0))
	i_time := C.rapica_info_date(info)

	rapica.info.company = int(C.rapica_info_company(info))
	rapica.info.deposit = int(C.rapica_info_deposit(info))
	rapica.info.date = time.Unix(int64(i_time), 0)

	// RapiCa属性情報(1)
	attr1 := (*C.rapica_attr1_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 0))
	a_time := C.rapica_attr_time(attr1)

	rapica.attr.datetime = time.Unix(int64(a_time), 0)
	rapica.attr.company = int(C.rapica_attr_company(attr1))
	rapica.attr.ticketno = int(C.rapica_attr_ticketno(attr1))
	rapica.attr.busstop = int(C.rapica_attr_busstop(attr1))
	rapica.attr.busline = int(C.rapica_attr_busline(attr1))
	rapica.attr.busno = int(C.rapica_attr_busno(attr1))

	// RapiCa属性情報(2)
	attr2 := (*C.rapica_attr2_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 1))
	rapica.attr.kind = int(C.rapica_attr_kind(attr2))
	rapica.attr.amount = int(C.rapica_attr_amount(attr2))
	rapica.attr.premier = int(C.rapica_attr_premier(attr2))
	rapica.attr.point = int(C.rapica_attr_point(attr2))
	rapica.attr.no = int(C.rapica_attr_no(attr2))
	rapica.attr.start_busstop = int(C.rapica_attr_start_busstop(attr2))
	rapica.attr.end_busstop = int(C.rapica_attr_end_busstop(attr2))

	// RapiCa属性情報(3)
	attr3 := (*C.rapica_attr3_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 2))
	rapica.attr.payment = int(C.rapica_attr_payment(attr3))

	// RapiCa属性情報(4)
	attr4 := (*C.rapica_attr4_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 3))
	rapica.attr.point2 = int(C.rapica_attr_point2(attr4))

	// RapiCa利用履歴
	last_time := C.time_t(rapica.attr.datetime.Unix())

	for i, _ := range currsys.svcdata(C.FELICA_SC_RAPICA_VALUE) {
		history := (*C.rapica_value_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_VALUE, i))
		h_time := C.rapica_value_datetime(history, last_time)
		if h_time == 0 {
			continue
		}

		value := RapicaValue{}
		value.datetime = time.Unix(int64(h_time), 0)
		value.company = int(C.rapica_value_company(history))
		value.busstop = int(C.rapica_value_busstop(history))
		value.busline = int(C.rapica_value_busline(history))
		value.busno = int(C.rapica_value_busno(history))
		value.kind = int(C.rapica_value_kind(history))
		value.amount = int(C.rapica_value_amount(history))
		value.st_value = -1
		value.ed_value = -1

		rapica.hist = append(rapica.hist, &value)
		last_time = h_time
	}

	// 乗車データと降車データの関連付けをする
	for i, value := range rapica.hist[:len(rapica.hist)-1] {
		pre_data := rapica.hist[i+1]

		if value.kind == C.RAPICA_KIND_GETOFF {
			// 降車
			for j, v := range rapica.hist[i+1:] {
				if v.kind == C.RAPICA_KIND_GETON {
					// 乗車を見つけた
					value.st_value = i + 1 + j
					v.ed_value = i
					break
				}
			}
		}

		value.payment = pre_data.amount - value.amount
	}

	// RapiCa積増情報
	for i, _ := range currsys.svcdata(C.FELICA_SC_RAPICA_CHARGE) {
		charge := (*C.rapica_charge_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_CHARGE, i))
		c_time := C.rapica_charge_date(charge)
		if c_time == 0 {
			continue
		}

		raw := RapicaCharge{}
		raw.date = time.Unix(int64(c_time), 0)
		raw.charge = int(C.rapica_charge_charge(charge))
		raw.premier = int(C.rapica_charge_premier(charge))
		raw.company = int(C.rapica_charge_company(charge))

		rapica.charges = append(rapica.charges, &raw)
	}
}

// カード情報を表示する
func (rapica *RapiCa) ShowInfo(cardinfo *CardInfo, extend bool) {
	// テーブルデータの読込み
	if rapica_tables == nil {
		rapica_tables, _ = load_yaml("rapica.yml")
	}

	// データの読込み
	rapica.Read(cardinfo)

	// 表示
	attr := rapica.attr

	fmt.Printf(`[発行情報]
  事業者: %v
  発行日: %s
  デポジット金額: %d円
`, rapica.info.company_name(), rapica.info.date.Format("2006-01-02"), rapica.info.deposit)

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
`, attr.datetime.Format("2006-01-02 15:04"),
		attr.company_name(),
		attr.ticketno, attr.busstop, attr.busline, attr.busno,
		attr.kind_name(),
		attr.amount, attr.premier, attr.point, attr.no, attr.start_busstop, attr.end_busstop,
		attr.payment, attr.point2)

	if extend {
		fmt.Println()
		fmt.Println("[利用履歴（元データ）]")
		fmt.Println("      日時      利用種別     残額         事業者 系統 / 停留所 (装置)")
		fmt.Println("  -------------------------------------------------------------------------------------")
		for _, value := range rapica.hist {
			fmt.Printf("   %s    %v  %8d円    %v %v / %v (%d)\n",
				value.datetime.Format("01/02 15:04"),
				value.kind_name(),
				value.amount,
				value.company_name(),
				value.busline_name(),
				value.busstop_name(),
				value.busno)
		}
	}

	fmt.Println()
	fmt.Println("[利用履歴]")
	fmt.Println("          日時       利用種別      利用料金        残額         事業者 系統 / 停留所 (装置)")
	fmt.Println("  ----------------------------------------------------------------------------------------------------------------------")
	for _, value := range rapica.hist {
		disp_payment := "---"
		disp_busstop := value.busstop_name()

		if 0 <= value.ed_value && value.payment == 0 {
			// 対応する降車データがあり利用金額が 0 ならば表示しない
			continue
		}

		if value.payment < 0 {
			// 積増
			disp_payment = fmt.Sprintf("(+%d円)", -value.payment)
		} else if 0 < value.payment {
			disp_payment = fmt.Sprintf("%d円", value.payment)
		}

		if 0 <= value.st_value {
			st_value := rapica.hist[value.st_value]
			disp_busstop = fmt.Sprintf("%v -> %v", st_value.busstop_name(), disp_busstop)
		}

		fmt.Printf("   %s    %v %14s\t%5d円    %v %v / %v (%d)\n",
			value.datetime.Format("2006-01-02 15:04"), value.kind_name(), disp_payment, value.amount,
			value.company_name(), value.busline_name(), disp_busstop, value.busno)
	}

	fmt.Println()
	fmt.Println("[積増情報]")
	fmt.Println("      日時       チャージ   プレミア    事業者")
	fmt.Println("  ------------------------------------------------")
	for _, raw := range rapica.charges {
		fmt.Printf("   %s %8d円 %8d円    %v\n",
			raw.date.Format("2006-01-02"), raw.charge, raw.premier, raw.company_name())
	}
}

// *** RapicaInfo メソッド
// 事業者名
func (info *RapicaInfo) company_name() interface{} {
	return rapica_disp_name("ATTR_COMPANY", info.company, 4)
}

// *** RapicaAttr メソッド
// 事業者名
func (attr *RapicaAttr) company_name() interface{} {
	return rapica_disp_name("ATTR_COMPANY", attr.company, 4)
}

// 利用種別
func (attr *RapicaAttr) kind_name() interface{} {
	return rapica_disp_name("ATTR_KIND", attr.kind&0xff0000, 6, attr.kind)
}

// *** RapicaValue メソッド
// 利用種別
func (value *RapicaValue) kind_name() interface{} {
	return rapica_disp_name("HIST_KIND", value.kind, 2)
}

// 事業者名
func (value *RapicaValue) company_name() interface{} {
	return rapica_disp_name("HIST_COMPANY", value.company>>4, 2, value.company)
}

// 停留所
func (value *RapicaValue) busstop_name() interface{} {
	return rapica_disp_name("BUSSTOP", value.busstop, 6)
}

// 系統名
func (value *RapicaValue) busline_name() interface{} {
	return rapica_disp_name("BUSLINE", (value.busstop&0xff0000)+value.busline, 4, value.busline)
}

// *** RapicaCharge メソッド
func (charge *RapicaCharge) company_name() interface{} {
	return rapica_disp_name("ATTR_COMPANY", charge.company, 4)
}

// ***
// RapiCaテーブル
var rapica_tables map[interface{}]interface{}

// RapiCaテーブルを検索して表示用の文字列を返す
func rapica_disp_name(name string, value int, base int, opt_values ...int) interface{} {
	return disp_name(rapica_tables, name, value, base, opt_values...)
}
