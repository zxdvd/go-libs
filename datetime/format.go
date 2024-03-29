package datetime

import (
	"strconv"
	"strings"
	"time"
)

type tok uint8

const (
	tokNop tok = iota
	tokM
	tokMo
	tokMM
	tokMMM
	tokMMMM
	tokD
	tokDo
	tokDD
	tokDDD
	tokDDDo
	tokDDDD
	tokd
	tokdo
	tokdd
	tokddd
	tokdddd
	tokE
	tokW
	tokY
	tokYY
	tokYYYY
	tokGG
	tokGGGG
	tokA
	toka
	tokH
	tokHH
	tokh
	tokhh
	tokk
	tokkk
	tokm
	tokmm
	toks
	tokss
	tokS
	tokSS
	tokSSS
	tokz
	tokzz
	tokZ
	tokZZ
)

var longDayNames = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var shortDayNames = []string{
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
}

var shortMonthNames = []string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

var longMonthNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

// learn from golang:src/time/format.go
// format learn from moment.js https://devdocs.io/moment/index#displaying-format
func nextStdChunk(layout string) (prefix string, std tok, suffix string) {
	for i := 0; i < len(layout); i++ {
		switch c := int(layout[i]); c {
		case 'M':
			if strings.HasPrefix(layout[i:], "MMMM") {
				return layout[0:i], tokMMMM, layout[i+4:]
			}
			if strings.HasPrefix(layout[i:], "MMM") {
				return layout[0:i], tokMMM, layout[i+3:]
			}
			if strings.HasPrefix(layout[i:], "MM") {
				return layout[0:i], tokMM, layout[i+2:]
			}
			if strings.HasPrefix(layout[i:], "Mo") {
				return layout[0:i], tokMo, layout[i+2:]
			}
			return layout[0:i], tokM, layout[i+1:]
		case 'D':
			if strings.HasPrefix(layout[i:], "DDDD") {
				return layout[0:i], tokDDDD, layout[i+4:]
			}
			if strings.HasPrefix(layout[i:], "DDDo") {
				return layout[0:i], tokDDDo, layout[i+4:]
			}
			if strings.HasPrefix(layout[i:], "DDD") {
				return layout[0:i], tokDDD, layout[i+3:]
			}
			if strings.HasPrefix(layout[i:], "DD") {
				return layout[0:i], tokDD, layout[i+2:]
			}
			if strings.HasPrefix(layout[i:], "Do") {
				return layout[0:i], tokDo, layout[i+2:]
			}
			return layout[0:i], tokD, layout[i+1:]
		case 'd':
			if strings.HasPrefix(layout[i:], "dddd") {
				return layout[0:i], tokdddd, layout[i+4:]
			}
			if strings.HasPrefix(layout[i:], "ddd") {
				return layout[0:i], tokddd, layout[i+3:]
			}
			if strings.HasPrefix(layout[i:], "dd") {
				return layout[0:i], tokdd, layout[i+2:]
			}
			if strings.HasPrefix(layout[i:], "do") {
				return layout[0:i], tokdo, layout[i+2:]
			}
			return layout[0:i], tokd, layout[i+1:]
		case 'E':
			return layout[0:i], tokE, layout[i+1:]
		case 'W':
			return layout[0:i], tokW, layout[i+1:]
		case 'Y':
			if strings.HasPrefix(layout[i:], "YYYY") {
				return layout[0:i], tokYYYY, layout[i+4:]
			}
			if strings.HasPrefix(layout[i:], "YY") {
				return layout[0:i], tokYY, layout[i+2:]
			}
		case 'G':
			if strings.HasPrefix(layout[i:], "GGGG") {
				return layout[0:i], tokGGGG, layout[i+4:]
			}
			if strings.HasPrefix(layout[i:], "GG") {
				return layout[0:i], tokGG, layout[i+2:]
			}
		case 'A': // AM PM
			return layout[0:i], tokA, layout[i+1:]
		case 'a': // am pm
			return layout[0:i], toka, layout[i+1:]
		case 'H': // 0 1 ... 22 23 and 00 01 ... 22 23
			if strings.HasPrefix(layout[i:], "HH") {
				return layout[0:i], tokHH, layout[i+2:]
			}
			return layout[0:i], tokH, layout[i+1:]
		case 'h': // 1 2 ... 11 12 and 01 02 ... 11 12
			if strings.HasPrefix(layout[i:], "hh") {
				return layout[0:i], tokhh, layout[i+2:]
			}
			return layout[0:i], tokh, layout[i+1:]
		case 'k': // 1 2 3...23 24 and 01 02 ... 23 24
			if strings.HasPrefix(layout[i:], "kk") {
				return layout[0:i], tokkk, layout[i+2:]
			}
			return layout[0:i], tokk, layout[i+1:]
		case 'm': // 0 1 ... 58 59 and 00 01 ... 58 59
			if strings.HasPrefix(layout[i:], "mm") {
				return layout[0:i], tokmm, layout[i+2:]
			}
			return layout[0:i], tokm, layout[i+1:]
		case 's': // 0 1 ... 58 59 and 00 01 ... 58 59
			if strings.HasPrefix(layout[i:], "ss") {
				return layout[0:i], tokss, layout[i+2:]
			}
			return layout[0:i], toks, layout[i+1:]
		case 'S': // 0 1 ... 8 9 and 00 01 98 99 and 000 001 ... 998 999
			if strings.HasPrefix(layout[i:], "SSS") {
				return layout[0:i], tokSSS, layout[i+3:]
			}
			if strings.HasPrefix(layout[i:], "SS") {
				return layout[0:i], tokSS, layout[i+2:]
			}
			return layout[0:i], tokS, layout[i+1:]
		case 'Z': // Z +08:00, ZZ +0800
			if strings.HasPrefix(layout[i:], "ZZ") {
				return layout[0:i], tokZZ, layout[i+2:]
			}
			return layout[0:i], tokZ, layout[i+1:]

		}
	}
	return layout, tokNop, ""
}

