/*------------------------------------------------------------------------*/
/**
 * @file	nicepass.h
 * @brief   nice-pass
 *
 * @author  M.Nukui
 * @date	2014-12-03
 *
 * Copyright (C) 2014 M.Nukui All rights reserved.
 */


#ifndef	PASORIKIT_NICEPASS_H
#define	PASORIKIT_NICEPASS_H


#include <stdint.h>


#define FELICA_POLLING_NICEPASS			0x040F		///< nice-passシステムコード

#define FELICA_SC_NICEPASS_ATTR			0x020B		///< nice-pass属性情報・サービスコード
#define FELICA_SC_NICEPASS_VALUE		0x030F		///< nice-pass利用履歴データ・サービスコード


#define nicepass_year(v)		(v->date[0] >> 1)								///< 年の取得
#define nicepass_month(v)		(((v->date[0] & 1) << 3) + (v->date[1] >> 5))	///< 月の取得
#define nicepass_day(v)			(v->date[1] & 0x1F)								///< 日の取得
#define nicepass_out_time(v)	(v->out_time * 10)								///< 通算分数の取得
#define nicepass_train(v)		((v->train[0] << 8) + v->train[1])				///< 装置番号の取得
#define nicepass_balance(v)		((v->balance[0] << 8) + v->balance[1])			///< 残額の取得

// generated bitfeilds2macro.rb (nicepass)

#define nicepass_in_station(v)	((v->station[0] << 12) + (v->station[1] << 4) + (v->station[2] >> 4))	///< bits 0-19: 乗車駅
#define nicepass_out_station(v)	(((v->station[2] & 0x0f) << 16) + (v->station[3] << 8) + v->station[4])	///< bits 20-39: 降車駅
#define nicepass_type(v)		(v->proc >> 4)									///< bits 0-3: 使用装置
#define nicepass_proc(v)		(v->proc & 0x0f)								///< bits 4-7: 処理種別
#define nicepass_use(v)			(((v->use[0] & 0x0f) << 8) + v->use[1])			///< bits 4-15: 利用金額（支払いはマイナス）


/// nice-pass履歴データ
typedef struct {
	uint8_t		date[2];		///< 年(7bit)/月(4bit)/日(5bit)
	uint8_t		out_time;		///< 降車時刻 00:00からの通算分数/10
	uint8_t		train[2];		///< 装置番号
	uint8_t		station[5];		///< 乗降駅

#pragma mark bitfeilds2macro(station)
#if 0
	in_station:20;	///< 乗車駅
	out_station:20;	///< 降車駅
#endif

	uint8_t		proc;			///< 処理

#pragma mark bitfeilds2macro(proc)
#if 0
	type:4;			///< 使用装置
	proc:4;			///< 処理種別
#endif

	uint8_t		use[2];			///< 利用金額 符号付12bit（10円単位）

#pragma mark bitfeilds2macro(use)
#if 0
	_:4;			///< 不明
	use:12;			///< 利用金額（支払いはマイナス）
#endif

	uint8_t		balance[2];		///< 残額
	uint8_t		unknown;		///< 不明
} nicepass_value_t;


#endif /* PASORIKIT_NICEPASS_H */
