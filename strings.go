package gtools

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/message"
	"golang.org/x/text/transform"
)

const charset string = "abcdefghzkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ3456789" // 随机因子

func GenerateRandomStr(length int) string {
	retStr := ""
	csLen := len(charset)
	for i := 0; i < length; i++ {
		randNum := GenerateRandom(0, csLen)
		retStr += string(charset[randNum])
	}

	return retStr
}

// GenerateMobileCaptcha 手机验证在4-8位之间,简单实现
func GenerateMobileCaptcha(length int) string {
	if length < 4 || length > 8 {
		return ""
	}

	minStr := "1" + strings.Repeat("0", length-1)
	maxStr := "1" + strings.Repeat("0", length)

	min, _ := strconv.Atoi(minStr)
	max, _ := strconv.Atoi(maxStr)

	captcha := GenerateRandom(min, max)
	return strconv.Itoa(captcha)
}

// url base64编码
func URLBase64Encode(data []byte) string {
	base64encodeBytes := base64.URLEncoding.EncodeToString(data)
	base64encodeBytes = strings.Replace(base64encodeBytes, "=", ".", -1)
	return base64encodeBytes
}

// url base64解码
func URLBase64Decode(data string) ([]byte, error) {
	data = strings.Replace(data, ".", "=", -1)
	decodeBytes, err := base64.URLEncoding.DecodeString(data)

	return decodeBytes, err
}

// base64编码
func Base64Encode(data []byte) string {
	base64encodeBytes := base64.StdEncoding.EncodeToString(data)
	return base64encodeBytes
}

// base64解码
func Base64Decode(data string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	return decodeBytes, err
}

// url 编码
func UrlEncode(data string) string {
	encode := url.QueryEscape(data)
	return encode
}

// url 解码
func UrlDecode(data string) (string, error) {
	decodeUrl, err := url.QueryUnescape(data)
	return decodeUrl, err
}

// 字串截取
func SubString(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}

// 字串截取
func SubStringByPos(s string, sPos, ePos int) string {
	runes := []rune(s)
	if ePos > len(runes) {
		ePos = len(runes)
	}

	return string(runes[sPos:ePos])
}

func StrTrim(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)

	return str
}

// StrReplace 在 origin 中搜索 search 组,替换成 replace
func StrReplace(origin string, search []string, replace string) (s string) {
	s = origin
	for _, find := range search {
		s = strings.Replace(s, find, replace, -1)
	}

	return
}

func TrimRealName(name string) string {
	// 将名字里面的标点符号替换成1个空格
	name = ReplaceInvalidRealName(name)
	reg := regexp.MustCompile(`[\pP]+?`)
	name = reg.ReplaceAllString(name, " ")

	// 将连续的空白替换成一个空格
	reg = regexp.MustCompile(`\s{2,}`)
	name = reg.ReplaceAllString(name, " ")

	name = strings.TrimSpace(name)

	return name
}

func ReplaceInvalidRealName(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z\\s]+")
	ret := reg.ReplaceAllString(name, "")
	return ret
}

// IsNumber 判断是否都是数字
func IsNumber(str string) (valid bool) {
	var exp, _ = regexp.Compile(`^[0-9]+$`)
	if exp.MatchString(str) {
		return true
	}

	return false
}

func ContainNumber(str string) (valid bool) {
	var exp, _ = regexp.Compile(`[\d]`)
	if exp.MatchString(str) {
		return true
	}

	return false
}

func Str2Int64(str string) (int64, error) {
	number, err := strconv.ParseInt(str, 10, 64)
	return number, err
}

func Int642Str(number int64) string {
	return strconv.FormatInt(number, 10)
}

func Str2Int(str string) (int, error) {
	number, err := strconv.ParseInt(str, 10, 0)
	return int(number), err
}

func Int2Str(number int) string {
	return strconv.FormatInt(int64(number), 10)
}

func Float2Str(f float32) string {
	return Float642Str(float64(f))
}

