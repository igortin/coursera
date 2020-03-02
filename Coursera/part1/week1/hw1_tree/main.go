package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

const (
	sep = string(os.PathSeparator)
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	Tree(out, path, printFiles, "")
	return nil
}

func Tree(out io.Writer, path string, printFiles bool, space string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	objSlice, _ := file.Readdir(-1)

	sort.SliceStable(objSlice, func(i, j int) bool {
		return objSlice[i].Name() < objSlice[j].Name()
	})

	objects := make([]os.FileInfo, 0)
	if !printFiles {
		for _, fd := range objSlice {
			if fd.IsDir() {
				objects = append(objects, fd)
			}
		}
	} else {
		objects = objSlice[:]
	}

	for index, fd := range objects {
		if index != len(objects)-1 {
			if !fd.IsDir() {
				fmt.Fprintf(out, "%s%s %s\n", space+"├───", fd.Name(), getS(fd.Size()))

				//fmt.Printf("%s%s %s\n", space+"├───", fd.Name(), getS(fd.Size()))
			} else {
				//fmt.Printf("%s%s\n", space+"├───", fd.Name())
				fmt.Fprintf(out, "%s%s\n", space+"├───", fd.Name())
			}
			Tree(out, path+sep+fd.Name(), printFiles, space+"│    ")
		} else {
			if !fd.IsDir() {
				//fmt.Printf("%s%s %s\n", space+"└───", fd.Name(), getS(fd.Size()))
				fmt.Fprintf(out, "%s%s %s\n", space+"└───", fd.Name(), getS(fd.Size()))
			} else {
				//fmt.Printf("%s%s\n", space+"└───", fd.Name())
				fmt.Fprintf(out, "%s%s\n", space+"└───", fd.Name())
			}
			Tree(out, path+sep+fd.Name(), printFiles, space+"       ")
		}
	}
	return nil
}

func getS(v int64) string {
	var s string
	if v == 0 {
		s = "(empty)"
	} else {
		s = "(" + strconv.FormatInt(v, 10) + "b" + ")"
	}
	return s
}
