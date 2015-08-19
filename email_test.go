package emailreplyparser

import (
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestReadsSimpleBody(t *testing.T) {
	e := email(t, "email_1_1")

	assert.Len(t, e.Fragments, 3)
	assert.Equal(t, []bool{false, false, false}, quoteds(e.Fragments))
	assert.Equal(t, []bool{false, true, true}, signatures(e.Fragments))
	assert.Equal(t, []bool{false, true, true}, hiddens(e.Fragments))
	assert.Equal(t, `Hi folks

What is the best way to clear a Riak bucket of all key, values after
running a test?
I am currently using the Java HTTP API.
`, e.Fragments[0].String())
	assert.Equal(t, "-Abhishek Kona\n\n", e.Fragments[1].String())
}

func TestReadsTopPost(t *testing.T) {
	e := email(t, "email_1_3")

	assert.Len(t, e.Fragments, 5)
	assert.Equal(t, []bool{false, false, true, false, false}, quoteds(e.Fragments))
	assert.Equal(t, []bool{false, true, true, true, true}, hiddens(e.Fragments))
	assert.Equal(t, []bool{false, true, false, false, true}, signatures(e.Fragments))
	assert.Regexp(t, regexp.MustCompile(`^Oh thanks.\n\nHaving`), e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`^-A`), e.Fragments[1].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^On [^\:]+\:`), e.Fragments[2].String())
	assert.Regexp(t, regexp.MustCompile(`^_`), e.Fragments[4].String())
}

func TestReadsBottomPost(t *testing.T) {
	e := email(t, "email_1_2")

	assert.Len(t, e.Fragments, 6)
	assert.Equal(t, []bool{false, true, false, true, false, false}, quoteds(e.Fragments))
	assert.Equal(t, []bool{false, false, false, false, false, true}, signatures(e.Fragments))
	assert.Equal(t, []bool{false, false, false, true, true, true}, hiddens(e.Fragments))
	assert.Equal(t, "Hi,", e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`^On [^\:]+\:`), e.Fragments[1].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^You can list`), e.Fragments[2].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^> `), e.Fragments[3].String())
	assert.Regexp(t, regexp.MustCompile(`^_`), e.Fragments[5].String())
}

func TestReadsInlineReplies(t *testing.T) {
	e := email(t, "email_1_8")

	assert.Len(t, e.Fragments, 7)
	assert.Equal(t, []bool{true, false, true, false, true, false, false}, quoteds(e.Fragments))
	assert.Equal(t, []bool{false, false, false, false, false, false, true}, signatures(e.Fragments))
	assert.Equal(t, []bool{false, false, false, false, true, true, true}, hiddens(e.Fragments))
	assert.Regexp(t, regexp.MustCompile(`^On [^\:]+\:`), e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^I will reply`), e.Fragments[1].String())
	assert.Regexp(t, regexp.MustCompile(`okay?`), e.Fragments[2].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^and under this.`), e.Fragments[3].String())
	assert.Regexp(t, regexp.MustCompile(`inline`), e.Fragments[4].String())
	assert.Equal(t, "\n", e.Fragments[5].String())
	assert.Equal(t, e.Fragments[6].String(), "--\nHey there, this is my signature\n")
}

func TestRecognizesDateStringAboveQuote(t *testing.T) {
	e := email(t, "email_1_4")

	assert.Regexp(t, regexp.MustCompile(`^Awesome`), e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^On`), e.Fragments[1].String())
	assert.Regexp(t, regexp.MustCompile(`Loader`), e.Fragments[1].String())
}

func TestAComplexBodyWithOnlyOneFragment(t *testing.T) {
	e := email(t, "email_1_5")

	assert.Len(t, e.Fragments, 1)
}

func TestReadsEmailWithCorrectSignature(t *testing.T) {
	e := email(t, "correct_sig")

	assert.Len(t, e.Fragments, 2)
	assert.Equal(t, []bool{false, false}, quoteds(e.Fragments))
	assert.Equal(t, []bool{false, true}, signatures(e.Fragments))
	assert.Equal(t, []bool{false, true}, hiddens(e.Fragments))
	assert.Regexp(t, regexp.MustCompile("^-- \nrick"), e.Fragments[1].String())
}

func TestDealsWithMultilineReplyHeaders(t *testing.T) {
	e := email(t, "email_1_6")

	assert.Regexp(t, regexp.MustCompile(`^I get`), e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^On`), e.Fragments[1].String())
	assert.Regexp(t, regexp.MustCompile(`Was this`), e.Fragments[1].String())
}

func TestDealsWithWindowsLineEndings(t *testing.T) {
	e := email(t, "email_1_7")

	assert.Regexp(t, regexp.MustCompile(`:\+1:`), e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^On`), e.Fragments[1].String())
	assert.Regexp(t, regexp.MustCompile(`Steps 0-2`), e.Fragments[1].String())
}

