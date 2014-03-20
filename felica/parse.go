package felica

import (
	"bufio"
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
func Read(path string) *CardInfo {
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
				currsys.idm = strings.Replace(match[1], " ", "", -1)
			},
		},

		// PMm
		{
			regexpes: []string{
				"(?i)PMm = *([0-9A-F]+)",
				"(?i)PMm :(( [0-9A-F]+)+)",
			},
			action: func(match []string) {
				currsys.pmm = strings.Replace(match[1], " ", "", -1)
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
				currsys.svccodes = append(currsys.svccodes, svccode)
				currsys.services[svccode] = [][]byte{}
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
				buf := hex2bin(data)
				currsys.services[svccode] = append(currsys.services[svccode], buf)
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

	if len(orphan.idm) != 0 || len(orphan.pmm) != 0 {
		for _, currsys := range cardinfo {
			currsys.idm = orphan.idm
			currsys.pmm = orphan.pmm
		}
	}

	return &cardinfo
}

// 空の SystemInfo を作成する
func empty_sysinfo() *SystemInfo {
	return &SystemInfo{svccodes: []string{}, services: make(ServiceInfo)}
}

// 16進文字列をバイナリに変換する
func hex2bin(hex string) []byte {
	buf := make([]byte, len(hex)/2)

	p := 0
	for i := 0; i < len(hex); i += 2 {
		b, _ := strconv.ParseUint(hex[i:i+2], 16, 8)
		buf[p] = byte(b)
		p += 1
	}

	return buf
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
