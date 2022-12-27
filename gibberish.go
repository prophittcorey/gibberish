package gibberish

import (
	"io"
	"strings"
)

var (
	DefaultCharset = []rune("abcdefghijklmnopqrstuvwxyz ")
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
func New(runesets ...[]rune) *Classifier {
	if len(runesets) == 0 {
		runesets = append(runesets, DefaultCharset)
	}

	classifier := &Classifier{runes: map[rune]struct{}{}}

	for _, runes := range runesets {
		for _, char := range runes {
			classifier.runes[char] = struct{}{}
		}
	}

	return classifier
}

// ngrams is a helper function that takes a slice of runes and returns the
// n-grams for the slice (into a slice of slice of runes).
func ngrams(n int, rs []rune) (runes [][]rune) {
	for i := 0; i < len(rs)-n+1; i++ {
		runes = append(runes, rs[i:i+n])
	}

	return runes
}
