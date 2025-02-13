//------------------------------------//
// Processes robots.txt into a struct //
//------------------------------------//

package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type urlValidator struct {
	allowedPatterns    []*regexp.Regexp
	disallowedPatterns []*regexp.Regexp
}

func (validator *urlValidator) validate(url string) bool {
	longestMatch := 0
	isValid := true
	for _, pattern := range validator.allowedPatterns {
		indices := pattern.FindStringIndex(url)
		if indices == nil {
			continue
		}
		matchLength := indices[1] - indices[0]
		if matchLength > longestMatch {
			longestMatch = matchLength
			isValid = true
		}
	}
	for _, pattern := range validator.disallowedPatterns {
		indices := pattern.FindStringIndex(url)
		if indices == nil {
			continue
		}
		matchLength := indices[1] - indices[0]
		if matchLength > longestMatch {
			longestMatch = matchLength
			isValid = false
		}
	}
	return isValid
}

func extractUserAgentBlocks(content string) (blocks map[string]string) {
	blocks = make(map[string]string)

	agentRegex := regexp.MustCompile("(?si)User-Agent:")
	matches := agentRegex.Split(content, -1)

	blockRegex := regexp.MustCompile("(?si)(.*?)(?:\r\n|\r|\n)(.*)") //Dodgy regex v2 (yes i know i can do this without regex, but regex = high iq)

	for _, match := range matches {
		captures := blockRegex.FindStringSubmatch(match)

		if len(captures) < 3 {
			continue
		}

		userAgent := strings.ToLower(strings.TrimSpace(captures[1]))
		directives := strings.TrimSpace(captures[2])

		blocks[userAgent] = directives
	}

	return blocks
}

func extractDavebotDirectives(userAgentBlocks map[string]string) (directives string) { //Yeah yeah, i know i suck at naming things
	userAgents := map[string]int{"*": 0, "davebot": 1, "davebot/0.1": 2} // attempts to match with most specific user agent

	currentAgentAccuracy := -1

	for k, v := range userAgentBlocks {
		agentAccuracy, ok := userAgents[k]
		if !ok || agentAccuracy < currentAgentAccuracy {
			continue
		}

		directives = v
		currentAgentAccuracy = agentAccuracy
	}

	return directives
}

func convertToRegex(pattern string) (regex *regexp.Regexp, err error) {
	pattern = "^" + pattern                          // match start of string
	pattern = strings.ReplaceAll(pattern, "*", ".+") //wildcard

	regex, err = regexp.Compile(pattern)

	return regex, err
}

func removeComments(s string) string {
	noComments := ""
	scanner := bufio.NewScanner(strings.NewReader(s))

	if err := scanner.Err(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}

	for scanner.Scan() {
		stringSplit := strings.Split(scanner.Text(), "#")
		noComments += stringSplit[0] + "\n"
	}
	return noComments
}

func generateUrlValidator(directives string) urlValidator {
	validator := urlValidator{make([]*regexp.Regexp, 0), make([]*regexp.Regexp, 0)}

	scanner := bufio.NewScanner(strings.NewReader(directives))
	directiveRegex := regexp.MustCompile("(?i)(.*?):(.*)")

	if err := scanner.Err(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}

	for scanner.Scan() {
		directive := directiveRegex.FindStringSubmatch(scanner.Text())

		if len(directive) < 3 {
			continue
		}

		name := strings.ToLower(strings.TrimSpace(directive[1]))
		value := strings.TrimSpace(directive[2])

		regex, err := convertToRegex(value)
		if err != nil {
			continue
		}

		switch name {
		case "disallow":
			validator.disallowedPatterns = append(validator.disallowedPatterns, regex)
		case "allow":
			validator.allowedPatterns = append(validator.allowedPatterns, regex)
		}
	}

	return validator
}

func ProcessRobotsTxt(content string) (validator urlValidator, sitemapUrl string) {
	content = removeComments(content)
	blocks := extractUserAgentBlocks(content)
	directives := extractDavebotDirectives(blocks)
	validator = generateUrlValidator(directives)

	return validator, ""
}
