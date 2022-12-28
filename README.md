# Gibberish

[![Go Reference](https://pkg.go.dev/badge/github.com/prophittcorey/gibberish.svg)](https://pkg.go.dev/github.com/prophittcorey/gibberish)

A golang package and command line tool for the analysis and identification of
gibberish text.

## Package Usage

Training and saving a classifier.

```golang
import "github.com/prophittcorey/gibberish"

classifier := New()

/* train */

classifier.Train(strings.NewReader(`
This is some example text. You should train on a lot of textual data for
whatever your use case is; in as many target languages as you need, etc.
`))

/* save */

if err := classifier.SaveFile(os.TempDir() + "/gibberish.classifier"); err != nil {
  t.Fatalf("failed to write classifier file; %s", err)
}
```

Loading and using a classifier.

```bash
import "github.com/prophittcorey/gibberish"

classifier := New()

/* load */

if err := classifier.LoadFile(os.TempDir() + "/gibberish.classifier"); err != nil {
  t.Fatalf("failed to load classifier file; %s", err)
}

/* test */

hasGibberish, probability := classifier.Predict("Joey")

fmt.Printf("Has gibberish %v w/%.2f probability", hasGibberish, probability)
```

## Tool Usage

The `gibberish` tool can be used to create and test classifiers.

Installing the tool can be done through `go` tool.

```bash
$ go install github.com/prophittcorey/gibberish/cmd/gibberish@latest
```

First, we need some data. The data you will require depends on your intended
use of the classifier. For most people, any large text data in your target
languages will work.

For our test case, let's use an English novel (Moby Dick).

```bash
$ wget -O /tmp/moby-dick.txt https://www.gutenberg.org/files/15/15-0.txt
```

Training a classifier and testing it out is easy.

```bash
$ gibberish --train "/tmp/moby-dick.txt" --classifier /tmp/text.classifier
$ gibberish --classifier /tmp/text.classifier --check "This looks like a good sentence." # Gibberish? False (99.65%)
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
