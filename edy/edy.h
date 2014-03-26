/*------------------------------------------------------------------------*/
/**
 * @file	edy.h
 * @brief   Edy
 *
 * @author  M.Nukui
 * @date	2008-04-02
 *
 * Copyright (C) 2008 M.Nukui All rights reserved.
 */


#ifndef	PASORIKIT_EDY_H
#define	PASORIKIT_EDY_H


#include <stdint.h>


#define FELICA_POLLING_EDY			0xFE00		///< Edyシステムコード

#define FELICA_SC_EDY_INFO			0x110B		///< Edyカード情報・サービスコード
#define FELICA_SC_EDY_LAST			0x1317		///< Edy残額情報・サービスコード
#define FELICA_SC_EDY_VALUE			0x170F		///< Edy利用履歴データ・サービスコード


#define edy_use(v)		((v->use[0] << 8) + v->use[1])									///< 出入金の取得
#define edy_rest(v)		((v->rest[0] << 8) + v->rest[1])								///< 残額の取得
#define edy_days(v)		(((v->date[0] << 8) + v->date[1]) >> 1)							///< 累積日数（2000年から）の取得
#define edy_sec(v)		(((v->date[1] & 1) << 16) + (v->date[2] << 8) + v->date[3])		///< 秒の取得


/// Edyカード情報（addr=0）
typedef struct {
	uint8_t		unknown1[2];	///< 不明
	uint8_t		edyno[8];		///< Edy番号
	uint8_t		unknown2[6];	///< 不明
} edy_info0_t;


/// Edy残額情報（最終利用状況）
typedef struct {
	uint8_t		rest[4];		///< 残額(LE)
	uint8_t		use[4];			///< 直近使用金額(LE) チャージのときは更新されない場合がある
	uint8_t		unknown[6];		///< 不明
	uint8_t		no[2];			///< 取引通番(LE)
} edy_last_t;


/// Edy履歴データ
typedef struct {
	uint8_t		type;			///< タイプ
	uint8_t		no[3];			///< 連番
	uint8_t		date[4];		///< 31-17bit（2000年からの通算日数） 16-0bit（秒）
	uint8_t		unknown1[2];	///< 不明
	uint8_t		use[2];			///< 入金／出金
	uint8_t		unknown2[2];	///< 不明
	uint8_t		rest[2];		///< 残額
} edy_value_t;



#endif /* PASORIKIT_EDY_H */
