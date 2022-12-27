package gibberish

import (
	"bufio"
	"io"
	"strings"
)

var (
	DefaultCharset = []rune("abcdefghijklmnopqrstuvwxyz ")
)

type Classifier struct {
	counts    map[rune]map[rune]int
	runes     map[rune]struct{}
	threshold float64
}

func (c *Classifier) Train(r io.Reader) error {
	// k := len(c.runes)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		runes := []rune(c.normalize(strings.TrimSpace(scanner.Text())))

		for _, grams := range ngrams(2, runes) {
			a, b := grams[0], grams[1]

			c.counts[a][b]++
		}
	}

	// Normalize the counts?

	// Pick a threshold.

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
		for _, r := range runes {
			classifier.runes[r] = struct{}{}
		}
	}

	classifier.counts = map[rune]map[rune]int{}

	for r := range classifier.runes {
		classifier.counts[r] = map[rune]int{}
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
