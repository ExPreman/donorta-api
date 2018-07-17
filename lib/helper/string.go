package helper

import (
	"fmt"
	"time"
	"regexp"
	"strings"
	"unicode"
	"math/rand"
	"encoding/json"
)

const STRING_CHARSET = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const STRING_NUMBER = "0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringRandomWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func StringRandom(length int) string {
	return StringRandomWithCharset(length, STRING_CHARSET)
}

func StringIsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func FormatToCurrency(n int) string {
	var s []string
	is := fmt.Sprintf("%d", n)

	for i := len(is); i > 0; i -= 3 {
		switch {
		case i >= 3:
			s = append([]string{is[i-3 : i]}, s...)
		case i < 3:
			s = append([]string{is[:i]}, s...)
		}
	}

	return strings.Join(s, ".")
}

func WhiteSpaceRemove(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func CleanHPNo (hp string) string {
	if len(hp) > 2 {
		hp = strings.Replace(hp, " ", "", -1)
		hp = strings.Replace(hp, "+", "", -1)
		hp = strings.Replace(hp, "-", "", -1)
		if hp[:1] == "0" {
			hp = "62"+ hp[1:]
		} else if notNumber := regexp.MustCompile("[a-zA-Z]+").MatchString; !notNumber(hp) && hp[:2] != "62" {
			hp = "62"+ hp
		}
	}
	return hp
}