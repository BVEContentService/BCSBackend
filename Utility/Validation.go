package Utility

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	REGEX_EMAIL = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	REGEX_GUID  = regexp.MustCompile("^[0-9a-f]{32}$") // "Cleaned" GUID - All lowercase and dashes removed
)

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func NormalizeGUID(guid string) string {
	guid = strings.Replace(guid, "-", "", -1)
	guid = strings.Replace(guid, "{", "", -1)
	guid = strings.Replace(guid, "}", "", -1)
	return strings.TrimSpace(strings.ToLower(guid))
}
