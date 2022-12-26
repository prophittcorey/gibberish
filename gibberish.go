package gibberish

import "io"

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
