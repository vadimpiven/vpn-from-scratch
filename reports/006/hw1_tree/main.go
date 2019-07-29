package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// filter deletes files from list if they shouldn't be printed, removes entries that break tests and should be ignored.
func filter(dirInfo []os.FileInfo, printFiles bool) []os.FileInfo {
	newDirInfo := make([]os.FileInfo, 0, len(dirInfo))
	for _, file := range dirInfo {
		name := file.Name()
		if name == ".DS_Store" || name == ".git" || name == ".idea" || (!file.IsDir() && !printFiles) {
			continue
		}
		newDirInfo = append(newDirInfo, file)
	}
	return newDirInfo
}

// write performs buffered write of entire string.
func write(out io.Writer, str string) (err error) {
	for n, err := io.WriteString(out, str); err == nil && n < len(str); {
		str = str[n:]
	}
	return err
}

// buildDirTree recursively outputs the directory tree
func buildDirTree(out io.Writer, path, prefix string, printFiles bool) error {
	dirInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	dirInfo = filter(dirInfo, printFiles)
	var newPrefix, line, name string
	for i, info := range dirInfo {
		if i == len(dirInfo)-1 {
			newPrefix, line = prefix+"\t", "└───"
		} else {
			newPrefix, line = prefix+"│\t", "├───"
		}
		name = info.Name()
		if info.IsDir() {
			if err = write(out, prefix+line+name+"\n"); err != nil {
				return err
			}
			err = buildDirTree(out, filepath.Join(path, name), newPrefix, printFiles) // path+string(os.PathSeparator)+name
		} else {
			if info.Size() > 0 {
				err = write(out, prefix+line+name+" ("+strconv.FormatInt(info.Size(), 10)+"b)\n")
			} else {
				err = write(out, prefix+line+name+" (empty)\n")
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// dirTree checks if path specified is pointing to directory and runs the recursive tree building function.
func dirTree(out io.Writer, path string, printFiles bool) error {
	p, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !p.Mode().IsDir() {
		return errors.New("given path is not a directory")
	}
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
