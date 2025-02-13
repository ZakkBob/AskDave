//----------------------------------------------------------//
// Processes urls into a struct (for consistent formatting) //
//----------------------------------------------------------//

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Protocol int

const (
	UnspecifiedProtocol Protocol = iota
	HttpProtocol
	HttpsProtocol
)

type url struct {
	protocol      Protocol //optional (default Unspecified)
	subdomain     string   //optional (default none)
	domain        string   //required
	tld           string   //required
	port          int      //optional (default based on protocol)
	path          []string //optional (default empty)
	trailingSlash bool     //optional (default false)
}

var protocolStringMap = map[Protocol]string{
	UnspecifiedProtocol: "",
	HttpProtocol:        "http",
	HttpsProtocol:       "https",
}

var stringProtocolMap = map[string]Protocol{
	"":      UnspecifiedProtocol,
	"http":  HttpProtocol,
	"https": HttpsProtocol,
}

func protocolToString(p Protocol) string {
	return protocolStringMap[p]
}

func stringToProtocol(s string) Protocol {
	return stringProtocolMap[s]
}

func (u *url) String() string {
	s := ""

	if u.protocol != UnspecifiedProtocol {
		s += protocolToString(u.protocol) + "://"
	}

	if u.subdomain != "" {
		s += u.subdomain + "."
	}

	s += u.domain + "."
	s += u.tld

	if (u.port != 0) && !(u.protocol == HttpProtocol && u.port == 80) && !(u.protocol == HttpsProtocol && u.port == 443) {
		s += ":" + strconv.Itoa(u.port)
	}

	if len(u.path) != 0 {
		for _, p := range u.path {
			s += "/" + p
		}
	}

	if u.trailingSlash {
		s += "/"
	}

	return s
}

func parseAbsoluteUrl(s string) (url, error) {
	original_s := s
	var parsed url

	protocolRegex := regexp.MustCompile("^(https?):\\/\\/(.*)")
	protocolMatches := protocolRegex.FindStringSubmatch(s)

	if len(protocolMatches) > 2 {
		protocol := protocolMatches[1]
		parsed.protocol = stringToProtocol(protocol)
		s = protocolMatches[2] //remove the protocol to simplify future regex (this pattern repeats btw)
	}

	subdomainRegex := regexp.MustCompile("^([a-z0-9-]+?)\\.(.*?\\..*)")
	subdomainMatches := subdomainRegex.FindStringSubmatch(s)

	if len(subdomainMatches) > 2 {
		subdomain := subdomainMatches[1]
		parsed.subdomain = subdomain
		s = subdomainMatches[2] //see what i mean
	}

	domainRegex := regexp.MustCompile("^([a-z0-9-]+?)\\.(.*)")
	domainMatches := domainRegex.FindStringSubmatch(s)

	if len(domainMatches) > 2 {
		domain := domainMatches[1]
		parsed.domain = domain
		s = domainMatches[2]
	} else {
		return parsed, fmt.Errorf("url: '%s' does not contain a domain!", original_s)
	}

	tldRegex := regexp.MustCompile("^([a-z0-9-]+?)([/:].*)?$")
	tldMatches := tldRegex.FindStringSubmatch(s)

	if len(tldMatches) > 2 {
		tld := tldMatches[1]
		parsed.tld = tld
		s = tldMatches[2]
	} else {
		return parsed, fmt.Errorf("url: '%s' does not contain a tld!", original_s)
	}

	portRegex := regexp.MustCompile("^:(.+?)(\\/.*)?$")
	portMatches := portRegex.FindStringSubmatch(s)

	if len(portMatches) > 2 {
		port, err := strconv.Atoi(portMatches[1])
		if err != nil {
			return parsed, fmt.Errorf("malformed port: '%s' is not a valid port. (%s)", portMatches[1], err.Error())
		}
		parsed.port = port
		s = portMatches[2]
	}

	pathRegex := regexp.MustCompile("^\\/(.+?)(\\/)?$")
	pathMatches := pathRegex.FindStringSubmatch(s)

	if len(pathMatches) > 2 {
		path := strings.Split(pathMatches[1], "/")
		parsed.path = path

		s = pathMatches[2]
	}

	if s == "/" {
		parsed.trailingSlash = true
	}

	return parsed, nil
}

func normalisePath(p []string) []string {
	var n []string
	for _, segment := range p {
		switch segment {
		case ".":
			continue
		case "..":
			if len(n) != 0 {
				n = n[:len(n)-1]
			}
		default:
			n = append(n, segment)
		}
	}
	return n
}

func normaliseUrl(u url) url {
	u.path = normalisePath(u.path)
	length := len(u.path)
	if length == 0 {
		return u
	}

	if u.path[length-1] == "" {
		u.path = u.path[:length-1]
		u.trailingSlash = true
	}
	return u
}

func parseRelativeUrl(s string, base url) (url, error) {
	absUrl, err := parseAbsoluteUrl(s)

	if err == nil && absUrl.protocol != UnspecifiedProtocol {
		return absUrl, nil
	}

	regex := regexp.MustCompile("(\\.?\\.?\\/)?(.+)")
	matches := regex.FindStringSubmatch(s)

	if len(matches) < 3 {
		return url{}, fmt.Errorf("invalid relative url: %s", s)
	}

	if matches[1] == "./" || matches[1] == "../" || matches[1] == "" {
		path := strings.Split(matches[0], "/")
		base.path = normalisePath(append(base.path, path...))
		return normaliseUrl(base), nil
	} else if matches[1] == "/" {
		base.path = strings.Split(matches[2], "/")
		return normaliseUrl(base), nil
	}

	if err != nil {
		return absUrl, nil
	}

	return url{}, fmt.Errorf("invalid relative url: %s", s)
}
