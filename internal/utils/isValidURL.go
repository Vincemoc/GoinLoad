package utils

import (
	"regexp"
)

type TargetDomain string

const (
	Cyberdrop TargetDomain = "cyberdrop"
	Ososedki  TargetDomain = "ososedki"
	Bunkr  TargetDomain = "bunkr"
	// Add more domains as needed
)

var acceptedPatterns = map[TargetDomain]string{
	Cyberdrop: `^https://cyberdrop\.me/a/([a-zA-Z0-9_-]+)$`,
	Ososedki:  `^https://ososedki\.com/photos/-(\d+_\d+)$`,
	Bunkr: `^https://bunkr\.sk/a/([a-zA-Z0-9_-]+)$`,
	// https://bunkr.sk/a/39Bi34i8
	// Add more patterns here if needed
}

func IsValidURL(url string) (TargetDomain, bool) {

	for domain, pattern := range acceptedPatterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(url)
		if len(match) > 1 {
			return domain, true
		}
	}

	return "", false
}
