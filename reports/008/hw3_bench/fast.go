package main

import (
	"bufio"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Close allows not to ignore potential error inside `defer`.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}

// readFileToChan reads line by line from given `path` and writes to `lines` channel.
func readFileToChan(path string, lines chan<- string) {
	defer close(lines)
	f, err := os.Open(path)
	defer Close(f)
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines <- s.Text()
		runtime.Gosched()
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
}

// FastSearch reads log output from `filePath` and outputs all unique users and browsers.
func FastSearch(out io.Writer) {
	lines := make(chan string, 100)
	go readFileToChan(filePath, lines)

	seenBrowsers, uniqueBrowsers, foundUsers, isAndroid, isMSIE, user, i, j :=
		make(map[string]struct{}, 128), 0, make([]byte, 0, 1024), false, false, &User{}, int64(0), 0
	for line := range lines {
		err := easyjson.Unmarshal([]byte(line), user)
		if err != nil {
			panic(err)
		}

		isAndroid, isMSIE = false, false
		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}
		}
		if isAndroid && isMSIE {
			j = strings.Index(user.Email, "@")
			foundUsers = append(foundUsers,[]byte("[")...)
			foundUsers = strconv.AppendInt(foundUsers, i, 10)
			foundUsers = append(foundUsers,[]byte("] ")...)
			foundUsers = append(foundUsers,[]byte(user.Name)...)
			foundUsers = append(foundUsers,[]byte(" <")...)
			foundUsers = append(foundUsers,[]byte(user.Email[:j])...)
			foundUsers = append(foundUsers,[]byte(" [at] ")...)
			foundUsers = append(foundUsers,[]byte(user.Email[j+1:])...)
			foundUsers = append(foundUsers,[]byte(">\n")...)
		}
		i++
	}

	err := write(out, "found users:\n"+string(foundUsers)+"\nTotal unique browsers "+strconv.Itoa(len(seenBrowsers))+"\n")
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
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					out.Browsers = append(out.Browsers, string(in.String()))
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
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