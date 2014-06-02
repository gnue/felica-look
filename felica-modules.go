package main

import (
	"github.com/gnue/felica-look/edy"
	"github.com/gnue/felica-look/felica"
	"github.com/gnue/felica-look/rapica"
	"github.com/gnue/felica-look/suica"
)

var felica_modules = []felica.Module{
	&suica.Module,
	&edy.Module,
	&rapica.Module,
}
