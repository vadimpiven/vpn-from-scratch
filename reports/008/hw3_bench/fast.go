package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

// Close allows not to ignore potential error inside `defer`.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}

// FastSearch reads log output from `filePath` and outputs all unique users and browsers.
func FastSearch(out io.Writer) {
	f, err := os.Open(filePath)
	defer Close(f)
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(f)
	seenBrowsers, foundUsers, isAndroid, isMSIE, user, i, j :=
		make([]string, 0, 128), make([]byte, 0, 16384), false, false, &User{}, uint64(0), 0
	for ; s.Scan(); i++ {
		if bytes.Contains(s.Bytes(), []byte("Android")) || bytes.Contains(s.Bytes(), []byte("MSIE")) {
			err := easyjson.Unmarshal(s.Bytes(), user)
			if err != nil {
				panic(err)
			}
			isAndroid, isMSIE = false, false
			for _, browser := range user.Browsers {
				if strings.Contains(browser, "Android") {
					isAndroid = true
					notSeenBefore := true
					for _, item := range seenBrowsers {
						if item == browser {
							notSeenBefore = false
							break
						}
					}
					if notSeenBefore {
						seenBrowsers = append(seenBrowsers, (browser + "\000")[:len(browser)])
					}
				} else if strings.Contains(browser, "MSIE") {
					isMSIE = true
					notSeenBefore := true
					for _, item := range seenBrowsers {
						if item == browser {
							notSeenBefore = false
							break
						}
					}
					if notSeenBefore {
						seenBrowsers = append(seenBrowsers, (browser + "\000")[:len(browser)])
					}
				}
			}
			if isAndroid && isMSIE {
				j = strings.Index(user.Email, "@")
				foundUsers = append(foundUsers, []byte("[")...)
				foundUsers = strconv.AppendUint(foundUsers, i, 10)
				foundUsers = append(foundUsers, []byte("] ")...)
				foundUsers = append(foundUsers, []byte(user.Name)...)
				foundUsers = append(foundUsers, []byte(" <")...)
				foundUsers = append(foundUsers, []byte(user.Email[:j])...)
				foundUsers = append(foundUsers, []byte(" [at] ")...)
				foundUsers = append(foundUsers, []byte(user.Email[j+1:])...)
				foundUsers = append(foundUsers, []byte(">\n")...)
			}
		}
	}
	if err = s.Err(); err != nil {
		panic(err)
	}

	err = write(out, "found users:\n"+string(foundUsers)+"\nTotal unique browsers "+strconv.Itoa(len(seenBrowsers))+"\n")
	if err != nil {
		panic(err)
	}
}

// write performs buffered write of entire string.
func write(out io.Writer, str string) (err error) {
	var n int
	for n, err = io.WriteString(out, str); err == nil && n < len(str); {
		str = str[n:]
	}
	return
}

// User holds the useful part of data from one line of log file.
type User struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

// UnmarshalEasyJSON implements easyjson.Unmarshaler interface for type User.
func (out *User) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = (out.Browsers)[:0]
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					out.Browsers = make([]string, 0, 4)
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					out.Browsers = append(out.Browsers, in.UnsafeString())
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = in.UnsafeString()
		case "name":
			out.Name = in.UnsafeString()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
