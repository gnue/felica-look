#include <time.h>
#include "nicepass.h"

// *** nice-pass利用履歴
time_t nicepass_value_datetime(const nicepass_value_t *value);	// 処理日時
int nicepass_value_train(const nicepass_value_t *value);		// 装置番号
int nicepass_value_in_station(const nicepass_value_t *value);	// 入場駅
int nicepass_value_out_station(const nicepass_value_t *value);	// 出場駅
int nicepass_value_type(const nicepass_value_t *value);			// 使用装置
int nicepass_value_proc(const nicepass_value_t *value);			// 処理種別
int nicepass_value_use(const nicepass_value_t *value);			// 利用金額
int nicepass_value_balance(const nicepass_value_t *value);		// 残額
