package urlparser

import (
	"log"
	"testing"
)

var (
	parser *UrlParser
	err    error
)

func init() {
	parser, err = Setup()

	if err != nil {
		log.Fatal("Unable to setup package!")
	}
}

func TestParse(t *testing.T) {
	testCases := map[string]*Result{
		"https://google.com/mail":                                        &Result{Scheme: "https", Subdomain: "", Domain: "google", Port: "", Suffix: "com", Path: "/mail"},
		"https://mail.google.com":                                        &Result{Scheme: "https", Subdomain: "mail", Domain: "google", Port: "", Suffix: "com", Path: ""},
		"http://www.somedomain.co.uk":                                    &Result{Scheme: "http", Subdomain: "www", Domain: "somedomain", Port: "", Suffix: "co.uk", Path: ""},
		"ftp://192.168.1.23:24/foo/bar/somefile.sql":                     &Result{Scheme: "ftp", Subdomain: "", Domain: "192.168.1.23", Port: "24", Suffix: "", Path: "/foo/bar/somefile.sql"},
		"https://some.subdomain.domain.us:81/images/test.jpg?w=720&q=90": &Result{Scheme: "https", Subdomain: "some.subdomain", Domain: "domain", Port: "81", Suffix: "us", Path: "/images/test.jpg?w=720&q=90"},
		"https://foo.bar.domain.noip.us/documents/catalogue.pdf":         &Result{Scheme: "https", Subdomain: "foo.bar", Domain: "domain", Port: "", Suffix: "noip.us", Path: "/documents/catalogue.pdf"},
		"blogspot.co.uk":                                                 &Result{Scheme: "", Subdomain: "", Domain: "", Port: "", Suffix: "blogspot.co.uk", Path: ""},
		"https://blah.blogspot.co.uk/some/path?custom_param=5&p2=hello":  &Result{Scheme: "https", Subdomain: "", Domain: "blah", Port: "", Suffix: "blogspot.co.uk", Path: "/some/path?custom_param=5&p2=hello"},
		"http://transporte.bo:8080":                                      &Result{Scheme: "http", Subdomain: "", Domain: "", Port: "8080", Suffix: "transporte.bo", Path: ""},
		"//sub.adomain.us/hello":                                         &Result{Scheme: "", Subdomain: "sub", Domain: "adomain", Port: "", Suffix: "us", Path: "/hello"},
		"http://foo.bar.blahsite.com/index.html":                         &Result{Scheme: "http", Subdomain: "foo.bar", Domain: "blahsite", Port: "", Suffix: "com", Path: "/index.html"},
	}

	for url, expected := range testCases {
		got := parser.Parse(url)
		if got.Scheme != expected.Scheme ||
			got.Subdomain != expected.Subdomain ||
			got.Domain != expected.Domain ||
			got.Port != expected.Port ||
			got.Suffix != expected.Suffix ||
			got.Path != expected.Path {
			t.Errorf("got %q, expected %q", got, expected)
		}
	}
}
