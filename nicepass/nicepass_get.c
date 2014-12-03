#include "nicepass_get.h"


// バイト列を int に変換する
static int bytes_to_int(const uint8_t bytes[], size_t len) {
	int value = 0;

	for (size_t i = 0; i < len; i++) {
		value = (value << 8) + bytes[i];
	}

	return value;
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
