/*------------------------------------------------------------------------*/
/**
 * @file	suica.h
 * @brief   Suica
 *
 * @author  M.Nukui
 * @date	2008-04-02
 *
 * Copyright (C) 2008 M.Nukui All rights reserved.
 */


#ifndef	PASORIKIT_SUICA_H
#define	PASORIKIT_SUICA_H


#include <stdint.h>


#define FELICA_POLLING_SUICA			0x0003		///< Suicaシステムコード
#define FELICA_POLLING_IRUCA			0xDE80		///< IruCaシステムコード
#define FELICA_POLLING_SAPICA			0x865E		///< SAPICAシステムコード
#define FELICA_POLLING_PASPY			0x8592		///< PASPYシステムコード

#define FELICA_SC_SUICA_VALUE			0x090F		///< Suica利用履歴データ・サービスコード
#define FELICA_SC_SUICA_TICKETING		0x184B		///< Suica発券改札データ・サービスコード


#define suica_year(v)		(v->date[0] >> 1)									///< 年の取得
#define suica_month(v)		(((v->date[0] & 1) << 3) + (v->date[1] >> 5))		///< 月の取得
#define suica_day(v)		(v->date[1] & 0x1F)									///< 日の取得
#define suica_balance(v)	((v->balance[1] << 8) + v->balance[0])				///< 残額の取得

#define suica_areacode(region, lineCode)	((region == 0)?((lineCode < 0x80)?0:1):((region == 1)?2:-1))	///< エリアコード


/// Suica履歴データ
typedef struct {
	uint8_t		type;			///< 端末種
	uint8_t		proc;			///< 処理
	uint8_t		unknown1[2];	///< 不明
	uint8_t		date[2];		///< 年(7bit)/月(4bit)/日(5bit) 
	uint8_t		in_station[2];	///< 入場駅（線区コード、駅順コード）
	uint8_t		out_station[2];	///< 出場駅（線区コード、駅順コード）
	uint8_t		balance[2];		///< 残額（リトルエンディアン） 
	uint8_t		no[3];			///< 連番
	uint8_t		region;			///< リージョン
} suica_value_t;


/// Suica発券データ
typedef struct {
	uint8_t		in_station[2];			///< 入場駅（線区コード、駅順コード）
	uint8_t		out_station[2];			///< 出場駅（線区コード、駅順コード）
	uint8_t		expire_date[2];			///< 有効期限、年(7bit)/月(4bit)/日(5bit) 
	uint8_t		ticketing_time[2];		///< 発券時刻、時(5bit)/分(6bit)/秒(5bit) ２秒単位
	uint8_t		kind;					///< 種別
	uint8_t		price;					///< 金額（10円単位）
	uint8_t		examination_device[2];	///< 改札装置番号 
	uint8_t		examination_station[2];	///< 改札実施駅（線区コード、駅順コード）
	uint8_t		examination_time[2];	///< 改札実施時刻、時(5bit)/分(6bit)/秒(5bit) ２秒単位
} suica_ticketing_t;


#endif /* PASORIKIT_SUICA_H */
