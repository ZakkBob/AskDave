//------------------------------------//
// Processes robots.txt into a struct //
//------------------------------------//

package robots

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/ZakkBob/AskDave/gocommon/url"
)

type UrlValidator struct {
	AllowedPatterns    []*regexp.Regexp `json:"allowed_patterns"`
	DisallowedPatterns []*regexp.Regexp `json:"disallowed_patterns"`
}

func (validator *UrlValidator) ValidateUrl(u *url.Url) bool {
	path := u.PathString()
	return validator.ValidatePath(path)
}

func (validator *UrlValidator) AllowedStrings() []string {
	allowed := make([]string, 0)
	for _, regex := range validator.AllowedPatterns {
		allowed = append(allowed, regex.String())
	}
	return allowed
}

func (validator *UrlValidator) DisallowedStrings() []string {
	disallowed := make([]string, 0)
	for _, regex := range validator.DisallowedPatterns {
		disallowed = append(disallowed, regex.String())
	}
	return disallowed
}

func (validator *UrlValidator) ValidatePath(path string) bool {
	longestMatch := 0
	isValid := true
	for _, pattern := range validator.AllowedPatterns {
		indices := pattern.FindStringIndex(path)
		if indices == nil {
			continue
		}
		matchLength := indices[1] - indices[0]
		if matchLength > longestMatch {
			longestMatch = matchLength
			isValid = true
		}
	}
	for _, pattern := range validator.DisallowedPatterns {
		indices := pattern.FindStringIndex(path)
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

func FromStrings(allowedS []string, disallowedS []string) (*UrlValidator, error) {
	var allowedR []*regexp.Regexp
	var disallowedR []*regexp.Regexp

	for _, s := range allowedS {
		r, err := regexp.Compile(s)
		if err != nil {
			return nil, fmt.Errorf("unable to parse allowedStrings: %w", err)
		}
		allowedR = append(allowedR, r)
	}

	for _, s := range disallowedS {
		r, err := regexp.Compile(s)
		if err != nil {
			return nil, fmt.Errorf("unable to parse disallowedStrings: %w", err)
		}
		disallowedR = append(disallowedR, r)
	}

	return &UrlValidator{
		AllowedPatterns:    allowedR,
		DisallowedPatterns: disallowedR,
	}, nil
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
	pattern = regexp.QuoteMeta(pattern)               //escape special chars
	pattern = "^" + pattern                           // match start of string
	pattern = strings.ReplaceAll(pattern, `\*`, ".*") //wildcard
	pattern = strings.ReplaceAll(pattern, `\$`, "$")  //endline

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

func generateUrlValidator(directives string) UrlValidator {
	validator := UrlValidator{make([]*regexp.Regexp, 0), make([]*regexp.Regexp, 0)}

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

		if value == "" {
			continue
		}

		regex, err := convertToRegex(value)
		if err != nil {
			continue
		}

		switch name {
		case "disallow":
			validator.DisallowedPatterns = append(validator.DisallowedPatterns, regex)
		case "allow":
			validator.AllowedPatterns = append(validator.AllowedPatterns, regex)
		}
	}

	return validator
}

func Parse(content string) (validator UrlValidator, sitemapUrl string) {
	content = removeComments(content)
	blocks := extractUserAgentBlocks(content)
	directives := extractDavebotDirectives(blocks)
	validator = generateUrlValidator(directives)

	return validator, ""
}