func Format(t time.Time, layout string) string {
	buflen := len(layout) + 10
	var b []byte
	var layoutBuffer [64]byte
	if buflen < 64 {
		b = layoutBuffer[:0]
	} else {
		b = make([]byte, 0, buflen)
	}
	for layout != "" {
		prefix, curTok, suffix := nextStdChunk(layout)
		if prefix != "" {
			b = append(b, prefix...)
		}
		if curTok == tokNop {
			break
		}
		layout = suffix

		switch curTok {
		case tokM:
			b = appendInt(b, int(t.Month()), 0)
		case tokMo:
			b = appendOrdStr(b, int(t.Month()))
		case tokMM:
			b = appendInt(b, int(t.Month()), 2)
		case tokMMM:
			s := shortMonthNames[int(t.Month())]
			b = append(b, s...)
		case tokMMMM:
			s := longMonthNames[t.Month()]
			b = append(b, s...)
		case tokD:
			b = appendInt(b, t.Day(), 0)
		case tokDo:
			b = appendOrdStr(b, int(t.Day()))
		case tokDD:
			b = appendInt(b, t.Day(), 2)
		case tokDDD:
			b = appendInt(b, t.YearDay(), 0)
		case tokDDDo:
			b = appendOrdStr(b, int(t.YearDay()))
		case tokDDDD:
			b = appendInt(b, t.YearDay(), 3)
		case tokd:
			b = appendInt(b, int(t.Weekday()), 0)
		case tokdo:
			b = appendOrdStr(b, int(t.Weekday()))
		case tokdd:
			s := shortDayNames[t.Weekday()]
			b = append(b, []byte(s)[0:2]...)
		case tokddd:
			s := shortDayNames[int(t.Weekday())]
			b = append(b, []byte(s)...)
		case tokdddd:
			s := longDayNames[int(t.Weekday())]
			b = append(b, []byte(s)...)
		case tokYY, tokGG:
			b = appendInt(b, t.Year()%100, 2)
		case tokYYYY, tokGGGG:
			b = appendInt(b, t.Year(), 4)
		case tokA:
			if t.Hour() < 12 {
				b = append(b, "AM"...)
			} else {
				b = append(b, "PM"...)
			}
		case toka:
			if t.Hour() < 12 {
				b = append(b, "am"...)
			} else {
				b = append(b, "pm"...)
			}
		case tokH:
			b = appendInt(b, t.Hour(), 0)
		case tokHH:
			b = appendInt(b, t.Hour(), 2)
		case tokh:
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			b = appendInt(b, h, 0)
		case tokhh:
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			b = appendInt(b, h, 2)
		case tokk:
			h := t.Hour()
			if h == 0 {
				h = 24
			}
			b = appendInt(b, h, 0)
		case tokkk:
			h := t.Hour()
			if h == 0 {
				h = 24
			}
			b = appendInt(b, h, 2)
		case tokm:
			b = appendInt(b, t.Minute(), 0)
		case tokmm:
			b = appendInt(b, t.Minute(), 2)
		case toks:
			b = appendInt(b, t.Second(), 0)
		case tokss:
			b = appendInt(b, t.Second(), 2)
		case tokS:
			nano := t.Nanosecond()
			b = appendInt(b, nano/1e8, 1)
		case tokSS:
			nano := t.Nanosecond()
			b = appendInt(b, nano/1e7, 2)
		case tokSSS:
			nano := t.Nanosecond()
			b = appendInt(b, nano/1e6, 3)
		case tokz, tokzz:
			zone, _ := t.Zone()
			b = append(b, zone...)
		case tokZ, tokZZ:
			_, offset := t.Zone()
			offset = offset / 60 // convert seconds to minutes
			if offset < 0 {
				b = append(b, '-')
				offset = -offset
			} else {
				b = append(b, '+')
			}
			b = appendInt(b, offset/60, 2)
			if curTok == tokZ {
				b = append(b, ':') // tokZ +08:00
			}
			b = appendInt(b, offset%60, 2)
		}
	}

	return string(b)
}

func appendInt(b []byte, x int, width int) []byte {
	var buffer [20]byte
	buf := buffer[:0]
	buf = strconv.AppendInt(buf, int64(x), 10)
	buf = prepad0(buf, width)
	b = append(b, buf...)
	return b
}

func prepad0(b []byte, width int) []byte {
	return prepad(b, width, '0')
}

func prepad(b []byte, width int, p byte) []byte {
	l := len(b)
	padlen := width - l
	if padlen <= 0 {
		return b
	}
	if width > cap(b) {
		buf := make([]byte, 0, width)
		for i := 0; i < width-l; i++ {
			buf = append(buf, p)
		}
		return append(buf, b...)
	}

	// width <= cap(b) then we can safe expand b to width length
	b = b[:width]
	// move the origin to end then we can prepend char at begin
	for i := 0; i < l; i++ {
		b[width-1-i] = b[l-1-i]
	}
	for i := 0; i < padlen; i++ {
		b[i] = p
	}
	return b
}

var ordNumMap = map[int]string{
	1: "st",
	2: "nd",
}

func appendOrdStr(b []byte, n int) []byte {
	b = appendInt(b, n, 0)
	return append(b, ordNumSuffix(n)...)
}

func ordNumSuffix(n int) string {
	if s, ok := ordNumMap[n]; ok {
		return s
	}
	return "th"
}
