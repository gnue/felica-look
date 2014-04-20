#include "suica_get.h"


// バイト列を int に変換する
static int bytes_to_int(const uint8_t bytes[], size_t len) {
	int value = 0;

	for (size_t i = 0; i < len; i++) {
		value = (value << 8) + bytes[i];
	}

	return value;
}


// *** Suica利用履歴
// 端末種
int suica_value_type(const suica_value_t *value) {
	return value->type;
}

// 利用年月日
time_t suica_value_date(const suica_value_t *value) {
	int day = suica_day(value);
	if (day == 0) return 0;

	struct tm tm = {
		.tm_mday = day,
		.tm_mon = suica_month(value) - 1,
		.tm_year = suica_year(value) + 2000 - 1900,
	};

	return mktime(&tm);
}

// 入場駅（線区コード、駅順コード）
int suica_value_in_station(const suica_value_t *value) {
	return bytes_to_int(value->in_station, sizeof(value->in_station));
}

// 出場駅（線区コード、駅順コード）
int suica_value_out_station(const suica_value_t *value) {
	return bytes_to_int(value->out_station, sizeof(value->out_station));
}

// 残額（リトルエンディアン）
int suica_value_balance(const suica_value_t *value) {
	return suica_balance(value);
}

// 連番
int suica_value_no(const suica_value_t *value) {
	return bytes_to_int(value->no, sizeof(value->no));
}