func Float642Str(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Str2Float64(s string) (f float64, err error) {
	f, err = strconv.ParseFloat(s, 64)

	return
}

func Str2Float(s string) (f float32, err error) {
	f64, err := Str2Float64(s)
	f = float32(f64)

	return
}

func JsonEncode(d interface{}) (jsonStr string, err error) {
	bson, err := json.Marshal(d)
	jsonStr = string(bson)

	return
}

func Unicode(rs string) string {
	jsonStr := ""
	for _, r := range rs {
		rInt := int(r)
		if rInt < 128 {
			jsonStr += string(r)
		} else {
			jsonStr += "\\u" + strconv.FormatInt(int64(rInt), 16)
		}
	}

	return jsonStr
}

func Escape(html string) string {
	return template.HTMLEscapeString(html)
}

func AddSlashes(str string) string {
	str = strings.Replace(str, `\`, `\\`, -1)
	str = strings.Replace(str, "'", `\'`, -1)
	str = strings.Replace(str, `"`, `\"`, -1)

	return str
}

func StripSlashes(str string) string {
	str = strings.Replace(str, `\'`, `'`, -1)
	str = strings.Replace(str, `\"`, `"`, -1)
	str = strings.Replace(str, `\\`, `\`, -1)

	return str
}

func RawUrlEncode(s string) (r string) {
	r = UrlEncode(s)
	r = strings.Replace(r, "+", "%20", -1)
	return
}

// 直接json.Marshal ，  会把 < > & 转成 unicode 编码
// JSONMarshal 解决直接json.Marshal 后单引号，双引号，< > & 符号的问题
func JSONMarshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)

	return buf.Bytes(), err
}

// ParseTableName 从SQL语句中解析主表名
func ParseTableName(sql string) (name string, err error) {
	re := regexp.MustCompile(`(?i).+FROM\s+(\S+)`)
	allMatch := re.FindAllStringSubmatch(sql, -1)

	if len(allMatch) == 0 {
		err = fmt.Errorf("parse sql has error")
		return
	}

	tableName := strings.Replace(allMatch[0][1], "`", "", -1)
	names := strings.Split(tableName, ".")
	name = names[len(names)-1]

	return
}

func RealNameMask(realName string) (maskStr string) {
	words := ([]rune)(realName)
	wordsLen := len(words)

	if wordsLen <= 1 {
		maskStr = realName
	} else if wordsLen == 2 {
		maskStr = fmt.Sprintf("%s*", string(words[0:1]))
	} else {
		var star string
		for i := 0; i < wordsLen-2; i++ {
			star += "*"
		}
		maskStr = fmt.Sprintf("%s%s%s", string(words[0:1]), star, string(words[wordsLen-1:]))
	}

	return
}

func ParseTargetList(str string) []string {
	list := make([]string, 0)
	if str == "" {
		return list
	}

	listStr := strings.Split(str, "\n")

	c := ","
	if strings.Contains(listStr[0], "\r") {
		c = "\r"
	} else if strings.Contains(listStr[0], ",") {
		c = ","
	}

	if len(listStr) == 1 {
		vec := strings.Split(listStr[0], c)

		list = append(list, strings.Trim(vec[0], " "))

		return list
	}

	for _, v := range listStr {
		vec := strings.Split(v, c)
		if len(vec) < 1 {
			continue
		}

		list = append(list, strings.Trim(vec[0], " "))
	}

	return list
}

func StrTenThousand2MoneyMul100(money string) int64 {
	price, _ := DecimalMoneyMul100(money)
	price *= 10000

	return price
}

// DecimalMoneyMul100 将带小数点的金额转换成数据库中的乘以100后的整数
func DecimalMoneyMul100(moneyStr string) (money int64, err error) {
	moneyBig, err := decimal.NewFromString(moneyStr)
	if err != nil {
		log.Println(fmt.Sprintf("[DecimalMoneyMul100] convert to math big exception, str: %s, err: %v", moneyStr, err))
		return
	}

	moneyMul := moneyBig.Mul(decimal.NewFromInt(100))
	money = moneyMul.IntPart()

	return
}

func SecretKeyMask(secretKey string) (str string) {
	keyLen := len(secretKey)
	if keyLen > 9 {
		str = fmt.Sprintf(`%s***%s`, SubString(secretKey, 0, 3), SubString(secretKey, keyLen-3, 3))
	} else {
		str = secretKey
	}

	return
}

/**
  13:23
*/
func ConvertTime2Secs(str string) int {
	arr := strings.Split(str, ":")
	hours, _ := Str2Int(arr[0])
	mins, _ := Str2Int(arr[1])
	return hours*60*60 + mins*60
}

// Snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// Camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data[:])
}

func SimpleNumStr2Num(s string) int {
	var ret int
	var base float64 = 1
	if strings.Contains(s, "w") || strings.Contains(s, "W") {
		base = 10000
		s = StrReplace(s, []string{"w", "W"}, "")
	}

	num, err := Str2Float64(s)
	if err != nil {
		log.Println(fmt.Sprintf("[SimpleNumStr2Num] str to num err: %v", err))
	}

	ret = int(num * base)

	return ret
}

func ExtractUrls(s string) ([]string, error) {
	var urlBox []string
	var err error

	findUrl := regexp.MustCompile(`(http\S*)`)
	allMatch := findUrl.FindAllStringSubmatch(s, -1)
	//log.Printf("[ExtractUrl] allMatch: %#v\n", allMatch)

	if len(allMatch) > 0 {
		for _, fu := range allMatch {
			urlBox = append(urlBox, fu[0])
		}
	} else {
		err = fmt.Errorf(`can not extract url from input: %s`, s)
	}

	return urlBox, err
}

func Gbk2Utf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}

	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}

	return d, nil
}

