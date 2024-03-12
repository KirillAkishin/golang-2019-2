package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	layer := 0
	workingLayers := *new([]bool)
	err := readDirectory(layer, out, path, printFiles, &workingLayers)
	return err
}

func readDirectory(layer int, out io.Writer, path string, printFiles bool, ptrToWrkLayers *[]bool) error {
	workingLayers := *ptrToWrkLayers
	workingLayers = append(workingLayers, true)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	listOfNames, err := setList(file, printFiles, path)
	if err != nil {
		return err
	}
	for cnt, name := range listOfNames {
		if cnt == len(listOfNames)-1 {
			workingLayers[layer] = false
		}

		currentPath := fmt.Sprintf("%v/%v", path, name)
		nameInfo, err := os.Stat(currentPath)
		if err != nil {
			return err
		}

		if nameInfo.IsDir() {
			prefix := setPrefix(layer, cnt, len(listOfNames), workingLayers)
			currentOut := fmt.Sprintf("%s%s\n", prefix, name)
			out.Write([]byte(currentOut))

			readDirectory(layer+1, out, currentPath, printFiles, &workingLayers)
		} else if printFiles {
			prefix := setPrefix(layer, cnt, len(listOfNames), workingLayers)
			postfix := setPostfix(nameInfo.Size())
			currentOut := fmt.Sprintf("%s%s%s\n", prefix, name, postfix)
			out.Write([]byte(currentOut))
		}
	}
	return nil
}

func setList(file *os.File, printFiles bool, path string) ([]string, error) {
	listOfNames, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	newListOfNames := *new([]string)
	if !printFiles {
		for _, name := range listOfNames {
			currentPath := fmt.Sprintf("%v/%v", path, name)
			nameInfo, _ := os.Stat(currentPath)
			if err != nil {
				return nil, err
			}
			if nameInfo.IsDir() {
				newListOfNames = append(newListOfNames, name)
			}
		}
	} else {
		newListOfNames = listOfNames
	}
	sort.Strings(newListOfNames)
	return newListOfNames, nil
}

func setPrefix(layer int, cnt int, lenList int, workingLayers []bool) (prefix string) {
	prefix = ""

	for i := 0; i < layer; i++ {
		if workingLayers[i] {
			prefix += "│"
		}
		prefix += "\t"
	}

	if cnt != lenList-1 {
		prefix += "├"
	} else {
		prefix += "└"
	}

	prefix += "───"
	return prefix
}

func setPostfix(size int64) (postfix string) {
	if size != 0 {
		postfix = fmt.Sprintf(" (%vb)", size)
	} else {
		postfix = fmt.Sprintf(" (empty)")
	}
	return postfix
}
