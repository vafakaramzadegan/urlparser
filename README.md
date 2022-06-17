# urlparser

urlparser extracts scheme, subdomain, domain, port, suffix, and path from a URL accurately using Public Suffix List (PSL).
urlparser is inspired by [tldextract](https://pypi.org/project/tldextract/) and [tldextract](https://github.com/joeguo/tldextract) packages, with modifications.

## Installation
Simply run:
```
$ go get github.com/vafakaramzadegan/urlparser
```

## Usage
```go
package main

import (
	"fmt"

	"github.com/vafakaramzadegan/urlparser"
)

func main() {
	parser, _ := urlparser.Setup()

	urls := []string{
		"https://google.com/mail",
		"ftp://192.168.1.23:24/foo/bar/somefile.sql",
		"http://foo.bar.blahsite.com/index.html",
	}
	for _, url := range urls {
		fmt.Println(parser.Parse(url))
	}
}
```

The output would be:
```go
&Result{Scheme: https, Subdomain: , Domain: google, Port: , Suffix: com, Path: /mail}
&Result{Scheme: ftp, Subdomain: , Domain: 192.168.1.23, Port: 24, Suffix: , Path: /foo/bar/somefile.sql}
&Result{Scheme: http, Subdomain: foo.bar, Domain: blahsite, Port: , Suffix: com, Path: /index.html}
```
