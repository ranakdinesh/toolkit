package toolkit

import (
	"errors"
	"regexp"
	"strings"
)

// CreateSlug is a function that takes a string and returns a slug. It removes all non-alphanumeric characters and replaces them with a dash. It also removes any leading or trailing dashes. It returns an error if the string is empty or if the slug is empty.
func (t *Tools) CreateSlug(s string) (string, error) {
	if s == "" {
		return "", errors.New("Empty string not Allowed")
	}
	var re = regexp.MustCompile("[^a-z\\d]+")

	slug := strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if slug == "" {
		return "", errors.New("Slug is empty")
	}
	return slug, nil
}
