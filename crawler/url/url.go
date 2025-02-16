//----------------------------------------------------------//
// Processes urls into a struct (for consistent formatting) //
//----------------------------------------------------------//

package url

import (
	"encoding/json"
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

type Url struct {
	Protocol      Protocol //optional (default Unspecified)
	Subdomain     string   //optional (default none)
	Domain        string   //required (unless ip is used)
	Tld           string   //required (unless ip is used)
	Port          int      //optional (default based on protocol)
	Path          []string //optional (default empty)
	TrailingSlash bool     //optional (default false)
}

func (u *Url) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *Url) UnmarshalJSON(data []byte) error {
	var s string
	json.Unmarshal(data, &s)
	newU, err := ParseAbs(s)
	if err != nil {
		return err
	}
	*u = newU
	return err
}

func (u *Url) IsFile() bool {
	fileRegex := regexp.MustCompile(`^.+\..+$`)
	pathLen := len(u.Path)
	if pathLen == 0 {
		return false
	}
	return fileRegex.MatchString(u.Path[pathLen-1])
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

func (u *Url) ProtocolString() string {
	return protocolToString(u.Protocol)
}

func (u *Url) PathString() string {
	s := ""
	for _, p := range u.Path {
		s += "/" + p
	}
	return s
}

func (u *Url) FQDN() string {
	s := ""

	if u.Subdomain != "" {
		s += u.Subdomain + "."
	}

	s += u.Domain + "."
	s += u.Tld

	return s
}

func (u *Url) String() string {
	s := ""

	if u.Protocol != UnspecifiedProtocol {
		s += protocolToString(u.Protocol) + "://"
	}

	if u.Subdomain != "" {
		s += u.Subdomain + "."
	}

	s += u.Domain + "."
	s += u.Tld

	if (u.Port != 0) && !(u.Protocol == HttpProtocol && u.Port == 80) && !(u.Protocol == HttpsProtocol && u.Port == 443) {
		s += ":" + strconv.Itoa(u.Port)
	}

	if len(u.Path) != 0 {
		s += u.PathString()
	}

	if u.TrailingSlash {
		s += "/"
	}

	return s
}

func ParseAbs(s string) (Url, error) {
	original_s := s
	var parsed Url

	protocolRegex := regexp.MustCompile(`^(https?):\/\/(.*)`)
	protocolMatches := protocolRegex.FindStringSubmatch(s)

	if len(protocolMatches) > 2 {
		protocol := protocolMatches[1]
		parsed.Protocol = stringToProtocol(protocol)
		s = protocolMatches[2] //remove the protocol to simplify future regex (this pattern repeats btw)
	}

	subdomainRegex := regexp.MustCompile(`^([a-z0-9-]+?)\.([a-z0-9-]*?\..*)`)
	subdomainMatches := subdomainRegex.FindStringSubmatch(s)

	if len(subdomainMatches) > 2 {
		subdomain := subdomainMatches[1]
		parsed.Subdomain = subdomain
		s = subdomainMatches[2] //see what i mean
	}

	domainRegex := regexp.MustCompile(`^([a-z0-9-]+?)\.(.*)`)
	domainMatches := domainRegex.FindStringSubmatch(s)

	if len(domainMatches) > 2 {
		domain := domainMatches[1]
		parsed.Domain = domain
		s = domainMatches[2]
	} else {
		return parsed, fmt.Errorf("url: '%s' does not contain a domain", original_s)
	}

	tldRegex := regexp.MustCompile(`^([a-z0-9-]+?)([/:].*)?$`)
	tldMatches := tldRegex.FindStringSubmatch(s)

	if len(tldMatches) > 2 {
		tld := tldMatches[1]
		parsed.Tld = tld
		s = tldMatches[2]
	} else {
		return parsed, fmt.Errorf("url: '%s' does not contain a tld", original_s)
	}

	portRegex := regexp.MustCompile(`^:(.+?)(\/.*)?$`)
	portMatches := portRegex.FindStringSubmatch(s)

	if len(portMatches) > 2 {
		port, err := strconv.Atoi(portMatches[1])
		if err != nil {
			return parsed, fmt.Errorf("malformed port: '%s' is not a valid port. (%s)", portMatches[1], err.Error())
		}
		parsed.Port = port
		s = portMatches[2]
	}

	pathRegex := regexp.MustCompile(`^\/(.+?)(\/)?$`)
	pathMatches := pathRegex.FindStringSubmatch(s)

	if len(pathMatches) > 2 {
		path := strings.Split(pathMatches[1], "/")
		parsed.Path = path
		s = pathMatches[2]
	}

	if s == "/" {
		parsed.TrailingSlash = true
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

func normalise(u Url) Url {
	u.Path = normalisePath(u.Path)
	length := len(u.Path)
	if length == 0 {
		return u
	}

	if u.Path[length-1] == "" {
		u.Path = u.Path[:length-1]
		u.TrailingSlash = true
	}
	return u
}

func ParseRel(s string, base Url) (Url, error) {
	absUrl, err := ParseAbs(s)

	if err == nil && absUrl.Protocol != UnspecifiedProtocol {
		return absUrl, nil
	}

	regex := regexp.MustCompile(`(\.?\.?\/)?(.+)`)
	matches := regex.FindStringSubmatch(s)

	if len(matches) < 3 {
		return Url{}, fmt.Errorf("invalid relative url: %s", s)
	}

	if matches[1] == "./" || matches[1] == "../" || matches[1] == "" {
		path := strings.Split(matches[0], "/")
		if base.IsFile() {
			base.Path = base.Path[:len(base.Path)-1]
		}
		base.Path = normalisePath(append(base.Path, path...))
		return normalise(base), nil
	} else if matches[1] == "/" {
		base.Path = strings.Split(matches[2], "/")
		return normalise(base), nil
	}

	if err != nil {
		return absUrl, nil
	}

	return Url{}, fmt.Errorf("invalid relative url: %s", s)
}
