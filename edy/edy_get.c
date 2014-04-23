#include "edy_get.h"


// バイト列(BE)を int に変換する
static int bytes_to_int(const uint8_t bytes[], size_t len) {
	int value = 0;

	for (size_t i = 0; i < len; i++) {
		value = (value << 8) + bytes[i];
	}

	return value;
}


// バイト列(LE)を int に変換する
static int le_to_int(const uint8_t bytes[], size_t len) {
	int value = 0;

	for (int i = len-1; 0 <= i; i--) {
		value = (value << 8) + bytes[i];
	}

	return value;
}


// *** Edy残額情報（最終利用状況）
// 残額(LE)
int edy_last_rest(const edy_last_t *last) {
	return le_to_int(last->rest, sizeof(last->rest));
}

// 直近使用金額(LE) チャージのときは更新されない場合がある
int edy_last_use(const edy_last_t *last) {
	return le_to_int(last->use, sizeof(last->use));
}

// 取引通番(LE)
int edy_last_no(const edy_last_t *last) {
	return le_to_int(last->no, sizeof(last->no));
}


// *** Edy履歴データ
// 処理日時
time_t edy_value_datetime(const edy_value_t *value) {
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
int edy_value_type(const edy_value_t *value) {
	return value->type;
}

// 入金／出金
int edy_value_use(const edy_value_t *value) {
	return edy_use(value);
}

// 残額
int edy_value_rest(const edy_value_t *value) {
	return edy_rest(value);
}

// 連番
int edy_value_no(const edy_value_t *value) {
	return bytes_to_int(value->no, sizeof(value->no));
}