// no need for does not modify input string test, not using a pointer

func TestReturnsOnlyTheVisibleFragmentsAsAString(t *testing.T) {
	e := email(t, "email_2_1")

	var visibleStrings []string
	for _, fragment := range e.Fragments {
		if !fragment.Hidden {
			visibleStrings = append(visibleStrings, fragment.String())
		}
	}
	assert.Equal(t, strings.TrimRightFunc(strings.Join(visibleStrings, "\n"), unicode.IsSpace), e.VisibleText())
}

func TestParseOutJustTopForOutlookReply(t *testing.T) {
	b := emailBody(t, "email_2_1")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Outlook with a reply", reply)
}

func TestParseOutJustTopForOutlookWithReplyDirectlyAboveLine(t *testing.T) {
	b := emailBody(t, "email_2_2")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Outlook with a reply directly above line", reply)
}

func TestParseOutSentFromIPhone(t *testing.T) {
	b := emailBody(t, "email_iPhone")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Here is another email", reply)
}

func TestParseOutSentFromBlackberry(t *testing.T) {
	b := emailBody(t, "email_BlackBerry")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Here is another email", reply)
}

func TestParseOutSendFromMultiwordMobileDevice(t *testing.T) {
	b := emailBody(t, "email_multi_word_sent_from_my_mobile_device")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Here is another email", reply)
}

func TestDoNotParseOutSendFromInRegularSentence(t *testing.T) {
	b := emailBody(t, "email_sent_from_my_not_signature")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Here is another email\n\nSent from my desk, is much easier then my mobile phone.", reply)
}

func TestRetainsBullets(t *testing.T) {
	b := emailBody(t, "email_bullets")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test 2 this should list second\n\nand have spaces\n\nand retain this formatting\n\n\n   - how about bullets\n   - and another", reply)
}

func TestParseReply(t *testing.T) {
	b := emailBody(t, "email_1_2")

	reply, err := ParseReply(b)
	if err != nil {
		t.Error(err)
	}

	e, err := Read(b)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, e.VisibleText(), reply)
}

func TestOneIsNotOn(t *testing.T) {
	e := email(t, "email_one_is_not_on")

	assert.Regexp(t, regexp.MustCompile(`One outstanding question`), e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^On Oct 1, 2012`), e.Fragments[1].String())
}

func TestMultipleOn(t *testing.T) {
	e := email(t, "greedy_on")

	assert.Regexp(t, regexp.MustCompile(`(?m)^On your remote host`), e.Fragments[0].String())
	assert.Regexp(t, regexp.MustCompile(`(?m)^On 9 Jan 2014`), e.Fragments[1].String())
	assert.Equal(t, []bool{false, true, false}, quoteds(e.Fragments))
	assert.Equal(t, []bool{false, false, false}, signatures(e.Fragments))
	assert.Equal(t, []bool{false, true, true}, hiddens(e.Fragments))

}

func TestPathologicalEmails(t *testing.T) {
	start := time.Now()
	_ = email(t, "pathological")
	assert.WithinDuration(t, start, time.Now(), time.Second)
}

func TestDoesntRemoveSignatureDelimiterInMidLine(t *testing.T) {
	e := email(t, "email_sig_delimiter_in_middle_of_line")

	assert.Len(t, e.Fragments, 1)
}

const testEmailsPath = "test_emails/"

func emailBody(t *testing.T, name string) string {
	body, err := ioutil.ReadFile(testEmailsPath + name + ".txt")
	if err != nil {
		t.Fatal(err)
	}

	return string(body)
}

func email(t *testing.T, name string) *Email {
	email, err := Read(emailBody(t, name))
	if err != nil {
		t.Fatal(err)
	}

	return email
}

func quoteds(fragments []*Fragment) []bool {
	quoteds := make([]bool, len(fragments))
	for i, fragment := range fragments {
		quoteds[i] = fragment.Quoted
	}

	return quoteds
}

func signatures(fragments []*Fragment) []bool {
	signatures := make([]bool, len(fragments))
	for i, fragment := range fragments {
		signatures[i] = fragment.Signature
	}

	return signatures
}

func hiddens(fragments []*Fragment) []bool {
	hiddens := make([]bool, len(fragments))
	for i, fragment := range fragments {
		hiddens[i] = fragment.Hidden
	}

	return hiddens
}
