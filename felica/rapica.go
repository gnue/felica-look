package felica

import (
	"fmt"
	"time"
)

/*
#include <time.h>
#include "rapica.h"

int bytes_to_int(const uint8_t bytes[], size_t len) {
	int value = 0;

	for (size_t i = 0; i < len; i++) {
		value = (value << 8) + bytes[i];
	}

	return value;
}

// *** RapiCa発行情報
// 事業者
int rapica_info_company(rapica_info_t *info) {
	return bytes_to_int(info->company, sizeof(info->company));
}

// 発行日
time_t rapica_info_date(rapica_info_t *info) {
	if (info->day == 0) return 0;

	struct tm tm = {
		.tm_mday = info->day,
		.tm_mon = info->month - 1,
		.tm_year = info->year + 2000 - 1900,
	};

	return mktime(&tm);
}

// デポジット
int rapica_info_deposit(rapica_info_t *info) {
	return bytes_to_int(info->deposit, sizeof(info->deposit));
}

// *** RapiCa属性情報(1)
// 直近処理日時
time_t rapica_attr_time(rapica_attr1_t *attr) {
	if (attr->day == 0) return 0;

	struct tm tm = {
		.tm_min = attr->minutes,
		.tm_hour = attr->hour,
		.tm_mday = attr->day,
		.tm_mon = attr->month - 1,
		.tm_year = attr->year + 2000 - 1900,
	};

	return mktime(&tm);
}

// 事業者
int rapica_attr_company(rapica_attr1_t *attr) {
	return bytes_to_int(attr->company, sizeof(attr->company));
}

// 整理券番号
int rapica_attr_ticketno(rapica_attr1_t *attr) {
	return attr->ticketno;
}

// 停留所
int rapica_attr_busstop(rapica_attr1_t *attr) {
	return bytes_to_int(attr->busstop, sizeof(attr->busstop));
}

// 系統
int rapica_attr_busline(rapica_attr1_t *attr) {
	return bytes_to_int(attr->busline, sizeof(attr->busline));
}

// 装置・車号？
int rapica_attr_busno(rapica_attr1_t *attr) {
	return bytes_to_int(attr->busno, sizeof(attr->busno));
}

// *** RapiCa属性情報(2)
// 利用種別
int rapica_attr_kind(rapica_attr2_t *attr) {
	return bytes_to_int(attr->kind, sizeof(attr->kind));
}

// 残額
int rapica_attr_amount(rapica_attr2_t *attr) {
	return bytes_to_int(attr->amount, sizeof(attr->amount));
}

// プレミア
int rapica_attr_premier(rapica_attr2_t *attr) {
	return bytes_to_int(attr->premier, sizeof(attr->premier));
}

// ポイント
int rapica_attr_point(rapica_attr2_t *attr) {
	return bytes_to_int(attr->point, sizeof(attr->point));
}

// 取引連番
int rapica_attr_no(rapica_attr2_t *attr) {
	return bytes_to_int(attr->no, sizeof(attr->no));
}

// 乗車停留所(整理券)番号
int rapica_attr_start_busstop(rapica_attr2_t *attr) {
	return attr->start_busstop;
}

// 降車停留所(整理券)番号
int rapica_attr_end_busstop(rapica_attr2_t *attr) {
	return attr->end_busstop;
}

// *** RapiCa属性情報(3)
// 利用金額
int rapica_attr_payment(rapica_attr3_t *attr) {
	return bytes_to_int(attr->payment, sizeof(attr->payment));
}

// *** RapiCa属性情報(4)
// ポイント？
int rapica_attr_point2(rapica_attr4_t *attr) {
	return bytes_to_int(attr->point, sizeof(attr->point));
}

// *** RapiCa履歴データ
// 処理日時
time_t rapica_value_datetime(rapica_value_t *value, time_t last_time) {
	struct tm last_tm;
	int date = rapica_date(value);
	int time = rapica_time(value);

	localtime_r(&last_time, &last_tm);
	int last_date = (last_tm.tm_mon + 1) * 100 + last_tm.tm_mday;
	int year = last_tm.tm_year;

	if (date > last_date) {
		// 年をまたいでいるので前年にする
		year--;
	}

	struct tm tm = {
		.tm_min = time % 100,
		.tm_hour = time / 100,
		.tm_mday = date % 100,
		.tm_mon = date / 100 - 1,
		.tm_year = year,
	};

	return mktime(&tm);
}

// 事業者
int rapica_value_company(rapica_value_t *value) {
	return value->company;
}

// 停留所
int rapica_value_busstop(rapica_value_t *value) {
	if (value->company == 0x40) {
		// いわさきグループ
		return bytes_to_int(value->as.iwasaki.busstop, sizeof(value->as.iwasaki.busstop));
	} else {
		// Rapica加盟局社
		return bytes_to_int(value->as.rapica.busstop, sizeof(value->as.rapica.busstop));
	}
}

// 系統
int rapica_value_busline(rapica_value_t *value) {
	if (value->company == 0x40) {
		// いわさきグループ
		return bytes_to_int(value->as.iwasaki.busline, sizeof(value->as.iwasaki.busline));
	} else {
		// Rapica加盟局社
		return bytes_to_int(value->as.rapica.busline, sizeof(value->as.rapica.busline));
	}
}

// 装置
int rapica_value_busno(rapica_value_t *value) {
	if (value->company == 0x40) {
		// いわさきグループ
		return bytes_to_int(value->as.iwasaki.busno, sizeof(value->as.iwasaki.busno));
	} else {
		// Rapica加盟局社
		return bytes_to_int(value->as.rapica.busno, sizeof(value->as.rapica.busno));
	}
}

// 利用種別
int rapica_value_kind(rapica_value_t *value) {
	return value->kind;
}

// 残額
int rapica_value_amount(rapica_value_t *value) {
	return bytes_to_int(value->amount, sizeof(value->amount));
}

// *** RapiCa積増データ
// 積増日付
time_t rapica_charge_date(rapica_charge_t *charge) {
	if (charge->day == 0) return 0;

	struct tm tm = {
		.tm_mday = charge->day,
		.tm_mon = charge->month - 1,
		.tm_year = charge->year + 2000 - 1900,
	};

	return mktime(&tm);
}

// 積増金額
int rapica_charge_charge(rapica_charge_t *charge) {
	return bytes_to_int(charge->charge, sizeof(charge->charge));
}

// プレミア
int rapica_charge_premier(rapica_charge_t *charge) {
	return bytes_to_int(charge->premier, sizeof(charge->premier));
}

// 事業者
int rapica_charge_company(rapica_charge_t *charge) {
	return bytes_to_int(charge->company, sizeof(charge->company));
}

*/
import "C"

