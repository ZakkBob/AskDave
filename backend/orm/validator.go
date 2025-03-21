package orm

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ZakkBob/AskDave/gocommon/robots"
)

func ValidatorByUrl(s string) (robots.UrlValidator, error) {
	query := `SELECT allowed_patterns, disallowed_patterns 
		FROM robots 
		JOIN site 
		ON robots.site_id = site.id 
		AND site.url = $1;`

	var allowedPatterns []string
	var disallowedPatterns []string

	row := dbpool.QueryRow(context.Background(), query, s)
	err := row.Scan(&allowedPatterns, &disallowedPatterns)
	if err != nil {
		return robots.UrlValidator{}, fmt.Errorf("unable to get validator for url '%s': %w", s, err)
	}

	var validator robots.UrlValidator

	for _, allowed := range allowedPatterns {
		r, err := regexp.Compile(allowed)
		if err != nil {
			return robots.UrlValidator{}, fmt.Errorf("unable to parse regex '%s': %w", allowed, err)
		}
		validator.AllowedPatterns = append(validator.AllowedPatterns, r)
	}

	for _, disallowed := range disallowedPatterns {
		r, err := regexp.Compile(disallowed)
		if err != nil {
			return robots.UrlValidator{}, fmt.Errorf("unable to parse regex '%s': %w", disallowed, err)
		}
		validator.DisallowedPatterns = append(validator.DisallowedPatterns, r)
	}

	return validator, nil

}
