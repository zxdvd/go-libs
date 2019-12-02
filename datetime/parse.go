package datetime

import (
	"errors"
	"time"
)

func Parse(tstr string, layout string) (time.Time, error) {
	var year, month, day, yearday, hour, minute, second, nano, zoneoffset int
	var pmSet, amSet bool
	var err error
	month = 1
	day = 1
	zoneoffset = -1
	for layout != "" {
		prefix, curTok, suffix := nextStdChunk(layout)
		tstr = tstr[len(prefix):]

		if curTok == tokNop {
			break
		}
		layout = suffix

		switch curTok {
		case tokM, tokMo:
			month, tstr, err = getnum(tstr, 2)
			if curTok == tokMo {
				tstr = tstr[2:] // skip the suffix of 1st, 2nd...
			}
		case tokMM:
			month, tstr, err = getnumExactly(tstr, 2)
		case tokMMM:
			month, tstr, err = lookup(shortMonthNames, tstr)
			month++
		case tokMMMM:
			month, tstr, err = lookup(longMonthNames, tstr)
			month++
		case tokD, tokDo:
			day, tstr, err = getnum(tstr, 3)
			if curTok == tokDo {
				tstr = tstr[2:] // skip the suffix of 1st, 2nd...
			}
		case tokDD:
			day, tstr, err = getnumExactly(tstr, 2)
		case tokDDD, tokDDDo:
			yearday, tstr, err = getnum(tstr, 3)
			if curTok == tokDDDo {
				tstr = tstr[2:] // skip the suffix of 1st, 2nd...
			}
		case tokDDDD:
			yearday, tstr, err = getnumExactly(tstr, 3)
		case tokddd: // TODO dealwith weekday
			_, tstr, err = lookup(shortDayNames, tstr)
		case tokdddd:
			_, tstr, err = lookup(longDayNames, tstr)
		case tokYY, tokGG:
			year, tstr, err = getnumExactly(tstr, 2)
			if year >= 69 { // Unix time starts Dec 31 1969 in some time zones
				year += 1900
			} else {
				year += 2000
			}
		case tokYYYY, tokGGGG:
			year, tstr, err = getnumExactly(tstr, 4)
		case tokA, toka:
			if len(tstr) < 2 {
				err = errBad
				break
			}
			var p string
			p, tstr = tstr[0:2], tstr[2:]
			switch p {
			case "PM", "pm":
				pmSet = true
			case "AM", "am":
				amSet = true
			default:
				err = errBad
			}
		case tokH:
			hour, tstr, err = getnum(tstr, 2)
		case tokHH:
			hour, tstr, err = getnumExactly(tstr, 2)
		case tokh:
			hour, tstr, err = getnum(tstr, 2)
		case tokhh:
			hour, tstr, err = getnumExactly(tstr, 2)
		case tokk:
			hour, tstr, err = getnum(tstr, 2)
			if hour == 24 {
				hour = 0
			}
		case tokkk:
			hour, tstr, err = getnumExactly(tstr, 2)
			if hour == 24 {
				hour = 0
			}
		case tokm:
			minute, tstr, err = getnum(tstr, 2)
		case tokmm:
			minute, tstr, err = getnumExactly(tstr, 2)
		case toks:
			second, tstr, err = getnum(tstr, 2)
		case tokss:
			second, tstr, err = getnumExactly(tstr, 2)
		case tokS:
			nano, tstr, err = getnumExactly(tstr, 1)
			nano = nano * 1e8
		case tokSS:
			nano, tstr, err = getnumExactly(tstr, 2)
			nano = nano * 1e7
		case tokSSS:
			nano, tstr, err = getnumExactly(tstr, 3)
			nano = nano * 1e6
		case tokz, tokzz:
			//TODO
		case tokZ, tokZZ:
			if len(tstr) >= 1 && tstr[0] == 'Z' { // UTC
				zoneoffset = 0
				break
			}
			var zoneSign, zoneHour, zoneMin string
			if curTok == tokZ { // +08:00 style
				if len(tstr) < 6 {
					err = errBad
					break
				}
				zoneSign, zoneHour, zoneMin = tstr[:1], tstr[1:3], tstr[4:6]
			} else { // +0800 style
				if len(tstr) < 5 {
					err = errBad
					break
				}
				zoneSign, zoneHour, zoneMin = tstr[:1], tstr[1:3], tstr[3:5]
			}
			var hr, m int
			hr, _, err = getnumExactly(zoneHour, 2)
			if err == nil {
				m, _, err = getnumExactly(zoneMin, 2)
			}
			zoneoffset = hr*3600 + m*60
			if zoneSign == "-" {
				zoneoffset = -zoneoffset
			} else if zoneSign != "+" {
				err = errors.New("wrong zone sign")
			}
		}
		if err != nil {
			return time.Time{}, err
		}
		if month <= 0 || month > 12 {
			return time.Time{}, errors.New("month, wrong range")
		}
		if day <= 0 {
			return time.Time{}, errors.New("day, wrong range")
		}
	}
	if pmSet && hour < 12 {
		hour += 12
	} else if amSet && hour == 12 {
		hour = 0
	}
	var t time.Time
	if yearday > 0 {
		t = time.Date(year, time.Month(1), 1, hour, minute, second, nano, time.UTC)
		t.AddDate(0, 0, yearday)
	}
	t = time.Date(year, time.Month(month), day, hour, minute, second, nano, time.UTC)
	if zoneoffset != -1 {
		t.Add(time.Second * time.Duration(-zoneoffset))
	}
	return t, nil

}

func getnum(s string, maxlen int) (int, string, error) {
	if !isDigit(s, 0) {
		return 0, s, errBad
	}
	val := int(s[0] - '0')
	var i int
	for i = 1; i < maxlen; i++ {
		if !isDigit(s, i) {
			break
		}
		val = 10*val + int(s[i]-'0')
	}
	return val, s[i:], nil
}

func getnumExactly(s string, length int) (int, string, error) {
	val := 0
	var i int
	for i = 0; i < length; i++ {
		if !isDigit(s, i) {
			return 0, s, errBad
		}
		val = 10*val + int(s[i]-'0')
	}
	return val, s[i:], nil
}

// copy from src/time/format.go
// isDigit reports whether s[i] is in range and is a decimal digit.

var errBad = errors.New("bad value for field")

func isDigit(s string, i int) bool {
	if len(s) <= i {
		return false
	}
	c := s[i]
	return '0' <= c && c <= '9'
}

func match(s1, s2 string) bool {
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 != c2 {
			// Switch to lower-case; 'a'-'A' is known to be a single bit.
			c1 |= 'a' - 'A'
			c2 |= 'a' - 'A'
			if c1 != c2 || c1 < 'a' || c1 > 'z' {
				return false
			}
		}
	}
	return true
}

func lookup(tab []string, val string) (int, string, error) {
	for i, v := range tab {
		if len(val) >= len(v) && match(val[0:len(v)], v) {
			return i, val[len(v):], nil
		}
	}
	return -1, val, errBad
}
