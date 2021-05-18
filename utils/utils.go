package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
)

func GetIntHour(startTime, endTime int64) (str []string, err error) {
	var timeArr []int64
	if endTime <= startTime {
		return nil, nil
	}

	var i int64 = 0
	for ; i < (endTime - startTime); i++ {
		if (startTime+i)%3600 == 0 {
			timeArr = append(timeArr, startTime+i)
			break
		}
	}
	if len(timeArr) < 1 {
		return nil, nil
	}

	for j := 1; startTime+int64(3600*j) < endTime; j++ {
		timeArr = append(timeArr, timeArr[0]+int64(3600*j))
	}
	timeArr = append(timeArr, endTime)
	// [0,1 1,2 2,3 3,4]
	var tr string
	for i := 0; i < len(timeArr); i++ {
		if len(timeArr) == i+1 {
			continue
		}
		if timeArr[i] <= timeArr[i+1] {
			tr = strconv.FormatInt(timeArr[i], 10)
			st := tr + "," + strconv.FormatInt(timeArr[i+1], 10)
			str = append(str, st)
		}
	}

	return
}

func GetIntDay(startTime, endTime int64) (str []string, err error) {
	var timeArr []int64
	if endTime <= startTime {
		return nil, errors.New("endTime <= startTime")
	}

	var i int64 = 0
	for ; i < (endTime - startTime); i++ {
		if (startTime+i)%3600 == 0 {
			timeArr = append(timeArr, startTime+i)
			break
		}
	}
	if len(timeArr) < 1 {
		return nil, errors.New("Time No IntHour")
	}

	for j := 1; startTime+int64(3600*j) < endTime; j++ {
		if j%24 == 0 {
			timeArr = append(timeArr, timeArr[0]+int64(3600*j))
		}
	}
	timeArr = append(timeArr, endTime)
	// [0,1 1,2 2,3 3,4]
	var tr string
	for i := 0; i < len(timeArr); i++ {
		if len(timeArr) == i+1 {
			continue
		}
		if timeArr[i] < timeArr[i+1] {
			tr = strconv.FormatInt(timeArr[i], 10)
			st := tr + "," + strconv.FormatInt(timeArr[i+1], 10)
			str = append(str, st)
		}
	}

	return
}

func TimeStringK(sTime int64) string {
	return time.Unix(sTime, 0).Format("2006-01-02 00:00:00")
}

func TableName(date time.Time) string {
	var tableName string
	newDateOne := int(date.Month())
	if int(date.Month()) < 10 {
		tableName = fmt.Sprintf("%v0%v", date.Year(), newDateOne)
	} else {
		tableName = fmt.Sprintf("%v%v", date.Year(), newDateOne)
	}

	return tableName
}

func Min(vals ...decimal.Decimal) decimal.Decimal {
	var min decimal.Decimal
	for _, val := range vals {
		if min.IsZero() || val.Cmp(min) <= 0 {
			min = val
		}
	}
	return min
}

func TimeInit() (startTime, endTime int64) {
	indexTime, err := time.ParseInLocation("2006-01-02 00:00:00", time.Now().Format("2006-01-02 00:00:00"), time.Local)
	if err != nil {
		return
	}

	startTime = indexTime.Unix()
	endTime = time.Now().Unix()

	return startTime, endTime
}

func Mining_GetIntDay(startTime, endTime int64, countTime int64) (str []string, err error) {
	var timeArr []int64
	if endTime <= startTime {
		return nil, errors.New("endTime <= startTime")
	}

	var i int64 = 0
	for ; i < (endTime - startTime); i++ {
		if (startTime+i)%3600 == 0 {
			timeArr = append(timeArr, startTime+i)
			break
		}
	}
	if len(timeArr) < 1 {
		return nil, errors.New("Time No IntHour")
	}

	var timeCount int
	switch countTime {
	case int64(24):
		timeCount = 24
	case int64(168):
		timeCount = 168
	case int64(720):
		timeCount = 720
	default:
		timeCount = 24
	}

	for j := 1; startTime+int64(3600*j) < endTime; j++ {
		if j%timeCount == 0 {
			timeArr = append(timeArr, timeArr[0]+int64(3600*j))
		}
	}
	timeArr = append(timeArr, endTime)
	// [0,1 1,2 2,3 3,4]
	var tr string
	for i := 0; i < len(timeArr); i++ {
		if len(timeArr) == i+1 {
			continue
		}
		if timeArr[i] < timeArr[i+1] {
			tr = strconv.FormatInt(timeArr[i], 10)
			st := tr + "," + strconv.FormatInt(timeArr[i+1], 10)
			str = append(str, st)
		}
	}

	return
}

