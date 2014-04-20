#include <time.h>
#include "rapica.h"

// *** RapiCa発行情報
int rapica_info_company(const rapica_info_t *info);			// 事業者
time_t rapica_info_date(const rapica_info_t *info);			// 発行日
int rapica_info_deposit(const rapica_info_t *info);			// デポジット

// *** RapiCa属性情報
time_t rapica_attr_time(const rapica_attr1_t *attr);		// 直近処理日時
int rapica_attr_company(const rapica_attr1_t *attr);		// 事業者
int rapica_attr_busstop(const rapica_attr1_t *attr);		// 停留所
int rapica_attr_busline(const rapica_attr1_t *attr);		// 系統
int rapica_attr_busno(const rapica_attr1_t *attr);			// 装置・車号？
int rapica_attr_kind(const rapica_attr2_t *attr);			// 利用種別
int rapica_attr_amount(const rapica_attr2_t *attr);			// 残額
int rapica_attr_premier(const rapica_attr2_t *attr);		// プレミア
int rapica_attr_point(const rapica_attr2_t *attr);			// ポイント
int rapica_attr_no(const rapica_attr2_t *attr);				// 取引連番
int rapica_attr_payment(const rapica_attr3_t *attr);		// 利用金額
int rapica_attr_point2(const rapica_attr4_t *attr);			// ポイント？

// *** RapiCa履歴データ
time_t rapica_value_datetime(const rapica_value_t *value, time_t last_time);	// 処理日時
int rapica_value_busstop(const rapica_value_t *value);		// 停留所
int rapica_value_busline(const rapica_value_t *value);		// 系統
int rapica_value_busno(const rapica_value_t *value);		// 装置
int rapica_value_amount(const rapica_value_t *value);		// 残額

// *** RapiCa積増データ
time_t rapica_charge_date(const rapica_charge_t *charge);	// 積増日付
int rapica_charge_charge(const rapica_charge_t *charge);	// 積増金額
int rapica_charge_premier(const rapica_charge_t *charge);	// プレミア
int rapica_charge_company(const rapica_charge_t *charge);	// 事業者
