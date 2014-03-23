#include "edy_get.h"


// バイト列を int に変換する
static int bytes_to_int(const uint8_t bytes[], size_t len) {
	int value = 0;

	for (size_t i = 0; i < len; i++) {
		value = (value << 8) + bytes[i];
	}

	return value;
}


// *** Edy履歴データ
// 処理日時
time_t edy_value_datetime(edy_value_t *value) {
	int days = edy_days(value);	// 累積日数（2000年から）
	int sec = edy_sec(value);

	if (days == 0 && sec == 0) return 0;

	struct tm tm = {
		.tm_mday = 1,
		.tm_year = 2000 - 1900,
	};

	time_t t = mktime(&tm);

	t += days * 24 * 60 * 60;
	t += sec;

	return t;
}

// タイプ
int edy_value_type(edy_value_t *value) {
	return value->type;
}

// 入金／出金
int edy_value_use(edy_value_t *value) {
	return edy_use(value);
}

// 残額
int edy_value_rest(edy_value_t *value) {
	return edy_rest(value);
}

// 連番
int edy_value_no(edy_value_t *value) {
	return bytes_to_int(value->no, sizeof(value->no));
}
