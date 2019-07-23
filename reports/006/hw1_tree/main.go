package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type PathInfo struct {
	Name  string
	IsDir bool
}

func filter(names []string, path string, printFiles bool) ([]PathInfo, error) {
	sort.Strings(names)
	info := make([]PathInfo, 0)
	for _, name := range names {
		if name == ".DS_Store" || name == ".git" || name == ".idea" {
			continue
		}
		p, err := os.Stat(path + string(os.PathSeparator) + name)
		if err != nil {
			return info, err
		}
		if p.Mode().IsDir() {
			info = append(info, PathInfo{name, true})
		} else if printFiles{
			if p.Size() > 0 {
				info = append(info, PathInfo{name + " (" + strconv.FormatInt(p.Size(), 10) + "b)", false})
			} else {
				info = append(info, PathInfo{name + " (empty)", false})
			}
		}
	}
	return info, nil
}

func write(out io.Writer, str string) error {
	for n, err := fmt.Fprint(out, str); n < len(str); {
		if err != nil {
			return err
		}
		str = str[n:]
	}
	return nil
}

func buildDirTree(out io.Writer, path, prefix string, printFiles bool) error {
	p, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !p.Mode().IsDir() {
		return errors.New("given path is not a directory")
	}
	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return err
	}
	names, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}
	info, err := filter(names, path, printFiles)
	if err != nil {
		return err
	}
	var newPrefix, line string
	for i, pathInfo := range info {
		if i == len(info) - 1 {
			newPrefix, line = prefix + "\t", "└───"
		} else {
			newPrefix, line = prefix + "│\t", "├───"
		}
		err = write(out, prefix + line + pathInfo.Name + "\n")
		if err != nil {
			return err
		}
		if pathInfo.IsDir {
			err = buildDirTree(out, path + string(os.PathSeparator) + pathInfo.Name, newPrefix, printFiles)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return buildDirTree(out, path, "", printFiles)
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
