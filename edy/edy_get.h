#include <time.h>
#include "edy.h"

// *** Edy残額情報（最終利用状況）
int edy_last_rest(const edy_last_t *last);				// 残額(LE)
int edy_last_use(const edy_last_t *last);				// 直近使用金額(LE) チャージのときは更新されない場合がある
int edy_last_no(const edy_last_t *last);				// 取引通番(LE)

// *** Edy履歴データ
time_t edy_value_datetime(const edy_value_t *value);	// 処理日時
int edy_value_type(const edy_value_t *value);			// タイプ
int edy_value_use(const edy_value_t *value);			// 入金／出金
int edy_value_rest(const edy_value_t *value);			// 残額
int edy_value_no(const edy_value_t *value);				// 連番
