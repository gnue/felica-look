package main

import (
	"./felica"
	"./rapica"
)

func felica_modules() []felica.Module {
	modules := []felica.Module{
		&rapica.Module,
	}

	return modules
}
