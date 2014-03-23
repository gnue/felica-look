#include <time.h>
#include "edy.h"

// *** Edy履歴データ
time_t edy_value_datetime(edy_value_t *value);	// 処理日時
int edy_value_type(edy_value_t *value);			// タイプ
int edy_value_use(edy_value_t *value);			// 入金／出金
int edy_value_rest(edy_value_t *value);			// 残額
int edy_value_no(edy_value_t *value);			// 連番
