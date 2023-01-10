# Gibberish

[![Go Reference](https://pkg.go.dev/badge/github.com/prophittcorey/gibberish.svg)](https://pkg.go.dev/github.com/prophittcorey/gibberish)

A golang package and command line tool for the analysis and identification of
gibberish text.

## Package Usage

Training and saving a classifier.

```golang
package main

import (
        "log"
        "os"
        "strings"

        "github.com/prophittcorey/gibberish"
)

func main() {
        classifier := gibberish.New()

        /* train */

        classifier.Train(strings.NewReader(`You did, Doctor, but none the less you must come round to my view, for otherwise I shall keep on piling fact upon fact on you until your reason breaks down under them and acknowledges me to be right. Now, Mr. Jabez Wilson here has been good enough to call upon me this morning, and to begin a narrative which promises to be one of the most singular which I have listened to for some time. You have heard me remark that the strangest and most unique things are very often connected not with the larger but with the smaller crimes, and occasionally, indeed, where there is room for doubt whether any positive crime has been committed. As far as I have heard, it is impossible for me to say whether the present case is an instance of crime or not, but the course of events is certainly among the most singular that I have ever listened to. Perhaps, Mr. Wilson, you would have the great kindness to recommence your narrative. I ask you not merely because my friend Dr. Watson has not heard the opening part but also because the peculiar nature of the story makes me anxious to have every possible detail from your lips. As a rule, when I have heard some slight indication of the course of events, I am able to guide myself by the thousands of other similar cases which occur to my memory. In the present instance I am forced to admit that the facts are, to the best of my belief, unique. I trust that I am not more dense than my neighbours, but I was always oppressed with a sense of my own stupidity in my dealings with Sherlock Holmes. Here I had heard what he had heard, I had seen what he had seen, and yet from his words it was evident that he saw clearly not only what had happened but what was about to happen, while to me the whole business was still confused and grotesque. As I drove home to my house in Kensington I thought over it all, from the extraordinary story of the red-headed copier of the Encyclopaedia down to the visit to Saxe-Coburg Square, and the ominous words with which he had parted from me. What was this nocturnal expedition, and why should I go armed? Where were we going, and what were we to do? I had the hint from Holmes that this smooth-faced pawnbroker's assistant was a formidable man--a man who might play a deep game. I tried to puzzle it out, but gave it up in despair and set the matter aside until night should bring an explanation.`))

        classifier.Label(
                /* good examples */
                strings.NewReader(`
                  This is a good sentence.
                  Another good sentence.
                  Here is a longer line that is good. Should help adjust the model.
                `),
                /* bad examples */
                strings.NewReader(`
                  A line that looks good but has xml;asxmalksmdlm12m1l2m1lm.
                  sadkasdlasjdlasjdla
                  asdlkaldjalsdjlasjd
                `),
        )

        /* save */

        if err := classifier.SaveFile(os.TempDir() + "/gibberish.classifier"); err != nil {
                log.Fatalf("failed to write classifier file; %s\n", err)
        }
}
```

Loading and using a classifier.

```golang
package main

import (
        "fmt"
        "log"
        "os"

        "github.com/prophittcorey/gibberish"
)

func main() {
        classifier := gibberish.New()

        /* load */

        if err := classifier.LoadFile(os.TempDir() + "/gibberish.classifier"); err != nil {
                log.Fatalf("failed to load classifier file; %s\n", err)
        }

        /* test */

        result := classifier.Analyze("This is a big sentence. With mlksamalkmdalskdmlaksmdlkasmkdlas.")

        if result.IsGibberish {
                fmt.Println("Looks like you got some gibberish text.")
        } else {
                fmt.Println("Looks like you got some good text.")
        }
}
```

## Tool Usage

The `gibberish` tool can be used to create and test classifiers.

Installing the tool can be done through `go` tool.

```bash
go install github.com/prophittcorey/gibberish/cmd/gibberish@latest
```

First, we need some data. The data you will require depends on your intended
use of the classifier. For most people, any large text data in your target
languages will work.

For our test case, let's use an English novel (Moby Dick).

```bash
wget -O /tmp/moby-dick.txt https://www.gutenberg.org/files/15/15-0.txt
```

Training a classifier and testing it out is easy. You will need some examples
of good and bad text. You can generate them on your own (only a few example
lines are needed).

```bash
gibberish --train "/tmp/moby-dick.txt" --good "/tmp/good.txt" --bad "/tmp/bad.txt" --classifier /tmp/english.classifier
gibberish --classifier /tmp/english.classifier --check "This looks like a good sentence." # $ => Good (205.11% / 68.58%)
```

A repl is included with the tool to enable quicker testing of classifiers.

```bash
gibberish --classifier /tmp/english.classifier --repl
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
