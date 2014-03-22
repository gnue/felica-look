package main

import (
	"./felica"
	"./rapica"
	"./suica"
)

func felica_modules() []felica.Module {
	modules := []felica.Module{
		&suica.Module,
		&rapica.Module,
	}

	return modules
}
