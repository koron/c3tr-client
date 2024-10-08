package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

const (
	EnglishToJapanese = "English to Japanese"
	JapaneseToEnglish = "Japanese to English"
)

var rxEtoJ = regexp.MustCompile(`(?i:e.*(?:to|2).*j)`)
var rxJtoE = regexp.MustCompile(`(?i:j.*(?:to|2).*e)`)

func regulateMode(mode string, content string) (string, error) {
	if strings.EqualFold(mode, EnglishToJapanese) || rxEtoJ.MatchString(mode) {
		return EnglishToJapanese, nil
	}
	if strings.EqualFold(mode, JapaneseToEnglish) || rxJtoE.MatchString(mode) {
		return JapaneseToEnglish, nil
	}
	if mode != "" && !strings.EqualFold(mode, "auto") {
		return "", fmt.Errorf("invalid mode: %s", mode)
	}

	var total, ascii int64
	for _, r := range content {
		if r >= ' ' {
			total++
		}
		if r >= ' ' && r < 128 {
			ascii++
		}
	}

	//log.Printf("ascii=%d total=%d\n", ascii, total)
	if ascii < 3*total/4 {
		return JapaneseToEnglish, nil
	}
	return EnglishToJapanese, nil
}

// Available styles.
// See also https://huggingface.co/webbigdata/C3TR-Adapter/discussions/1#669e6ef419d0f96d8a77128b
const (
	Casual               = "casual"
	Formal               = "formal"
	Technical            = "technical"
	Journalistic         = "journalistic"
	WebFiction           = "web-fiction"
	Business             = "business"
	NSFW                 = "nsfw"
	EducationalCasual    = "educational-casual"
	AcademicPresentation = "academic-presentation"
	Slang                = "slang"
	SNSCasual            = "sns-casual"
)

var ValidWritingStyles = []string{
	Casual,
	Formal,
	Technical,
	Journalistic,
	WebFiction,
	Business,
	NSFW,
	EducationalCasual,
	AcademicPresentation,
	Slang,
	SNSCasual,
}

func IsValidWritingStyles(s string) bool {
	for _, valid := range ValidWritingStyles {
		if s == valid {
			return true
		}
	}
	return false
}

type PromptParam struct {
	Mode         string
	WritingStyle string
	Styles       map[string]string
	Content      string
}

//go:embed prompt.tmpl
var promptTmplRaw string

var promptTmpl = template.Must(template.New("prompt").Parse(promptTmplRaw))

func (pp PromptParam) Generate() (string, error) {
	bb := &bytes.Buffer{}
	err := promptTmpl.Execute(bb, pp)
	if err != nil {
		return "", err
	}
	return bb.String(), nil
}
