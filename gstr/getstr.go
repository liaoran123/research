package gstr

import "strings"

func Do(text string, lstr string, rstr string, llast bool, rlast bool) string {

	lpos, rpos := -1, -1
	if lstr == "" {
		lpos = 0
	} else {
		if !llast {
			lpos = strings.Index(text, lstr)
		} else {
			lpos = strings.LastIndex(text, lstr)
		}
		if lpos != -1 {
			lpos = lpos + len(lstr)
		}
	}
	if lpos == -1 {
		return ""
	}
	tmptext := text[lpos:]
	if rstr == "" {
		return tmptext
	} else {
		if !rlast {
			rpos = strings.Index(tmptext, rstr)
		} else {
			rpos = strings.LastIndex(tmptext, rstr)
		}

	}
	if rpos == -1 {
		return ""
	}
	return tmptext[:rpos]
}

func LStr(text, str string) string {
	return Do(text, "", str, false, false)
}
func RStr(text, str string) string {
	return Do(text, str, "", false, false)
}
func Mstr(text, ls, rs string) string {
	return Do(text, ls, rs, false, false)
}
