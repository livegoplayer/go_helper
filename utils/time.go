package utils

import (
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	Y         = "2006" // 年
	Y_CN      = Y + "年"
	M         = "01" // 月
	M_CN      = M + "月"
	D         = "02" // 日
	D_CN      = D + "日"
	Hour      = "15" // 时
	Hour_CN   = Hour + "时"
	Min       = "04" // 分
	Min_CN    = Min + "分"
	Sec       = "05" // 秒
	Sec_CN    = Sec + "秒"
	YM        = Y + M
	MD        = M + D
	YMD       = YM + D
	Y_M       = Y + "-" + M
	M_D       = M + "-" + D
	Y_M_D     = Y_M + "-" + D
	YM_CN     = Y_CN + M_CN
	MD_CN     = M_CN + D_CN
	YMD_CN    = YM_CN + D_CN
	YMDHM_CN  = YMD_CN + Hour_CN + Min_CN
	YMDHMS_CN = YMDHM_CN + Sec_CN
)

const (
	Day        = 24 * time.Hour // 一天
	ThreeDays  = 3 * Day        // 三天
	Week       = 7 * Day        // 一周
	Feb        = 28 * Day       // 二月
	LeapFeb    = 29 * Day       // 闰二月
	SmallMonth = 30 * Day       // 小月
	ThirtyDays = SmallMonth     // 三十天
	Month      = SmallMonth     // 月
	BigMonth   = 31 * Day       // 大月
	Year       = 365 * Day      // 一年
	LeapYear   = 366 * Day      // 闰年
)

const (
	Secc      = 1
	MinuteSec = 60 * Secc
	HourSec   = 60 * MinuteSec
	DaySec    = 24 * HourSec // 一天
	WeekSec   = 7 * DaySec   // 一周
	MonthSec  = 30 * DaySec  // 一个月
	YearSec   = 365 * DaySec // 一年
)

func IsCrossDay(begin int64, ends ...int64) bool {
	end := time.Now().Unix()
	if len(ends) > 0 {
		end = ends[0]
	}
	bstr := time.Unix(begin, 0).Format(YMD)
	estr := time.Unix(end, 0).Format(YMD)
	return bstr < estr
}

func DatetimeToUnix(date string) int64 {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", date, loc)
	if err != nil {
		t, _ = time.ParseInLocation("2006-01-02", date, loc)
		return t.Unix()
	}
	return t.Unix()
}

func DateToUnix(date string) int64 {
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02", date, loc)
	return t.Unix()
}

func UnixToDate(unix int64) string {
	return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
}

func TodayBegin() time.Time {
	return Begin(time.Now())
}

func Begin(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Local")
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func TodayEnd() time.Time {
	return End(time.Now())
}

func End(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Local")
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, loc)
}

// 获取最近几天相同时刻
func GetLatestDaysBaseNow(days int, desc ...bool) []time.Time {
	return GetLatestDays(time.Now(), days, desc...)
}

// 获取自某个时刻最近几天的相同时刻（包含当日）
// base 起点时刻
// days 最近的几天，< 0 表示过去几天；> 0 表示之后几天
// desc[0] = true 倒序输出
func GetLatestDays(base time.Time, days int, desc ...bool) (dates []time.Time) {
	isAfter := days > 0
	isDesc := false
	if len(desc) > 0 && desc[0] == true {
		isDesc = true
	}

	abs := int(math.Abs(float64(days)))
	if abs == 0 {
		abs++
	}
	dates = make([]time.Time, abs)

	for i := range dates {
		if isAfter {
			if isDesc {
				// desc
				dates[i] = base.AddDate(0, 0, abs-1-i)
			} else {
				// asc
				dates[i] = base.AddDate(0, 0, i)
			}
		} else {
			if isDesc {
				// desc
				dates[i] = base.AddDate(0, 0, -i)
			} else {
				// asc
				dates[i] = base.AddDate(0, 0, -abs+1+i)
			}
		}
	}

	return dates
}

// 获取该时间当天的起止时间点
func GetThisDayBeginEnd(t time.Time) (begin, end time.Time) {
	return Begin(t), End(t)
}

func DurationUntilTomorrow() time.Duration {
	loc, _ := time.LoadLocation("Local")
	tomorrow := time.Now().Add(time.Hour * 24)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, loc)
	return time.Until(tomorrow)
}

func GetNowUnix() int64 {
	return time.Now().Unix()
}

func GetNowUnixMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetNowUnixNano() int64 {
	return time.Now().UnixNano()
}

// 获取本周所在年份、第几周，起止时间点
func WeekInfo(t ...time.Time) (year, isoWeek int, weekBegin, weekEnd time.Time) {
	now := time.Now()
	if len(t) > 0 && !t[0].IsZero() {
		now = t[0]
	}

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekBegin = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	year, isoWeek = weekBegin.ISOWeek()
	weekEnd = weekBegin.AddDate(0, 0, 7).Add(time.Hour * 24).Add(-time.Second)

	return
}

// 简易的年月日转字符串
// 使用前需自行校验 year, month, day 有效性
func YmdStr(year, month, day int, sep ...string) string {
	sym := "-"
	if len(sep) > 0 {
		sym = sep[0]
	}

	ymd := []string{strconv.Itoa(year), strconv.Itoa(month), strconv.Itoa(day)}

	return strings.Join(ymd, sym)
}

func SubDay(time1, time2 time.Time) int64 {
	return int64(time1.Sub(time2).Hours() / 24)
}

func GetDayHourNum(t int64) int {
	return time.Unix(t, 0).Hour()
}

func GetOneDaySec(t int64) int64 {
	return t % int64(Day)
}
