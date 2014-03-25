#include <time.h>
#include "rapica.h"

// *** RapiCa発行情報
int rapica_info_company(rapica_info_t *info);		// 事業者
time_t rapica_info_date(rapica_info_t *info);		// 発行日
int rapica_info_deposit(rapica_info_t *info);		// デポジット

// *** RapiCa属性情報
time_t rapica_attr_time(rapica_attr1_t *attr);		// 直近処理日時
int rapica_attr_company(rapica_attr1_t *attr);		// 事業者
int rapica_attr_busstop(rapica_attr1_t *attr);		// 停留所
int rapica_attr_busline(rapica_attr1_t *attr);		// 系統
int rapica_attr_busno(rapica_attr1_t *attr);		// 装置・車号？
int rapica_attr_kind(rapica_attr2_t *attr);			// 利用種別
int rapica_attr_amount(rapica_attr2_t *attr);		// 残額
int rapica_attr_premier(rapica_attr2_t *attr);		// プレミア
int rapica_attr_point(rapica_attr2_t *attr);		// ポイント
int rapica_attr_no(rapica_attr2_t *attr);			// 取引連番
int rapica_attr_payment(rapica_attr3_t *attr);		// 利用金額
int rapica_attr_point2(rapica_attr4_t *attr);		// ポイント？

// *** RapiCa履歴データ
time_t rapica_value_datetime(rapica_value_t *value, time_t last_time);	// 処理日時
int rapica_value_busstop(rapica_value_t *value);	// 停留所
int rapica_value_busline(rapica_value_t *value);	// 系統
int rapica_value_busno(rapica_value_t *value);		// 装置
int rapica_value_amount(rapica_value_t *value);		// 残額

// *** RapiCa積増データ
time_t rapica_charge_date(rapica_charge_t *charge);	// 積増日付
int rapica_charge_charge(rapica_charge_t *charge);	// 積増金額
int rapica_charge_premier(rapica_charge_t *charge);	// プレミア
int rapica_charge_company(rapica_charge_t *charge);	// 事業者
