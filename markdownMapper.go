package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Entry struct {
	level int
	text  string
}

func generateLink(file string, header string) string {
	headerText := strings.TrimSpace(strings.TrimLeft(header, "#"))
	headerID := strings.Replace(strings.ToLower(headerText), " ", "-", -1)
	link := fmt.Sprintf("[%s](./%s#%s)", headerText, file, headerID)
	return link
}

func generateTOC(file string) map[string][]Entry {
	toc := make(map[string][]Entry)
	f, _ := os.Open(file)
	scanner := bufio.NewScanner(f)

	headerRegex := regexp.MustCompile(`^(#{1,3})\s+`)
	mainHeading := ""

	for scanner.Scan() {
		line := scanner.Text()
		if headerRegex.MatchString(line) {
			level := len(headerRegex.FindStringSubmatch(line)[1])
			link := generateLink(file, line)
			indent := strings.Repeat("    ", level-1)
			if level == 1 {
				mainHeading = link
			} else if mainHeading != "" {
				toc[mainHeading] = append(toc[mainHeading], Entry{level: level, text: indent + "- " + link})
			}
		}
	}

	if mainHeading == "" {
		mainHeading = fmt.Sprintf("[%s](./%s)", strings.TrimSuffix(file, filepath.Ext(file)), file)
		toc[mainHeading] = []Entry{}
	}

	return toc
}

func main() {
	tocFilename := "table_of_contents.md"
	tocMap := make(map[string][]Entry)

	files, _ := ioutil.ReadDir(".")
	for _, f := range files {
		file := f.Name()
		if strings.HasSuffix(file, ".md") && file != tocFilename {
			fileToc := generateTOC(file)
			for mainHeading, subheadings := range fileToc {
				tocMap[mainHeading] = append(tocMap[mainHeading], subheadings...)
			}
		}
	}

	sortedToc := make([]string, 0, len(tocMap))
	for mainHeading := range tocMap {
		sortedToc = append(sortedToc, mainHeading)
	}
	sort.Slice(sortedToc, func(i, j int) bool {
		return strings.ToLower(sortedToc[i]) < strings.ToLower(sortedToc[j])
	})

	tocFile, _ := os.Create(tocFilename)
	defer tocFile.Close()

	tocWriter := bufio.NewWriter(tocFile)
	tocWriter.WriteString("Table of Contents:\n")
	for _, mainHeading := range sortedToc {
		tocWriter.WriteString("- " + mainHeading + "\n")
		for _, entry := range tocMap[mainHeading] {
			tocWriter.WriteString(entry.text + "\n")
		}
	}
	tocWriter.Flush()
}
