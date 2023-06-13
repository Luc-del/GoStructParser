package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const file = "./input_struct.txt"

func main() {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "json") {
			parsed := parseLine(line).toYAML()
			fmt.Println(parsed)
		}
	}
}

type field struct {
	name     string
	dataType string
	tag      string
	format   string
}

func newField(name, dataType, tag string) field {
	t, f := convertType(dataType)
	return field{
		name:     name,
		dataType: t,
		tag:      tag,
		format:   f,
	}
}

func convertType(t string) (string, string) {
	switch t {
	case "string":
		return t, ""
	case "bool":
		return "boolean", ""
	case "int":
		return "integer", ""
	case "int64":
		return "integer", "int64"
	case "null.String":
		return "[ \"null\", \"string\" ]", ""
	case "time.Time":
		return "string", "date-time"
	}

	fmt.Println("no corresponding value for", t, "you might want to look it up")
	return t, ""
}

func (f field) toYAML() string {
	res := fmt.Sprintf("%s:\n\ttype: %s", f.tag, f.dataType)
	if f.format != "" {
		res = fmt.Sprintf("%s\n\tformat: %s", res, f.format)
	}
	return res
}

func parseLine(line string) field {
	line = strings.ReplaceAll(line, "\t", " ")
	line = regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")
	line = strings.TrimSpace(line)
	s := strings.Split(line, " ")
	return newField(s[0], s[1], parseTags(s[2]))
}

func parseTags(tags string) string {
	tags = strings.TrimPrefix(tags, "`json:\"")
	sp := strings.Split(tags, "\"")
	return sp[0]
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
