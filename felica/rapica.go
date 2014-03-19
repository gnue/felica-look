package felica

/*
#include "rapica.h"
*/
import "C"

// RapiCa/鹿児島市交通局
type RapiCa struct {
}

// カード名
func (rapica *RapiCa) Name() string {
	return "RapiCa"
}

// システムコード
func (rapica *RapiCa) SystemCode() uint64 {
	return C.FELICA_POLLING_RAPICA
}

// カード情報を表示する
func (rapica *RapiCa) ShowInfo(cardinfo *CardInfo, extend bool) {
}
