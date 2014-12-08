#include "nicepass_get.h"


// バイト列を int に変換する
static int bytes_to_int(const uint8_t bytes[], size_t len) {
	int value = 0;

	for (size_t i = 0; i < len; i++) {
		value = (value << 8) + bytes[i];
	}

	return value;
}

// *** nice-pass残額
// チャージ金額残額
int nicepass_amount_charge(const nicepass_amount_t *amount) {
	return bytes_to_int(amount->charge, sizeof(amount->charge));
}

// プレミアム残額種別
int nicepass_amount_premium_kind(const nicepass_amount_t *amount) {
	return nicepass_premium_kind(amount);
}

// プレミアム残額
int nicepass_amount_premium(const nicepass_amount_t *amount) {
	return nicepass_premium(amount);
}

// *** nice-pass属性情報(2)
// 乗車日時
time_t nicepass_attr_in_time(const nicepass_attr2_t *attr) {
	int day = nicepass_day(attr);
	if (day == 0) return 0;

	struct tm tm = {
		//.tm_sec = nicepass_sec(attr->in_time),
		.tm_min = nicepass_min(attr->in_time),
		.tm_hour = nicepass_hour(attr->in_time),
		.tm_mday = day,
		.tm_mon = nicepass_month(attr) - 1,
		.tm_year = nicepass_year(attr) + 2000 - 1900,
	};

	return mktime(&tm);
}

// 降車日時
time_t nicepass_attr_out_time(const nicepass_attr2_t *attr) {
	int day = nicepass_day(attr);
	if (day == 0) return 0;

	struct tm tm = {
		//.tm_sec = nicepass_sec(attr->out_time),
		.tm_min = nicepass_min(attr->out_time),
		.tm_hour = nicepass_hour(attr->out_time),
		.tm_mday = day,
		.tm_mon = nicepass_month(attr) - 1,
		.tm_year = nicepass_year(attr) + 2000 - 1900,
	};

	return mktime(&tm);
}

// 使用装置
int nicepass_attr_type(const nicepass_attr2_t *attr) {
	return nicepass_type(attr);
}

// 処理種別
int nicepass_attr_proc(const nicepass_attr2_t *attr) {
	return nicepass_proc(attr);
}

// 直近利用金額
int nicepass_attr_use(const nicepass_attr2_t *attr) {
	return bytes_to_int(attr->use, sizeof(attr->use)) * 10;
}

// 直近残額
int nicepass_attr_balance(const nicepass_attr2_t *attr) {
	return nicepass_balance(attr);
}

// *** nice-pass属性情報(3)
// 乗車駅
int nicepass_attr_in_station(const nicepass_attr3_t *attr) {
	return nicepass_in_station(attr);
}

// 降車駅
int nicepass_attr_out_station(const nicepass_attr3_t *attr) {
	return nicepass_out_station(attr);
}

// 取引通番
int nicepass_attr_no(const nicepass_attr3_t *attr) {
	return bytes_to_int(attr->no, sizeof(attr->no));
}

// *** nice-pass利用履歴
// 処理日時
time_t nicepass_value_datetime(const nicepass_value_t *value) {
	int day = nicepass_day(value);
	if (day == 0) return 0;

	int min= nicepass_out_time(value);

	struct tm tm = {
		.tm_min = min % 60,
		.tm_hour = min / 60,
		.tm_mday = day,
		.tm_mon = nicepass_month(value) - 1,
		.tm_year = nicepass_year(value) + 2000 - 1900,
	};

	return mktime(&tm);
}

// 装置番号
int nicepass_value_train(const nicepass_value_t *value) {
	return nicepass_train(value);
}

// 乗車駅
int nicepass_value_in_station(const nicepass_value_t *value) {
	return nicepass_in_station(value);
}

// 降車駅
int nicepass_value_out_station(const nicepass_value_t *value) {
	return nicepass_out_station(value);
}

// 使用装置
int nicepass_value_type(const nicepass_value_t *value) {
	return nicepass_type(value);
}

// 処理種別
int nicepass_value_proc(const nicepass_value_t *value) {
	return nicepass_type(value);
}

// 利用金額種別
int nicepass_value_use_kind(const nicepass_value_t *value) {
	return nicepass_use_kind(value);
}

// 利用金額
int nicepass_value_use(const nicepass_value_t *value) {
	int16_t use = nicepass_use(value);

	if (use & 0x400) {
		use |= 0xf000;
	}

	return use * 10;
}

// 残額
int nicepass_value_balance(const nicepass_value_t *value) {
	return nicepass_balance(value);
}
