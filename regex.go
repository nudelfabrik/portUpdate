package portUpdate

import (
	"regexp"
)

var regexDate *regexp.Regexp
var regexAuthor *regexp.Regexp
var regexAffectsLine *regexp.Regexp
var regexAffects *regexp.Regexp
var regexDescr *regexp.Regexp

func RegexCompile() {
	regexDate = regexp.MustCompile(`(([0-9]{8})):\n`)
	regexAuthor = regexp.MustCompile(`AUTHORS?: (.*)\n`)
	regexAffectsLine = regexp.MustCompile(`AFFECTS:.*\n`)
	regexAffects = regexp.MustCompile(`[A-Za-z0-9-*]*/[A-Za-z0-9-*]*`)
	regexDescr = regexp.MustCompile(`^.*\n.*\n.*\n.*\n((?s).*)`)
}
