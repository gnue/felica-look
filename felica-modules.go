package main

import (
	"./edy"
	"./felica"
	"./rapica"
	"./suica"
)

func felica_modules() []felica.Module {
	modules := []felica.Module{
		&suica.Module,
		&edy.Module,
		&rapica.Module,
	}

	return modules
}
