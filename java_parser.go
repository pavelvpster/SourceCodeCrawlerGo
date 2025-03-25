package main

import (
	"regexp"
	"strings"
)

func GetJavaClass(s string) (bool, string, []string) {
	re := regexp.MustCompile(`\s*(public)?\s*(final|abstract)?\s*(class|interface|enum)\s*(?P<className>\w*)\s*(?P<extendsOrImplements>[\w,\s]*)?{`)
	if !re.MatchString(s) {
		return false, "", nil
	}

	var className, extendsOrImplements string
	match := re.FindStringSubmatch(s)
	for i, groupName := range re.SubexpNames() {
		if i != 0 {
			groupText := strings.TrimSpace(match[i])
			if len(groupText) == 0 {
				continue
			}

			switch groupName {
			case "className":
				className = groupText
			case "extendsOrImplements":
				extendsOrImplements = groupText
			}
		}
	}

	parentClassNames := make([]string, 0)
	if len(extendsOrImplements) > 0 {
		parentClassNames = append(parentClassNames, parseExtends(extendsOrImplements)...)
		parentClassNames = append(parentClassNames, parseImplements(extendsOrImplements)...)
	}

	return len(className) > 0, className, parentClassNames
}

func parseExtends(s string) []string {
	re := regexp.MustCompile(`\s*extends\s*(?P<baseClassName>\w*)`)
	if !re.MatchString(s) {
		return nil
	}

	parentClassNames := make([]string, 0)
	match := re.FindStringSubmatch(s)
	for i, groupName := range re.SubexpNames() {
		if i != 0 {
			groupText := strings.TrimSpace(match[i])
			if len(groupText) == 0 {
				continue
			}

			switch groupName {
			case "baseClassName":
				parentClassNames = append(parentClassNames, groupText)
			}
		}
	}

	return parentClassNames
}

func parseImplements(s string) []string {
	re := regexp.MustCompile(`implements\s*(?<firstInterface>\w*)(?P<otherInterfaces>(?:\s*,\s*\w+)+)?`)
	if !re.MatchString(s) {
		return nil
	}

	parentClassNames := make([]string, 0)
	match := re.FindStringSubmatch(s)
	for i, groupName := range re.SubexpNames() {
		if i != 0 {
			groupText := strings.TrimSpace(match[i])
			if len(groupText) == 0 {
				continue
			}

			switch groupName {
			case "firstInterface":
				parentClassNames = append(parentClassNames, groupText)
			case "otherInterfaces":
				for _, otherInterface := range strings.Split(groupText, ",") {
					t := strings.TrimSpace(otherInterface)
					if len(t) == 0 {
						continue
					}
					parentClassNames = append(parentClassNames, t)
				}
			}
		}
	}

	return parentClassNames
}
