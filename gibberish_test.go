package gibberish

import (
	"strings"
	"testing"
)

var trainingdata = `
To Sherlock Holmes she is always the woman. I have seldom heard him mention her under any other name. In his eyes she eclipses and predominates the whole of her sex. It was not that he felt any emotion akin to love for Irene Adler. All emotions, and that one particularly, were abhorrent to his cold, precise but admirably balanced mind. He was, I take it, the most perfect reasoning and observing machine that the world has seen, but as a lover he would have placed himself in a false position. He never spoke of the softer passions, save with a gibe and a sneer. They were admirable things for the observer--excellent for drawing the veil from men's motives and actions. But for the trained reasoner to admit such intrusions into his own delicate and finely adjusted temperament was to introduce a distracting factor which might throw a doubt upon all his mental results. Grit in a sensitive instrument, or a crack in one of his own high-power lenses, would not be more disturbing than a strong emotion in a nature such as his. And yet there was but one woman to him, and that woman was the late Irene Adler, of dubious and questionable memory.

I had seen little of Holmes lately. My marriage had drifted us away from each other. My own complete happiness, and the home-centred interests which rise up around the man who first finds himself master of his own establishment, were sufficient to absorb all my attention, while Holmes, who loathed every form of society with his whole Bohemian soul, remained in our lodgings in Baker Street, buried among his old books, and alternating from week to week between cocaine and ambition, the drowsiness of the drug, and the fierce energy of his own keen nature. He was still, as ever, deeply attracted by the study of crime, and occupied his immense faculties and extraordinary powers of observation in following out those clues, and clearing up those mysteries which had been abandoned as hopeless by the official police. From time to time I heard some vague account of his doings: of his summons to Odessa in the case of the Trepoff murder, of his clearing up of the singular tragedy of the Atkinson brothers at Trincomalee, and finally of the mission which he had accomplished so delicately and successfully for the reigning family of Holland. Beyond these signs of his activity, however, which I merely shared with all the readers of the daily press, I knew little of my former friend and companion.

One night--it was on the twentieth of March, 1888--I was returning from a journey to a patient (for I had now returned to civil practice), when my way led me through Baker Street. As I passed the well-remembered door, which must always be associated in my mind with my wooing, and with the dark incidents of the Study in Scarlet, I was seized with a keen desire to see Holmes again, and to know how he was employing his extraordinary powers. His rooms were brilliantly lit, and, even as I looked up, I saw his tall, spare figure pass twice in a dark silhouette against the blind. He was pacing the room swiftly, eagerly, with his head sunk upon his chest and his hands clasped behind him. To me, who knew his every mood and habit, his attitude and manner told their own story. He was at work again. He had risen out of his drug-created dreams and was hot upon the scent of some new problem. I rang the bell and was shown up to the chamber which had formerly been in part my own.
`

var gooddata = `
This is a normal sentence.
Here is another sentence that is normal.
A man, a plan, a codebase.
`

var baddata = `
xxxxxxxxxxxxxxxxxxxxxx
This looks terrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrribleeee!!!!
lkqwlkdjqwldjqwldjqwldjqwk
A junk lkrejtlerjtlrejtlkerjtlkerjtl sentence.
zxcvwerjasc
This sentence lamsxamklm123121l2mmlmdlamdla is bad?
212pok12o3k12pk312
asxpaksxpakxpakxapksx
This sentence kmlsm1l2m12lwm12lmasdasldajslal is bad?
`

func TestClassifier(t *testing.T) {
	classifier := New()

	err := classifier.Train(strings.NewReader(trainingdata))

	if err != nil {
		t.Fatalf("failed to train data; %s", err)
	}

	err = classifier.Feed(strings.NewReader(gooddata), strings.NewReader(baddata))

	if err != nil {
		t.Fatalf("failed to feed classifier; %s", err)
	}

	//
	// The "good".
	//

	good := []string{
		"Sherlock holmes returned.",
		"This should be good.",
		"This sentence is completely good, at least-- it seems to be to me.",
	}

	for _, sentence := range good {
		if ok, prob := classifier.Gibberish(sentence); ok {
			t.Fatalf("falsely flagged good sentence as gibberish; %s at %f (%f)", sentence, prob, classifier.threshold)
		}
	}

	//
	// The "bad".
	//

	bad := []string{
		"zzzzzzzzz",
		"lwekqwleklqwjelqwjel",
		"21p3k12po3kp12kpk",
	}

	for _, sentence := range bad {
		if ok, prob := classifier.Gibberish(sentence); !ok {
			t.Fatalf("failed to detect a gibberish sentence; %s at %f (%f)", sentence, prob, classifier.threshold)
		}
	}
}

func TestClassifierFromFile(t *testing.T) {
	classifier := New()

	err := classifier.LoadFile("testfiles/classifier.gob.gz")

	if err != nil {
		t.Fatalf("failed to load classifier; %s", err)
	}

	err = classifier.Feed(strings.NewReader(gooddata), strings.NewReader(baddata))

	if err != nil {
		t.Fatalf("failed to feed classifier; %s", err)
	}

	//
	// The "good".
	//

	good := []string{
		"Sherlock holmes returned.",
		"This should be good.",
		"This sentence is completely good, at least-- it seems to be to me.",
	}

	for _, sentence := range good {
		if ok, prob := classifier.Gibberish(sentence); ok {
			t.Fatalf("falsely flagged good sentence as gibberish; %s at %f (%f)", sentence, prob, classifier.threshold)
		}
	}

	//
	// The "bad".
	//

	bad := []string{
		"zzzzzzzzz",
		"This sentence sdsadasdasdjhkahjdkadhkjsdh is bad?",
		"21p3k12po3kp12kpk",
	}

	for _, sentence := range bad {
		if ok, prob := classifier.Gibberish(sentence); !ok {
			t.Fatalf("failed to detect a gibberish sentence; %s at %f (%f)", sentence, prob, classifier.threshold)
		}
	}
}

func TestClassifierNormalize(t *testing.T) {
	classifier := New([]rune("abc "))

	if cleaned := classifier.normalize("_-dza dfdb ddc"); cleaned != "a b c" {
		t.Fatalf("failed to clean stringl got %s", cleaned)
	}
}
