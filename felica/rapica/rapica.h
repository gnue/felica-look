/*------------------------------------------------------------------------*/
/**
 * @file	rapica.h
 * @brief   RapiCa
 *
 * @author  M.Nukui
 * @date	2014-03-18
 *
 * Copyright (C) 2014 M.Nukui All rights reserved.
 */


#ifndef	PASORIKIT_RAPICA_H
#define	PASORIKIT_RAPICA_H


#include <stdint.h>


#define FELICA_POLLING_RAPICA			0x8194		///< RapiCa/鹿児島市交通局

#define FELICA_SC_RAPICA_INFO			0x000B		///< RapiCa発行情報・サービスコード
#define FELICA_SC_RAPICA_ATTR			0x004B		///< RapiCa属性情報・サービスコード
#define FELICA_SC_RAPICA_VALUE			0x008F		///< RapiCa利用履歴データ・サービスコード
#define FELICA_SC_RAPICA_CHARGE			0x00CF		///< RapiCa積増情報データ・サービスコード

#define RAPICA_KIND_CREATE				0x00		///< 作成？
#define RAPICA_KIND_REGISTER			0x10		///< 登録？
#define RAPICA_KIND_GETON				0x30		///< 乗車
#define RAPICA_KIND_CHARGE				0x40		///< 積増
#define RAPICA_KIND_GETOFF				0x41		///< 降車


#define rapica_is_iwasaki(v)			((v->company & 0xf0) == 0x40) ///< いわさきグループか？

// generated bitfeilds2macro.rb (rapica)

#define rapica_date(v)	((v->datetime[0] << 4) + (v->datetime[1] >> 4))		///< bits 0-11: 月*100+日
#define rapica_time(v)	(((v->datetime[1] & 0x0f) << 8) + v->datetime[2])	///< bits 12-23: 時*100+分


/// RapiCa発行データ
typedef struct {
	uint8_t		company[2];		///< 事業者
	uint8_t		unkown1[3];		///< 不明
	uint8_t		year;			///< 年(+2000)
	uint8_t		month;			///< 月
	uint8_t		day;			///< 日
	uint8_t		unkown2[4];		///< 不明
	uint8_t		deposit[2];		///< デポジット
	uint8_t		unkown3[2];		///< 不明
} rapica_info_t;


/// RapiCa属性データ(1)
typedef struct {
	uint8_t		year;			///< 年(+2000)
	uint8_t		month;			///< 月
	uint8_t		day;			///< 日
	uint8_t		hour;			///< 時
	uint8_t		minutes;		///< 分
	uint8_t		company[2];		///< 事業者
	uint8_t		ticketno;		///< 整理券番号
	uint8_t		busstop[3];		///< 停留所
	uint8_t		busline[2];		///< 系統
	uint8_t		busno[3];		///< 装置
} rapica_attr1_t;


/// RapiCa属性データ(2)
typedef struct {
	uint8_t		kind[3];		///< 利用種別
	uint8_t		amount[3];		///< 残額
	uint8_t		premier[2];		///< プレミア
	uint8_t		point[2];		///< ポイント
	uint8_t		unkown[1];		///< 不明
	uint8_t		no[3];			///< 取引連番
	uint8_t		start_busstop;	///< 乗車停留所(整理券)番号
	uint8_t		end_busstop;	///< 降車停留所(整理券)番号
} rapica_attr2_t;


/// RapiCa属性データ(3)
typedef struct {
	uint8_t		unkown1[1];		///< 不明
	uint8_t		payment[2];		///< 利用金額
	uint8_t		unkown2[13];	///< 不明
} rapica_attr3_t;


/// RapiCa属性データ(4)
typedef struct {
	uint8_t		point[2];		///< ポイント？
	uint8_t		unkown[14];		///< 不明
} rapica_attr4_t;


/// RapiCa履歴データ
typedef struct {
	uint8_t		datetime[3];	///< 日付(12bit)/時刻(12bit)

#pragma mark bitfeilds2macro(data)
#if 0
	date:12;	///< 月*100+日
	time:12;	///< 時*100+分
#endif

	uint8_t		company;			///< 事業者
	union {
		/// Rapica加盟局社
		struct {
			uint8_t		busstop[3];	///< 停留所
			uint8_t		busline[2];	///< 系統
			uint8_t		busno[3];	///< 装置
		} rapica;

		/// いわさきグループ
		struct {
			uint8_t		unkown[1];	///< 不明
			uint8_t		busstop[3];	///< 停留所
			uint8_t		busline[2];	///< 系統
			uint8_t		busno[2];	///< 装置
		} iwasaki;
	} as;
	uint8_t		kind;			///< 利用種別
	uint8_t		amount[3];	///< 残額
} rapica_value_t;


/// RapiCa積増データ
typedef struct {
	uint8_t		year;			///< 年(+2000)
	uint8_t		month;			///< 月
	uint8_t		day;			///< 日
	uint8_t		charge[2];		///< 積増金額
	uint8_t		premier[2];		///< プレミア
	uint8_t		unkown[7];		///< 不明
	uint8_t		company[2];		///< 事業者
} rapica_charge_t;


#endif /* PASORIKIT_RAPICA_H */
