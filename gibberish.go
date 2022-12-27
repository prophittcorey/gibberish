package gibberish

import (
	"bufio"
	"io"
	"log"
	"math"
	"strings"
)

var (
	DefaultCharset = []rune("abcdefghijklmnopqrstuvwxyz ")
)

type Classifier struct {
	counts    map[rune]map[rune]float64
	runes     map[rune]struct{}
	threshold float64
}

func (c *Classifier) Train(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		runes := []rune(c.normalize(strings.TrimSpace(scanner.Text())))

		for _, grams := range ngrams(2, runes) {
			a, b := grams[0], grams[1]

			c.counts[a][b]++
		}
	}

	// Normalize the counts to log probabilities.
	for a, transitions := range c.counts {
		s := float64(len(transitions))

		for b, counts := range transitions {
			c.counts[a][b] = math.Log(counts / s)
		}
	}

	// TODO: Pick a threshold automatically?

	return nil
}

func (c *Classifier) Check(junk string) (bool, error) {
	log.Println(c.avg([]rune(junk)), c.threshold)
	return c.avg([]rune(junk)) > c.threshold, nil
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

// avg is the average transition probability for a slice of runes.
func (c *Classifier) avg(runes []rune) float64 {
	log := 0.0
	counts := 0.0

	for _, grams := range ngrams(2, runes) {
		a, b := grams[0], grams[1]

		log += c.counts[a][b]

		counts++
	}

	if counts == 0.0 {
		counts = 1.0
	}

	return math.Exp(log / counts)
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

	classifier.counts = map[rune]map[rune]float64{}

	for r := range classifier.runes {
		classifier.counts[r] = map[rune]float64{}

		// Assume we've seen 10 of each rune. Acts as a kind of prior/smoothing.
		for k := range classifier.runes {
			classifier.counts[r][k] = 10.0
		}
	}

	// TODO: Should we try to determine this automatically?
	classifier.threshold = 0.85

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
