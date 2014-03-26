#include <time.h>
#include "edy.h"

// *** Edy残額情報（最終利用状況）
int edy_last_rest(edy_last_t *last);			// 残額(LE)
int edy_last_use(edy_last_t *last);				// 直近使用金額(LE) チャージのときは更新されない場合がある
int edy_last_no(edy_last_t *last);				// 取引通番(LE)

// *** Edy履歴データ
time_t edy_value_datetime(edy_value_t *value);	// 処理日時
int edy_value_type(edy_value_t *value);			// タイプ
int edy_value_use(edy_value_t *value);			// 入金／出金
int edy_value_rest(edy_value_t *value);			// 残額
int edy_value_no(edy_value_t *value);			// 連番
