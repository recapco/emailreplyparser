package emailreplyparser

import (
	"bufio"
	"regexp"
	"strings"
	"unicode"
)

// Parses an email body string and creates a new Email struct.
func Read(text string) (*Email, error) {
	e := &Email{}

	err := e.read(text)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// Parses an email body string and returns the parsed reply string.
func ParseReply(text string) (string, error) {
	e := Email{}

	err := e.read(text)
	if err != nil {
		return "", err
	}

	return e.VisibleText(), nil
}

// Email is a container of text Fragments.
type Email struct {
	Fragments []*Fragment

	fragment     *Fragment
	foundVisible bool
}

// Returns the combined text of all non-hidden fragments in the email.
func (e *Email) VisibleText() string {
	var fragmentStrings []string
	for _, fragment := range e.Fragments {
		if !fragment.Hidden {
			fragmentStrings = append(fragmentStrings, fragment.String())
		}
	}

	return strings.TrimRightFunc(strings.Join(fragmentStrings, "\n"), unicode.IsSpace)
}

var wroteRegexp = regexp.MustCompile(`(?ms)^On\s(.+?)wrote:$`)
var negativeWroteRegexp = regexp.MustCompile(`(?ms)^On.*On\s.+?wrote:$`)
var underscoreRegexp = regexp.MustCompile(`(?ms)([^\n])(\n_{7}_+)$`)

func (e *Email) read(text string) error {
	text = strings.Replace(text, "\r\n", "\n", -1)

	wroteMatches := wroteRegexp.FindAllString(text, -1)
	for _, wroteMatch := range wroteMatches {
		if negativeWroteRegexp.MatchString(wroteMatch) {
			continue
		}

		text = strings.Replace(text, wroteMatch, strings.Replace(wroteMatch, "\n", " ", -1), -1)
	}

	text = underscoreRegexp.ReplaceAllString(text, "$1\n$2")

	text = reverse(text)

	scanner := bufio.NewScanner(strings.NewReader(text)) // note: will error on lines longer than 65536
	for scanner.Scan() {
		e.scanLine(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	e.finishFragment()

	e.fragment = nil

	for i, j := 0, len(e.Fragments)-1; i < j; i, j = i+1, j-1 {
		e.Fragments[i], e.Fragments[j] = e.Fragments[j], e.Fragments[i]
	}

	return nil
}

var signatureRegexp = regexp.MustCompile(`(?ms)(--\s*$|__\s*$|\w-$)|(^(\w+\s*){1,3} ym morf tneS$)`)
var quotedRegexp = regexp.MustCompile(`(>+)$`)
var quoteHeaderRegexp = regexp.MustCompile(`^:etorw.*nO$`)

func (e *Email) scanLine(line string) {
	line = strings.TrimSuffix(line, "\n")

	if !signatureRegexp.MatchString(line) {
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
	}

	isQuoted := quotedRegexp.MatchString(line)

	if e.fragment != nil && line == "" && len(e.fragment.lines) > 0 {
		if signatureRegexp.MatchString(e.fragment.lines[len(e.fragment.lines)-1]) {
			e.fragment.Signature = true
			e.finishFragment()
		}
	}

	if e.fragment != nil && ((e.fragment.Quoted == isQuoted) ||
		(e.fragment.Quoted && (quoteHeaderRegexp.MatchString(line) || line == ""))) {
		e.fragment.lines = append(e.fragment.lines, line)
	} else {
		e.finishFragment()
		e.fragment = newFragment(isQuoted, line)
	}
}

func (e *Email) finishFragment() {
	if e.fragment != nil {
		e.fragment.finish()

		if !e.foundVisible {
			if e.fragment.Quoted || e.fragment.Signature || strings.TrimSpace(e.fragment.String()) == "" {
				e.fragment.Hidden = true
			} else {
				e.foundVisible = true
			}
		}

		e.Fragments = append(e.Fragments, e.fragment)
		e.fragment = nil
	}
}

// Fragment is a part of an email which can be marked as quoted, hidden, or as part of a signature.
type Fragment struct {
	Quoted    bool
	Signature bool
	Hidden    bool

	lines   []string
	content string
}

func newFragment(quoted bool, firstLine string) *Fragment {
	return &Fragment{
		lines:  []string{firstLine},
		Quoted: quoted,
	}
}

func (f *Fragment) finish() {
	f.content = strings.Join(f.lines, "\n")
	f.lines = nil

	f.content = reverse(f.content)
}

// Returns the text content of a Fragment.
func (f *Fragment) String() string {
	return f.content
}

func reverse(input string) string { // note: this will not work well with combining characters
	reversed := []rune(input)

	for i, j := 0, len(reversed)-1; i < len(reversed)/2; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return string(reversed)
}
