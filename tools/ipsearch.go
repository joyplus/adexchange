package tools

import (
	zh "code.google.com/p/go.text/encoding/simplifiedchinese"
	"code.google.com/p/go.text/transform"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var aryDirectCities []string
var arySpecialProvince []string

type Ip2LocationReq struct {
	ip string //查询的ip
}
type Ip2LocationResp struct {
	ok      bool
	ip      string
	country string
	area    string
}

const queryLength = 2

var queryPool chan Ip2LocationReq
var recodePool chan Ip2LocationResp
var queryMutex sync.RWMutex

func init() {
	aryDirectCities = []string{"北京市", "天津市", "上海市", "重庆市", "香港", "澳门"}
	arySpecialProvince = []string{"内蒙古", "广西", "西藏", "宁夏", "新疆"}
}

func startQueryService(dbfile string) {

	file, err := os.Open(dbfile)

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 32)

	// header
	file.ReadAt(buf[0:8], 0)

	indexStart := int64(binary.LittleEndian.Uint32(buf[0:4]))
	indexEnd := int64(binary.LittleEndian.Uint32(buf[4:8]))

	log.Printf("Index range: %d - %d", indexStart, indexEnd)

	for {

		req, eor := <-queryPool

		if !eor {
			log.Fatal("empty query.")
		}

		//log.Printf("%#v",req)

		ip := net.ParseIP(req.ip)
		ip4 := make([]byte, 4)
		ip4 = ip.To4() // := &net.IPAddr{IP:ip}
		//log.Printf("IP4: %#v",ip4)

		//log.Printf("%#v",req.ip)

		//二分法查找
		maxLoop := int64(32)
		head := indexStart //+ 8
		tail := indexEnd   //+ 8

		//是否找到
		got := false
		rpos := int64(0)

		for ; maxLoop >= 0 && len(ip4) == 4; maxLoop-- {
			idxNum := (tail - head) / 7
			pos := head + int64(idxNum/2)*7

			//pos += maxLoop*7

			file.ReadAt(buf[0:7], pos)

			// startIP
			_ip := binary.LittleEndian.Uint32(buf[0:4])

			//log.Printf("%d - INs:%d POS:%d %#v %d.%d.%d.%d",maxLoop,idxNum,pos,buf[0:7],_ip&0xff,_ip>>8&0xff,_ip>>16&0xff,_ip>>24&0xff)

			//记录位置
			_buf := append(buf[4:7], 0x0) // 3byte + 1byte(0x00)
			rpos = int64(binary.LittleEndian.Uint32(_buf))
			//log.Printf("%d %#v",rpos,_buf)

			file.ReadAt(buf[0:4], rpos)

			_ip2 := binary.LittleEndian.Uint32(buf[0:4])

			//log.Printf("IP_END:%#v %d.%d.%d.%d",buf[0:4],_ip2&0xff,_ip2>>8&0xff,_ip2>>16&0xff,_ip2>>24&0xff)

			//查询的ip被转成大端了
			_ipq := binary.BigEndian.Uint32(ip4)

			if _ipq > _ip2 {
				head = pos
				continue
			}

			if _ipq < _ip {
				tail = pos
				continue
			}

			// got

			got = true

			break

		}

		loc := Ip2LocationResp{
			ok:      false,
			ip:      req.ip,
			country: "",
			area:    "",
		}
		if got {
			_loc := getIpLocation(file, rpos)

			var tr *transform.Reader
			tr = transform.NewReader(strings.NewReader(_loc.country), zh.GBK.NewDecoder())

			if s, err := ioutil.ReadAll(tr); err == nil {
				loc.country = string(s)
			}

			tr = transform.NewReader(strings.NewReader(_loc.area), zh.GBK.NewDecoder())

			if s, err := ioutil.ReadAll(tr); err == nil {
				loc.area = string(s)
			}

			loc.ok = _loc.ok

		}

		recodePool <- loc

	}

}

