#include <time.h>
#include "suica.h"

// *** Suica利用履歴
int suica_value_type(suica_value_t *value);			// 端末種
int suica_value_proc(suica_value_t *value);			// 処理
time_t suica_value_date(suica_value_t *value);		// 利用年月日
int suica_value_in_station(suica_value_t *value);	// 入場駅（線区コード、駅順コード）
int suica_value_out_station(suica_value_t *value);	// 出場駅（線区コード、駅順コード）
int suica_value_balance(suica_value_t *value);		// 残額（リトルエンディアン）
int suica_value_no(suica_value_t *value);			// 連番
int suica_value_region(suica_value_t *value);		// 処理
