// Package gibberish enables you to train a model to determine if strings
// contain gibberish text.
package gibberish

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

var (
	// DefaultRuneSet are the basic runes used by the default classifier.
	// Basically, it avoids punctuation. If you have a special purpose like
	// classifying email addresses you may choose to include other runes like
	// '@', '.', etc.
	DefaultRuneSet = []rune("abcdefghijklmnopqrstuvwxyz ")
)

// Analysis is the results of Analyze. Contains information about the text that
// was checked.
type Analysis struct {
	Threshold   float64
	Probability float64
	IsGibberish bool
}

type Classifier struct {
	threshold float64
	counts    map[rune]map[rune]float64
	runes     map[rune]struct{}
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

	c.threshold = 50.0

	return nil
}

// Label good and bad data. The classifier will adjust its threshold
// for classifying data based on what you feed it.
func (c *Classifier) Label(good io.Reader, bad io.Reader) error {
	var maxbad float64

	scanner := bufio.NewScanner(bad)

	for scanner.Scan() {
		if p := c.avg([]rune(c.normalize(scanner.Text()))); p > maxbad {
			maxbad = p
		}
	}

	/* convervative, take the maximum bad probability as the bar */
	c.threshold = maxbad

	if maxbad == 0.0 {
		return fmt.Errorf("error: something isn't right; maxbad == 0.0")
	}

	return nil
}

func (c *Classifier) Analyze(junk string) *Analysis {
	probability := c.avg([]rune(c.normalize(junk)))

	return &Analysis{
		Probability: probability,
		Threshold:   c.threshold,
		IsGibberish: probability < c.threshold,
	}
}

// Load takes an io.Reader and decodes it. This can be used to read a
// classifier from a file, virtual file, bytes, etc.
func (c *Classifier) Load(reader io.Reader) error {
	gz, err := gzip.NewReader(reader)

	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(gz)

	if err := gz.Close(); err != nil {
		return err
	}

	return decoder.Decode(c)
}

// LoadFile wraps Load.
func (c *Classifier) LoadFile(file string) error {
	f, err := os.Open(file)

	if err != nil {
		return err
	}

	defer (func() {
		if err := f.Close(); err != nil {
			log.Printf("warning: failed to close file; %s", err)
		}
	})()

	return c.Load(f)
}

// Save serializes a classifier and writes it out to an io.Writer.
func (c Classifier) Save(writer io.Writer) error {
	gz := gzip.NewWriter(writer)

	if err := gob.NewEncoder(gz).Encode(&c); err != nil {
		return err
	}

	return gz.Close()
}

// SaveFile wraps Save.
func (c Classifier) SaveFile(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer (func() {
		if err := f.Close(); err != nil {
			log.Printf("warning: failed to close file; %s", err)
		}
	})()

	return c.Save(f)
}

func (c *Classifier) GobDecode(buf []byte) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(buf))

	if err := decoder.Decode(&c.counts); err != nil {
		return err
	}

	if err := decoder.Decode(&c.runes); err != nil {
		return err
	}

	return decoder.Decode(&c.threshold)
}

func (c Classifier) GobEncode() ([]byte, error) {
	w := &bytes.Buffer{}

	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(c.counts); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.runes); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.threshold); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// Normalize's a string for the given classifier. Removes any runes that are
// not part of the classifier's runeset.
func (c *Classifier) normalize(s string) string {
	var sb strings.Builder

	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
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
		runesets = append(runesets, DefaultRuneSet)
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
