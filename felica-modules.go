package main

import (
	"./edy"
	"./felica"
	"./rapica"
	"./suica"
)

var felica_modules = []felica.Module{
	&suica.Module,
	&edy.Module,
	&rapica.Module,
}
