#include <time.h>
#include "suica.h"

// *** Suica利用履歴
int suica_value_type(const suica_value_t *value);			// 端末種
int suica_value_proc(const suica_value_t *value);			// 処理
time_t suica_value_date(const suica_value_t *value);		// 利用年月日
int suica_value_in_station(const suica_value_t *value);		// 入場駅（線区コード、駅順コード）
int suica_value_out_station(const suica_value_t *value);	// 出場駅（線区コード、駅順コード）
int suica_value_balance(const suica_value_t *value);		// 残額（リトルエンディアン）
int suica_value_no(const suica_value_t *value);				// 連番
int suica_value_region(const suica_value_t *value);			// 処理
