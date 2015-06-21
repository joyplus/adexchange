package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//password hash function
func Pwdhash(str string) string {
	return Strtomd5(str)
}

func StringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

//得到当前时间
func GetCurrentTime() string {

	return time.Now().Format("20060102150405")
}

//生成OTP
func GenerateOTP() string {
	nums := generateRandomNumber(100000, 999999, 1)
	return fmt.Sprintf("%d", nums[0])
}

//生成OTP Sequence Number
func GenerateSequenceNumberForOTP(otp string) string {
	return GetMd5String(otp)
}

//生成Security token
func GenerateSecurityToken(mobileNumber string) string {
	return GetMd5String(mobileNumber + GetCurrentTime())
}

//生成count个[start,end)结束的不重复的随机数
func generateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

//生成订单号
//REQ: 询价单
//RES: 报价单
//TRX: 订单
func GenerateOrderNumber(purchaseType string) string {
	nums := generateRandomNumber(10000, 99999, 1)
	return purchaseType + GetCurrentTime() + fmt.Sprintf("%d", nums[0])
}

func GenerateBid(prefix string) string {
	nums := generateRandomNumber(1, 10000, 1)
	return prefix + GetCurrentTime() + fmt.Sprintf("%d", nums[0])
}

func ConvertStrToInt(s string) int {

	i, err := strconv.Atoi(s)
	if err != nil {
		i = 0
	}

	return i
}

func ConvertIntToString(i int) string {

	return strconv.Itoa(i)
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DivisionInt(first int, second int) float32 {
	if second == 0 || first == 0 {
		return 0.0
	} else {
		f := float32(first) / float32(second)
		return f
	}
}

func EscapeCtrl(ctrl []byte) (esc []byte) {
	u := []byte(`\u0000`)
	for i, ch := range ctrl {
		if ch <= 31 {
			if esc == nil {
				esc = append(make([]byte, 0, len(ctrl)+len(u)), ctrl[:i]...)
			}
			esc = append(esc, u...)
			hex.Encode(esc[len(esc)-2:], ctrl[i:i+1])
			continue
		}
		if esc != nil {
			esc = append(esc, ch)
		}
	}
	if esc == nil {
		return ctrl
	}
	return esc
}
