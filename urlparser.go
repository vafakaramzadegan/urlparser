package urlparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// the regular expressions used in parsing the url.
var (
	urlSchemeRegex    = regexp.MustCompile(`^(([a-zA-Z0-9]*?):{0,1}){0,1}\/\/`)
	ipv4Regex         = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)`)
	fullDomainRegex   = regexp.MustCompile(`^((\w+\.)+(\w+)+)`)
	portRegex         = regexp.MustCompile(`^:(\d+)`)
	publicSuffixRegex = regexp.MustCompile(`^[*.]*([a-z].+)`)
	pathRegex         = regexp.MustCompile(`^(\/.*)`)
)

// this list is used for proper suffix extraction.
const (
	publicSuffixListUrl = "https://publicsuffix.org/list/public_suffix_list.dat"
	cacheFn             = "publicSuffixList.dat"
)

// will contain all the entries found in public_suffix_list.dat.
type publicSuffixList struct {
	list []string
}

var suffixList publicSuffixList = publicSuffixList{}

type UrlParser struct{}

var url string

type Result struct {
	Scheme    string
	Subdomain string
	Domain    string
	Port      string
	Suffix    string
	Path      string
}

// download suffixList if not present, then load it into memory
func Setup() (*UrlParser, error) {
	err := suffixList.downloadSuffixList()
	if err != nil {
		return &UrlParser{}, err
	}
	err = suffixList.processSuffixList()
	if err != nil {
		return &UrlParser{}, err
	}

	return &UrlParser{}, nil
}

// download suffixList into a file
func (sl *publicSuffixList) downloadSuffixList() error {
	invalidFile := false

	f, err := os.OpenFile(cacheFn, 0, 0644)

	if errors.Is(err, os.ErrNotExist) {
		invalidFile = true
	} else {
		fi, err := f.Stat()
		f.Close()

		if err != nil {
			return err
		}

		if fi.Size() == 0 {
			invalidFile = true
		}
	}

	if invalidFile == true {
		out, err := os.Create(cacheFn)
		if err != nil {
			return err
		}

		resp, err := http.Get(publicSuffixListUrl)
		if err != nil || resp.StatusCode != 200 {
			return err
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

		out.Close()
		resp.Body.Close()
	}

	return nil
}

// load suffixList into memory and remove comments and other useless lines
func (sl *publicSuffixList) processSuffixList() error {
	file, err := os.Open("publicSuffixList.dat")
	if err != nil {
		return err
	}
	defer func() error {
		if err = file.Close(); err != nil {
			return err
		}
		return nil
	}()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		tld := publicSuffixRegex.FindStringSubmatch(scanner.Text())
		if tld != nil {
			sl.list = append(sl.list, tld[1])
		}
	}
	return nil
}

// try to extract the scheme from url
func (res *Result) evalScheme() {
	sch := urlSchemeRegex.FindStringSubmatch(url)
	if sch != nil {
		res.Scheme = sch[2]
		url = urlSchemeRegex.ReplaceAllString(url, "")
	}
}

// try to eval domain and subdomains if present
func (res *Result) evalDomain() {
	// check for IP address
	dom := ipv4Regex.FindStringSubmatch(url)
	if dom != nil {
		res.Domain = dom[0]
		return
	}

	dom = fullDomainRegex.FindStringSubmatch(url)
	if dom != nil {
		domList := strings.Split(dom[0], ".")

		var concat string
		for i := len(domList) - 1; i >= 0; i-- {
			concat = strings.TrimSuffix(
				domList[i]+"."+concat, ".",
			)
			for _, suffix := range suffixList.list {
				if concat == suffix {
					tmp := strings.Replace(dom[0], fmt.Sprintf(".%s", suffix), "", -1)
					if tmp == suffix {
						res.Domain = ""
						res.Subdomain = ""
						res.Suffix = suffix
						break
					}

					tmp2 := strings.Split(tmp, ".")
					res.Domain = tmp2[len(tmp2)-1]
					res.Subdomain = strings.Join(tmp2[:len(tmp2)-1], ".")
					res.Suffix = suffix

					break
				}
			}
		}
	}
}

// extract port from url
func (res *Result) evalPort() {
	url = fullDomainRegex.ReplaceAllString(url, "")
	p := portRegex.FindStringSubmatch(url)
	if p != nil {
		res.Port = p[1]
		url = portRegex.ReplaceAllString(url, "")
	}
}

// extract path from url
func (res *Result) evalPath() {
	p := pathRegex.FindStringSubmatch(url)
	if p != nil {
		res.Path = p[0]
	} else {
		res.Path = ""
	}
	url = ""
}

func (te *UrlParser) Parse(url_to_parse string) *Result {
	url = strings.TrimSpace(url_to_parse)

	res := Result{}
	res.evalScheme()
	res.evalDomain()
	res.evalPort()
	res.evalPath()

	return &res
}
