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

// カードに対応するモジュールを探す
func find_module(cardinfo felica.CardInfo, modules []felica.Module) felica.Module {
	for _, m := range modules {
		if m.IsCard(cardinfo) {
			return m
		}
	}

	return nil
}

// コード・リストを文字列のリストにする
func codes_to_strings(codes []uint16) []string {
	var list []string

	for _, code := range codes {
		list = append(list, fmt.Sprintf("%04X", code))
	}

	return list
}

// カード情報を簡易出力する
func show_info(cardinfo felica.CardInfo) {
	for syscode, currsys := range cardinfo {
		fmt.Printf("SYSTEM CODE: %04X\n", syscode)
		fmt.Println("  IDm: ", currsys.IDm)
		fmt.Println("  PMm: ", currsys.PMm)
		fmt.Println("  SERVICE CODES: ", codes_to_strings(currsys.ServiceCodes))
	}
}

// カード情報をダンプ出力する
func dump_info(cardinfo felica.CardInfo) {
	for syscode, currsys := range cardinfo {
		fmt.Printf("SYSTEM CODE: %04X\n", syscode)
		fmt.Println("  IDm: ", currsys.IDm)
		fmt.Println("  PMm: ", currsys.PMm)

		for svccode, data := range currsys.Services {
			fmt.Printf("  SERVICE CODE: %04X\n", svccode)

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
				engine := m.Bind(cardinfo)

				fmt.Printf("%s:\n", engine.Name())
				engine.ShowInfo(&options)
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