// RapiCa/鹿児島市交通局
type RapiCa struct {
}

// カード名
func (rapica *RapiCa) Name() string {
	return "RapiCa"
}

// システムコード
func (rapica *RapiCa) SystemCode() uint64 {
	return C.FELICA_POLLING_RAPICA
}

// カード情報を表示する
func (rapica *RapiCa) ShowInfo(cardinfo *CardInfo, extend bool) {

	// システムデータの取得
	currsys := cardinfo.sysinfo(rapica.SystemCode())

	// RapiCa発行情報
	info := (*C.rapica_info_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_INFO, 0))

	i_company := C.rapica_info_company(info)
	i_time := C.rapica_info_date(info)
	i_deposit := C.rapica_info_deposit(info)

	i_date := time.Unix(int64(i_time), 0)

	// RapiCa属性情報(1)
	attr1 := (*C.rapica_attr1_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 0))
	a_time := C.rapica_attr_time(attr1)
	a_company := C.rapica_attr_company(attr1)
	a_ticketno := C.rapica_attr_ticketno(attr1)
	a_busstop := C.rapica_attr_busstop(attr1)
	a_busline := C.rapica_attr_busline(attr1)
	a_busno := C.rapica_attr_busno(attr1)

	a_datetime := time.Unix(int64(a_time), 0)

	// RapiCa属性情報(2)
	attr2 := (*C.rapica_attr2_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 1))
	a_kind := C.rapica_attr_kind(attr2)
	a_amount := C.rapica_attr_amount(attr2)
	a_premier := C.rapica_attr_premier(attr2)
	a_point := C.rapica_attr_point(attr2)
	a_no := C.rapica_attr_no(attr2)
	a_start_busstop := C.rapica_attr_start_busstop(attr2)
	a_end_busstop := C.rapica_attr_end_busstop(attr2)

	// RapiCa属性情報(3)
	attr3 := (*C.rapica_attr3_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 2))
	a_payment := C.rapica_attr_payment(attr3)

	// RapiCa属性情報(4)
	attr4 := (*C.rapica_attr4_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_ATTR, 3))
	a_point2 := C.rapica_attr_point2(attr4)

	// 表示
	fmt.Printf(`[発行情報]
  事業者: 0x%04X
  発行日: %s
  デポジット金額: %d円
`, i_company, i_date.Format("2006-01-02"), i_deposit)

	fmt.Printf(`[属性情報]
  直近処理日時:	%s
  事業者:	0x%04X
  整理券番号:	%d
  停留所:	0x%06X
  系統:		0x%04X
  装置・車号？:	%d
  利用種別:	0x%04X
  残額:		%d円
  プレミア:	%d円
  ポイント:	%dpt
  取引連番:	%d
  乗車停留所(整理券)番号: %d
  降車停留所(整理券)番号: %d
  利用金額:	%d円
  ポイント？:	%dpt
`, a_datetime.Format("2006-01-02 15:04"), a_company, a_ticketno, a_busstop, a_busline, a_busno,
		a_kind, a_amount, a_premier, a_point, a_no, a_start_busstop, a_end_busstop,
		a_payment, a_point2)

	// RapiCa利用履歴
	fmt.Println("[利用履歴]")
	last_time := C.time_t(a_datetime.Unix())
	for i, _ := range currsys.svcdata(C.FELICA_SC_RAPICA_VALUE) {
		history := (*C.rapica_value_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_VALUE, i))
		h_time := C.rapica_value_datetime(history, last_time)
		if h_time == 0 {
			continue
		}

		h_company := C.rapica_value_company(history)
		h_busstop := C.rapica_value_busstop(history)
		h_busline := C.rapica_value_busline(history)
		h_busno := C.rapica_value_busno(history)
		h_kind := C.rapica_value_kind(history)
		h_amount := C.rapica_value_amount(history)

		h_datetime := time.Unix(int64(h_time), 0)

		fmt.Printf("  %s  0x%04X  残額:%5d円\t0x%04X 0x%04X / 0x%06X (%d)\n", h_datetime.Format("01/02 15:04"), h_kind, h_amount,
			h_company, h_busline, h_busstop, h_busno)
		last_time = h_time
	}

	// RapiCa積増情報
	fmt.Println("[積増情報]")
	for i, _ := range currsys.svcdata(C.FELICA_SC_RAPICA_CHARGE) {
		charge := (*C.rapica_charge_t)(currsys.svcdata_ptr(C.FELICA_SC_RAPICA_CHARGE, i))
		c_time := C.rapica_charge_date(charge)
		if c_time == 0 {
			continue
		}

		c_charge := C.rapica_charge_charge(charge)
		c_premier := C.rapica_charge_premier(charge)
		c_company := C.rapica_charge_company(charge)

		c_date := time.Unix(int64(c_time), 0)

		fmt.Printf("  %s 積増金額:%d円 プレミア:%d円  0x%04X\n", c_date.Format("2006-01-02"), c_charge, c_premier, c_company)
	}
}
