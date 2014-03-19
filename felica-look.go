package main

import (
	"./felica"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
)

// コマンドの使い方
func usage() {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", path.Base(cmd))
	flag.PrintDefaults()
	os.Exit(0)
}

func find_module(cardinfo *felica.CardInfo, modules []felica.Module) felica.Module {
	for syscode, _ := range *cardinfo {
		code, _ := strconv.ParseUint(syscode, 16, 0)

		for _, m := range modules {
			if m.SystemCode() == code {
				return m
			}
		}
	}

	return nil
}

// カード情報を簡易出力する
func show_info(cardinfo *felica.CardInfo) {
	for syscode, currsys := range *cardinfo {
		fmt.Println("SYSTEM CODE: ", syscode)
		fmt.Println("  IDm: ", currsys.IDm())
		fmt.Println("  PMm: ", currsys.PMm())
		fmt.Println("  SERVICE CODES: ", currsys.ServiceCodes())
	}
}

// カード情報をダンプ出力する
func dump_info(cardinfo *felica.CardInfo) {
	for syscode, currsys := range *cardinfo {
		fmt.Println("SYSTEM CODE: ", syscode)
		fmt.Println("  IDm: ", currsys.IDm())
		fmt.Println("  PMm: ", currsys.PMm())

		for svccode, data := range currsys.Services() {
			fmt.Println("  SERVICE CODE: ", svccode)

			for _, v := range data {
				fmt.Printf("      %X\n", v)
			}
		}
	}
}

func main() {
	dump := flag.Bool("d", false, "dump")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help || len(flag.Args()) == 0 {
		usage()
	}

	modules := []felica.Module{
	}

	for _, v := range flag.Args() {
		cardinfo := felica.Read(v)

		if *dump {
			dump_info(cardinfo)
		} else {
			m := find_module(cardinfo, modules)
			if m != nil {
				m.ShowInfo(cardinfo)
			} else {
				show_info(cardinfo)
			}
		}

	}
}
