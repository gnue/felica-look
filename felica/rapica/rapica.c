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
int rapica_attr_on_busstop(rapica_attr2_t *attr) {
	return attr->on_busstop;
}

// 降車停留所(整理券)番号
int rapica_attr_off_busstop(rapica_attr2_t *attr) {
	return attr->off_busstop;
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
	if (rapica_is_iwasaki(value)) {
		// いわさきグループ
		return bytes_to_int(value->as.iwasaki.busstop, sizeof(value->as.iwasaki.busstop));
	} else {
		// Rapica加盟局社
		return bytes_to_int(value->as.rapica.busstop, sizeof(value->as.rapica.busstop));
	}
}

// 系統
int rapica_value_busline(rapica_value_t *value) {
	if (rapica_is_iwasaki(value)) {
		// いわさきグループ
		return bytes_to_int(value->as.iwasaki.busline, sizeof(value->as.iwasaki.busline));
	} else {
		// Rapica加盟局社
		return bytes_to_int(value->as.rapica.busline, sizeof(value->as.rapica.busline));
	}
}

// 装置
int rapica_value_busno(rapica_value_t *value) {
	if (rapica_is_iwasaki(value)) {
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
