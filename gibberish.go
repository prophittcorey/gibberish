package gibberish

import "io"

/*

classifier := gibberish.New()

classifier.Train("some big text")

if ok, err := classifier.Check("something"); err == nil && ok {
	// Gibberish
}

*/

type Classifier struct {
}

func (c *Classifier) Train(r io.Reader) error {
	return nil
}

func (c *Classifier) Check(junk string) (bool, error) {
	return false, nil
}

// New creates a nrew classifier that is ready for use.
func New() *Classifier {
	return &Classifier{}
}