func MobileMask(mobile string) string {
	if len(mobile) != 11 {
		log.Println(fmt.Sprintf("input wrong mobile: %s", mobile))
		return ""
	}

	return fmt.Sprintf(`%s%s%s`, mobile[0:3], strings.Repeat(`*`, 5), mobile[8:])
}

func RegRemoveScript(in string) string {
	reg, _ := regexp.Compile(`\<(?i)script[\S\s]+?\</(?i)script\>`)
	in = reg.ReplaceAllString(in, "")
	return in
}

func NicknameMask(nickname string) string {
	var s = nickname

	rStr := []rune(nickname)
	rLen := len(rStr)

	if rLen > 0 && rLen < 2 {
		s = fmt.Sprintf(`%s*`, string(rStr[0:1]))
	} else if rLen > 2 && rLen <= 4 {
		s = fmt.Sprintf(`%s%s%s`, string(rStr[0:1]), strings.Repeat("*", rLen-2), string(rStr[rLen-1:]))
	} else if rLen > 4 {
		s = fmt.Sprintf(`%s%s%s`, string(rStr[0:1]), strings.Repeat("*", rLen-3), string(rStr[rLen-2:]))
	}

	return s
}

func TrimTags(s string) string {
	re := regexp.MustCompile(`#(\S+)`)
	out := strings.TrimSpace(re.ReplaceAllString(s, ""))
	return out
}

func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

func MoneyDisplay(money int64) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

func HumanMoney(money int64) string {
	str := MoneyDisplay(money)
	length := len(str)
	if length < 4 {
		return str
	}

	arr := strings.Split(str, ".") // 用小数点符号分割字符串,为数组接收
	length1 := len(arr[0])
	if length1 < 4 {
		return str
	}
	count := (length1 - 1) / 3

	for i := 0; i < count; i++ {
		arr[0] = arr[0][:length1-(i+1)*3] + "," + arr[0][length1-(i+1)*3:]
	}

	return strings.Join(arr, ".") // 将一系列字符串连接为一个字符串，之间用sep来分隔。
}

// NumberFormat 输出格式化数字，千分位以逗号分割
func NumberFormat(number interface{}) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprint(number)
}
