package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type ParentChild struct {
	parentClassName string
	className       string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run ./main.go <PATH>\n")
		return
	}
	path := os.Args[1]
	fmt.Printf("SourceCodeCrawler. path: %s\n", path)

	javaFiles := make([]string, 0)
	collectJavaFiles := func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil || fileInfo.IsDir() {
			return nil
		}
		if strings.HasSuffix(fileInfo.Name(), ".java") {
			javaFiles = append(javaFiles, path)
		}
		return nil
	}
	err := filepath.Walk(path, collectJavaFiles)
	if err != nil {
		fmt.Printf("Error collecting java files: %v\n", err)
	}
	fmt.Printf("%d java files found\n", len(javaFiles))

	classMap := make(map[string][]string)
	classesChannel := make(chan ParentChild)
	go func() {
		for {
			parentChild, ok := <-classesChannel
			if !ok {
				break
			}
			t := classMap[parentChild.parentClassName]
			t = append(t, parentChild.className)
			classMap[parentChild.parentClassName] = t
		}
	}()

	waitCroup := sync.WaitGroup{}
	for _, javaFile := range javaFiles {
		waitCroup.Add(1)
		go processFile(javaFile, &waitCroup, classesChannel)
	}
	waitCroup.Wait()
	close(classesChannel)

	for parentClassName, className := range classMap {
		fmt.Printf("%s: %v\n", parentClassName, className)
	}
}

func processFile(path string, waitCroup *sync.WaitGroup, consumer chan ParentChild) {
	fmt.Printf("Process file. path: %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error open file. path: %s\n", path)
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		ok, className, parentClassNames := getClass(line)
		if ok {
			fmt.Printf("Class found. className: %s\n", className)
			for _, parentClassName := range parentClassNames {
				consumer <- ParentChild{parentClassName, className}
			}
		}
	}

	file.Close()

	waitCroup.Done()
}

func getClass(line string) (bool, string, []string) {
	re := regexp.MustCompile(`\s*(public)?\s*(final|abstract)?\s*(class|interface|enum)\s*(?P<className>\w*)\s*((?:extends|implements)\s*(?P<parentClasses>.*))?{`)
	if !re.MatchString(line) {
		return false, "", nil
	}

	var className, parentClassesString string
	match := re.FindStringSubmatch(line)
	for i, groupName := range re.SubexpNames() {
		if i != 0 {
			t := match[i]
			switch groupName {
			case "className":
				className = t
			case "parentClasses":
				parentClassesString = t
			}
		}
	}

	parentClassNames := make([]string, 0)
	if len(parentClassesString) > 0 {
		for _, parentClassName := range strings.Split(parentClassesString, ",") {
			parentClassNames = append(parentClassNames, strings.TrimSpace(parentClassName))
		}
	}
	return len(className) > 0, className, parentClassNames
}