func ExecCmd(cmdStr string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func QueryTime(startTime int64, endTime int64) (timeStr []string, err error) {
	var timeArr []int64
	if endTime <= startTime {
		return nil, errors.New("endTime <= startTime")
	}

	var i int64 = 0
	for ; i < (endTime - startTime); i++ {
		if (startTime+i)%3600 == 0 {
			timeArr = append(timeArr, startTime+i)
			break
		}
	}
	if len(timeArr) < 1 {
		return nil, errors.New("Time No IntHour")
	}

	for j := 1; startTime+int64(3600*j) < endTime; j++ {
		if j%24 == 0 {
			timeArr = append(timeArr, timeArr[0]+int64(3600*j))
		}
	}
	timeArr = append(timeArr, endTime)
	// [0,1 1,2 2,3 3,4]
	var tr string
	for i := 0; i < len(timeArr); i++ {
		if len(timeArr) == i+1 {
			continue
		}
		if timeArr[i] < timeArr[i+1] {
			tr = strconv.FormatInt(timeArr[i], 10)
			st := tr + "," + strconv.FormatInt(timeArr[i+1], 10)
			timeStr = append(timeStr, st)
		}
	}

	return timeStr, err
}

func DateTimeF(str string) (time.Time, error) {
	identifyingDate, err := time.ParseInLocation("2006-01-02", str, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return identifyingDate, nil
}

func DataTable() string {
	var tableName string
	d0 := TimeStr("2020-12-01 00:00:00").Day()
	d1 := time.Now().AddDate(0, 0, -time.Now().Day()+1).Day()
	if d0 == d1 {
		tableName = time.Now().AddDate(0, 0, -5).Format("200601")
	} else {
		tableName = time.Now().Format("200601")
	}

	return tableName
}

func TimeStr(sTime string) time.Time {
	indexTime, err := time.ParseInLocation("2006-01-02 00:00:00", sTime, time.Local)
	if err != nil {
		return time.Time{}
	}
	return indexTime
}

func ParseInLocation() {
	indexTime111, err22 := time.ParseInLocation("2006-01-02 00:00:00", time.Now().Format("2006-01-02 00:00:00"), time.Local)
	if err22 != nil {
		fmt.Println("\n", err22)
	}

	fmt.Println("\n", time.Now().Unix())
	fmt.Println("\n", indexTime111.Unix())
}

func TimeIn64(sTime string) int64 {
	indexTime, err := time.ParseInLocation("2006-01-02 00:00:00", sTime, time.Local)
	if err != nil {
		return 0
	}
	return indexTime.Unix()
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func CheckAgency(date, date1 string) bool {
	indexTime1, err := time.ParseInLocation("2006-01-02 00:00:00", fmt.Sprintf("%v 00:00:00", date), time.Local)
	if err != nil {
		return false
	}

	indexTime2, err := time.ParseInLocation("2006-01-02 00:00:00", fmt.Sprintf("%v 00:00:00", date1), time.Local)
	if err != nil {
		return false
	}

	ite := decimal.RequireFromString(FloatToString(indexTime2.Sub(indexTime1).Hours()))
	idy := decimal.NewFromInt(int64(24 * 90))

	if ite.LessThan(idy) {
		return false
	} else {
		return true
	}
}

//生成随机字符串 数字+大写+小写
func GetRandomS(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomlowS(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func StrToMD5(str string) (result string) {
	md5Ctx1 := md5.New()
	md5Ctx1.Write([]byte(str))
	result = fmt.Sprintf("%x", md5Ctx1.Sum(nil))
	return
}
