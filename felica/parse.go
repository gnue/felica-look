package felica

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// カード情報
type CardInfo map[string]*SystemInfo

// システム情報
type SystemInfo struct {
	idm      string
	pmm      string
	services ServiceInfo
}

// サービス情報
type ServiceInfo map[string]([][]byte)

// SystemInfoメンバーの getter
func (sysinfo SystemInfo) IDm() string {
	return sysinfo.idm
}

func (sysinfo SystemInfo) PMm() string {
	return sysinfo.pmm
}

func (sysinfo SystemInfo) Services() ServiceInfo {
	return sysinfo.services
}

func (sysinfo SystemInfo) ServiceCodes() []string {
	codes := make([]string, 0, len(sysinfo.services))

	for svccode, _ := range sysinfo.services {
		codes = append(codes, svccode)
	}

	return codes
}

// FeliCaダンプファイルを読込む
func Read(path string) *CardInfo {
	cardinfo := CardInfo{}

	// IDmの正規表現
	re_idm := []*regexp.Regexp{
		regexp.MustCompile("(?i)IDm = *([0-9A-F]+)"),
		regexp.MustCompile("(?i)IDm :(( [0-9A-F]+)+)"),
	}

	// PMmの正規表現
	re_pmm := []*regexp.Regexp{
		regexp.MustCompile("(?i)PMm = *([0-9A-F]+)"),
		regexp.MustCompile("(?i)PMm :(( [0-9A-F]+)+)"),
	}

	// システムコードの正規表現
	re_syscode := []*regexp.Regexp{
		regexp.MustCompile("(?i)^# FELICA SYSTEM_CODE = *([0-9A-F]+)"),
		regexp.MustCompile("(?i)^# System code: ([0-9A-F]+)"),
	}

	// サービスコードの正規表現
	re_svccode := []*regexp.Regexp{
		regexp.MustCompile("(?i)^# [0-9A-F]+:[0-9A-F]+:([0-9A-F]+) #[0-9A-F]+"),
		regexp.MustCompile("(?i)# Serivce code = *([0-9A-F]+)"),
	}

	// データの正規表現
	re_data := []*regexp.Regexp{
		regexp.MustCompile("(?i)^ *[0-9A-F]+:[0-9A-F]+:([0-9A-F]+):[0-9A-F]+:([0-9A-F]{32})"),
		regexp.MustCompile("(?i)^ *([0-9A-F]+):[0-9A-F]+(( [0-9A-F]+){16})"),
	}

	file, err := os.Open(path)
	if err != nil {
		// エラー処理をする
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var svccode string
	orphan := &SystemInfo{idm: "", pmm: "", services: make(ServiceInfo)}
	currsys := orphan

	for scanner.Scan() {
		line := scanner.Text()

		// IDm
		for _, re := range re_idm {
			match := re.FindStringSubmatch(line)
			if match != nil {
				idm := match[1]
				currsys.idm = strings.Replace(idm, " ", "", -1)
			}
		}

		// PMm
		for _, re := range re_pmm {
			match := re.FindStringSubmatch(line)
			if match != nil {
				pmm := match[1]
				currsys.pmm = strings.Replace(pmm, " ", "", -1)
			}
		}

		// システムコード
		for _, re := range re_syscode {
			match := re.FindStringSubmatch(line)
			if match != nil {
				syscode := match[1]
				currsys = &SystemInfo{idm: "", pmm: "", services: make(ServiceInfo)}
				cardinfo[syscode] = currsys
			}
		}

		// サービスコード
		for _, re := range re_svccode {
			match := re.FindStringSubmatch(line)
			if match != nil {
				svccode = match[1]
				currsys.services[svccode] = [][]byte{}
			}
		}

		// データ
		for _, re := range re_data {
			match := re.FindStringSubmatch(line)
			if match != nil {
				data := match[2]
				data = strings.Replace(data, " ", "", -1)
				buf := hex2bin(data)
				currsys.services[svccode] = append(currsys.services[svccode], buf)
			}
		}
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
