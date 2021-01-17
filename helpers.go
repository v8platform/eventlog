package eventlog

import (
	"strconv"
	"time"
)

const unixTimeInSeconds = 62135596801

func SecondsToUnixTime(seconds int) time.Time {

	if seconds == 0 {
		return time.Time{}
	}

	seconds = seconds - unixTimeInSeconds

	return time.Unix(int64(seconds), 0).UTC()

}

func From16To10(str string) int64 {

	val, _ := strconv.ParseInt(str, 16, 64)

	return val
	//var val int64
	//
	//for i := 0; i < len(str)-1; i++ {
	//
	//	n := str[i:i+1]
	//
	//	idx := strings.Index(srt16, strings.ToUpper(n))
	//	if idx == -1 {
	//		continue
	//	}
	//
	//	val += int64(idx) * int64(math.Pow(16, float64(len(str)) - float64(i+1)))
	//
	//}
	//
	//return val
}