func getIpLocation(file *os.File, offset int64) (loc Ip2LocationResp) {

	buf := make([]byte, 1024)

	file.ReadAt(buf[0:1], offset+4)

	mod := buf[0]

	//log.Printf("C1 FLAG: %#v", mod)

	countryOffset := int64(0)
	areaOffset := int64(0)

	if 0x01 == mod {
		countryOffset = _readLong3(file, offset+5)
		//log.Printf("Redirect to: %#v",countryOffset);

		file.ReadAt(buf[0:1], countryOffset)

		mod2 := buf[0]

		//log.Printf("C2 FLAG: %#v", mod2)

		if 0x02 == mod2 {
			loc.country = _readString(file, _readLong3(file, countryOffset+1))
			areaOffset = countryOffset + 4
		} else {
			loc.country = _readString(file, countryOffset)
			areaOffset = countryOffset + int64(len(loc.country)) + 1
			// areaOffset = file.Seek(0,1) // got the pos
			// log.Printf("cPOS: %#v aPOS: %#v err: %#v",countryOffset,areaOffset,err3)

		}

		loc.area = _readArea(file, areaOffset)

	} else if 0x02 == mod {
		loc.country = _readString(file, _readLong3(file, offset+5))
		loc.area = _readArea(file, offset+8)
	} else {
		loc.country = _readString(file, offset+4)
		areaOffset = offset + 4 + int64(len(loc.country)) + 1
		//areaOffset,_ = file.Seek(0,1)

		loc.area = _readArea(file, areaOffset)
	}

	loc.ok = true

	return
}
func _readArea(file *os.File, offset int64) string {
	buf := make([]byte, 4)

	file.ReadAt(buf[0:1], offset)

	mod := buf[0]

	//log.Printf("A FLAG: %#v", mod)

	if 0x01 == mod || 0x02 == mod {
		areaOffset := _readLong3(file, offset+1)
		if areaOffset == 0 {
			return ""
		} else {
			return _readString(file, areaOffset)
		}
	}
	return _readString(file, offset)
}

func _readLong3(file *os.File, offset int64) int64 {
	buf := make([]byte, 4)
	file.ReadAt(buf, offset)
	buf[3] = 0x00

	return int64(binary.LittleEndian.Uint32(buf))
}

func _readString(file *os.File, offset int64) string {

	buf := make([]byte, 1024)
	got := int64(0)

	for ; got < 1024; got++ {
		file.ReadAt(buf[got:got+1], offset+got)

		if buf[got] == 0x00 {
			break
		}
	}

	return string(buf[0:got])
}

func QueryIP(ipStr string) (string, string) {

	if len(ipStr) == 0 {
		return "", ""
	}
	//todo local test mode
	if strings.EqualFold("[", ipStr) {
		return "上海市", "上海市"
	}

	queryMutex.Lock()
	queryPool <- Ip2LocationReq{
		ip: ipStr,
	}

	record := <-recodePool
	queryMutex.Unlock()

	province, city := getLocation(record)
	return province, city

}

func getLocation(record Ip2LocationResp) (province string, city string) {

	strRegion := record.country

	for _, directCity := range aryDirectCities {
		if strings.Contains(strRegion, directCity) {
			return directCity, directCity
		}
	}

	for _, province := range arySpecialProvince {
		if strings.Contains(strRegion, province) {

			rs := []rune(record.country)
			rsProvince := []rune(province)
			return string(rs[0:len(rsProvince)]), string(rs[len(rsProvince):])
		}
	}

	if strings.Contains(strRegion, "省") {

		result := strings.Split(strRegion, "省")

		return result[0] + "省", result[1]
	}

	return "", ""

}

func Init(dbfile string) {

	queryPool = make(chan Ip2LocationReq, queryLength)
	recodePool = make(chan Ip2LocationResp, queryLength)

	//启动查询进程
	go startQueryService(dbfile)

}
