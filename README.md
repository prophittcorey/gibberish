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

classifier.Label(
    /* good examples */
    strings.NewReader(`
      This is a good sentence.
      Another good sentence.
    `),
    /* bad examples */
    strings.NewReader(`
      sadkasdlasjdlasjdla
      asdlkaldjalsdjlasjd
    `),
)

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

result := classifier.Analyze("Joey")

if result.IsGibberish {
  fmt.Println("Looks like you got some gibberish text.")
}
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

Training a classifier and testing it out is easy. You will need some examples
of good and bad text. You can generate them on your own (only a few example
lines are needed).

```bash
$ gibberish --train "/tmp/moby-dick.txt" --good "/tmp/good.txt" --bad "/tmp/bad.txt" --classifier /tmp/english.classifier
$ gibberish --classifier /tmp/english.classifier --check "This looks like a good sentence."
$ => Good (205.11% / 68.58%)
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
