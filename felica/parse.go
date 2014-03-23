package felica

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 正規表現とアクション
type re_action struct {
	_regexpes []*regexp.Regexp
	regexpes  []string
	action    func(match []string)
}

// FeliCaダンプファイルを読込む
func Read(path string) CardInfo {
	cardinfo := CardInfo{}

	file, err := os.Open(path)
	if err != nil {
		// エラー処理をする
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var svccode string
	orphan := empty_sysinfo()
	currsys := orphan

	actions := [](*re_action){
		// IDm
		{
			regexpes: []string{
				"(?i)IDm = *([0-9A-F]+)",
				"(?i)IDm :(( [0-9A-F]+)+)",
			},
			action: func(match []string) {
				currsys.IDm = strings.Replace(match[1], " ", "", -1)
			},
		},

		// PMm
		{
			regexpes: []string{
				"(?i)PMm = *([0-9A-F]+)",
				"(?i)PMm :(( [0-9A-F]+)+)",
			},
			action: func(match []string) {
				currsys.PMm = strings.Replace(match[1], " ", "", -1)
			},
		},

		// システムコード
		{
			regexpes: []string{
				"(?i)^# FELICA SYSTEM_CODE = *([0-9A-F]+)",
				"(?i)^# System code: ([0-9A-F]+)",
			},
			action: func(match []string) {
				syscode := match[1]
				currsys = empty_sysinfo()
				cardinfo[syscode] = currsys
			},
		},

		// サービスコード
		{
			regexpes: []string{
				"(?i)^# [0-9A-F]+:[0-9A-F]+:([0-9A-F]+) #[0-9A-F]+",
				"(?i)# Serivce code = *([0-9A-F]+)",
			},
			action: func(match []string) {
				svccode = match[1]
				currsys.ServiceCodes = append(currsys.ServiceCodes, svccode)
				currsys.Services[svccode] = [][]byte{}
			},
		},

		// felica_dump サービスコード
		{
			regexpes: []string{
				"(?i)^# ([0-9A-F]{4}):([0-9A-F]{4}) ",
			},
			action: func(match []string) {
				code, _ := strconv.ParseInt(match[1], 16, 0)
				attr, _ := strconv.ParseInt(match[2], 16, 0)
				svccode = fmt.Sprintf("%04X", code<<6+attr)
				currsys.ServiceCodes = append(currsys.ServiceCodes, svccode)
				currsys.Services[svccode] = [][]byte{}
			},
		},

		// データ
		{
			regexpes: []string{
				"(?i)^ *[0-9A-F]+:[0-9A-F]+:([0-9A-F]+):[0-9A-F]+:([0-9A-F]{32})",
				"(?i)^ *([0-9A-F]+):[0-9A-F]+(( [0-9A-F]+){16})",
			},
			action: func(match []string) {
				data := match[2]
				data = strings.Replace(data, " ", "", -1)
				buf, _ := hex.DecodeString(data)
				currsys.Services[svccode] = append(currsys.Services[svccode], buf)
			},
		},

		// felica_dump データ
		{
			regexpes: []string{
				"(?i)^  [0-9A-F]{4}:[0-9A-F]{4}:([0-9A-F]{32})",
			},
			action: func(match []string) {
				data := match[1]
				data = strings.Replace(data, " ", "", -1)
				buf, _ := hex.DecodeString(data)
				currsys.Services[svccode] = append(currsys.Services[svccode], buf)
			},
		},
	}

	// 正規表現のコンパイル
	re_action_compile(actions)

	for scanner.Scan() {
		line := scanner.Text()
		re_match_action(line, actions, true)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	if len(orphan.IDm) != 0 || len(orphan.PMm) != 0 {
		for _, currsys := range cardinfo {
			currsys.IDm = orphan.IDm
			currsys.PMm = orphan.PMm
		}
	}

	return cardinfo
}

// 空の SystemInfo を作成する
func empty_sysinfo() *SystemInfo {
	return &SystemInfo{ServiceCodes: []string{}, Services: make(ServiceInfo)}
}

// 正規表現をコンパイルする
func re_action_compile(actions [](*re_action)) {
	for _, a := range actions {
		a._regexpes = make([]*regexp.Regexp, len(a.regexpes))

		for i, s := range a.regexpes {
			a._regexpes[i] = regexp.MustCompile(s)
		}

	}
}

// 正規表現に一致したら対応するアクションを実行する
func re_match_action(text string, actions [](*re_action), is_break bool) {
	for _, a := range actions {
		for _, re := range a._regexpes {
			match := re.FindStringSubmatch(text)
			if match != nil {
				a.action(match)

				if is_break {
					// 残りは実行しない
					return
				}
			}
		}
	}
}
