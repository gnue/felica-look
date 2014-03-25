package main

import (
	"./felica"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// コマンドの使い方
func usage() {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, "usage: %s [options] [file...]\n", filepath.Base(cmd))
	flag.PrintDefaults()
	os.Exit(0)
}

func find_module(cardinfo felica.CardInfo, modules []felica.Module) felica.Module {
	for syscode, _ := range cardinfo {
		for _, m := range modules {
			if m.SystemCode() == syscode {
				return m
			}
		}
	}

	return nil
}

// カード情報を簡易出力する
func show_info(cardinfo felica.CardInfo) {
	for syscode, currsys := range cardinfo {
		fmt.Printf("SYSTEM CODE: %04X\n", syscode)
		fmt.Println("  IDm: ", currsys.IDm)
		fmt.Println("  PMm: ", currsys.PMm)
		fmt.Println("  SERVICE CODES: ", currsys.ServiceCodes)
	}
}

// カード情報をダンプ出力する
func dump_info(cardinfo felica.CardInfo) {
	for syscode, currsys := range cardinfo {
		fmt.Printf("SYSTEM CODE: %04X\n", syscode)
		fmt.Println("  IDm: ", currsys.IDm)
		fmt.Println("  PMm: ", currsys.PMm)

		for svccode, data := range currsys.Services {
			fmt.Println("  SERVICE CODE: ", svccode)

			for _, v := range data {
				fmt.Printf("      %X\n", v)
			}
		}
	}
}

func main() {
	extend := flag.Bool("e", false, "extend information")
	hex := flag.Bool("x", false, "with hex dump")
	dump := flag.Bool("d", false, "dump")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help {
		usage()
	}

	options := felica.Options{Extend: *extend, Hex: *hex}
	modules := felica_modules()

	show := func(path string) {
		cardinfo := felica.Read(path)

		if *dump {
			dump_info(cardinfo)
		} else {
			m := find_module(cardinfo, modules)
			if m != nil {
				fmt.Printf("%s:\n", m.Name())
				m.ShowInfo(cardinfo, &options)

				// モジュールを初期状態にする
				modules = felica_modules()
			} else {
				show_info(cardinfo)
			}
		}
	}

	if len(flag.Args()) == 0 {
		show("")
	} else {
		for _, v := range flag.Args() {
			show(v)
		}
	}
}
