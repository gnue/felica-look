#include <time.h>
#include "nicepass.h"

// *** nice-pass残額
int nicepass_amount_charge(const nicepass_amount_t *amount);		// チャージ金額残額
int nicepass_amount_premium_kind(const nicepass_amount_t *amount);	// プレミアム残額種別
int nicepass_amount_premium(const nicepass_amount_t *amount);		// プレミアム残額

// *** nice-pass属性情報
time_t nicepass_attr_in_time(const nicepass_attr2_t *attr);		// 乗車日時
time_t nicepass_attr_out_time(const nicepass_attr2_t *attr);	// 降車日時
int nicepass_attr_type(const nicepass_attr2_t *attr);			// 使用装置
int nicepass_attr_proc(const nicepass_attr2_t *attr);			// 処理種別
int nicepass_attr_use(const nicepass_attr2_t *attr);			// 直近利用金額
int nicepass_attr_balance(const nicepass_attr2_t *attr);		// 直近残額
int nicepass_attr_in_station(const nicepass_attr3_t *attr);		// 乗車駅
int nicepass_attr_out_station(const nicepass_attr3_t *attr);	// 降車駅
int nicepass_attr_no(const nicepass_attr3_t *attr);				// 取引通番

// *** nice-pass利用履歴
time_t nicepass_value_datetime(const nicepass_value_t *value);	// 処理日時
int nicepass_value_train(const nicepass_value_t *value);		// 装置番号
int nicepass_value_in_station(const nicepass_value_t *value);	// 乗車駅
int nicepass_value_out_station(const nicepass_value_t *value);	// 降車駅
int nicepass_value_type(const nicepass_value_t *value);			// 使用装置
int nicepass_value_proc(const nicepass_value_t *value);			// 処理種別
int nicepass_value_use_kind(const nicepass_value_t *value);		// 利用金額種別
int nicepass_value_use(const nicepass_value_t *value);			// 利用金額
int nicepass_value_balance(const nicepass_value_t *value);		// 残額
