package gibberish

import (
	"io"
	"strings"
)

/*

classifier := gibberish.New()

classifier.Train(strings.NewReader("some big text"))

if ok, err := classifier.Check("something"); err == nil && ok {
	// Gibberish
}

*/

const (
	DefaultCharset = "abcdefghijklmnopqrstuvwxyz "
)

type Classifier struct {
	runes map[rune]struct{}
}

func (c *Classifier) Train(r io.Reader) error {
	return nil
}

func (c *Classifier) Check(junk string) (bool, error) {
	return false, nil
}

// Normalize's a string for the given classifier. Removes any runes that are
// not part of the classifier's runeset.
func (c *Classifier) normalize(s string) string {
	var sb strings.Builder

	for _, r := range s {
		if _, ok := c.runes[r]; ok {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// New creates a new classifier that is ready for use.
func New(runsets ...string) *Classifier {
	if len(runsets) == 0 {
		runsets = append(runsets, DefaultCharset)
	}

	classifier := &Classifier{runes: map[rune]struct{}{}}

	for _, runes := range runsets {
		for _, char := range runes {
			classifier.runes[char] = struct{}{}
		}
	}

	return classifier
}
