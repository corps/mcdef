package mcdef

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindTerms(t *testing.T) {
	actualTerms, actualText := FindTerms(testText, JapaneseWordSplitter)

	assert.Equal(t, len(expectedTerms), len(actualTerms))
	assert.Equal(t, expectedText, actualText, "Expected non definition text was not equal to actual!")

	for _, expectedTerm := range expectedTerms {
		actualTerm := findTerm(expectedTerm, actualTerms)
		if actualTerm == nil {
			assert.Fail(t, fmt.Sprintf("Expected to find term with reference %s but did not!", expectedTerm.Reference))
			continue
		}

		assert.Equal(t, expectedTerm.Text, actualTerm.Text)
		assert.Equal(t, expectedTerm.Definition, actualTerm.Definition)
		assert.Equal(t, len(expectedTerm.Splits), len(actualTerm.Splits))

		actualSplitSet := make(map[string]bool)
		for _, split := range actualTerm.Splits {
			actualSplitSet[split] = true
		}

		for _, split := range expectedTerm.Splits {
			if !actualSplitSet[split] {
				assert.Fail(t, fmt.Sprintf("Expected to find %s split for term reference %s!\n", split, expectedTerm.Reference))
			}
		}
	}
}

func TestFindTermsIssue(t *testing.T) {
	actualTerms, _ := FindTerms(otherTestText, JapaneseWordSplitter)

	assert.Equal(t, []Term{Term{}, Term{}}, actualTerms)
}

func findTerm(expectedTerm Term, actualTerms []Term) *Term {
	for _, actualTerm := range actualTerms {
		if actualTerm.Reference == expectedTerm.Reference {
			return &actualTerm
		}
	}
	return nil
}

var expectedTerms = []Term{
	Term{Text: "伝えられ", Reference: "c", Definition: "伝えられ \u003d\u003e　伝える", Splits: []string{"られ", "伝", "え"}},
	Term{Text: "マヘンドラパルバタ", Reference: "b", Definition: "マヘンドラパルバタ", Splits: []string{"マヘ", "ラパ", "タ", "ルバ", "ンド"}},
	Term{Text: "アンコールワット", Reference: "アンコールワット",
		Definition: "《「寺院町」の意》アンコールにある石造寺院遺跡。12世紀初め、クメール王朝スールヤバルマン2世の治下に建立。1992年、アンコールの他の遺跡とともに世界遺産（文化遺産）に登録された。",
		Splits:     []string{"ット", "ルワ", "アンコー"}},
	Term{Text: "密林", Reference: "密林", Definition: "みつ‐りん【密林】\n樹木などがすきまのないほど生い茂っている林。",
		Splits: []string{"密", "林"}},
	Term{Text: "碑文", Reference: "a", Definition: "ひ‐ぶん【碑文】 \n石碑に彫りつけた文章。碑銘。",
		Splits: []string{"碑", "文"}},
}

var otherTestText = `
 Here's some [completely](abc) different thing. 

 [abc]: /completely 
 Ho ho. 


 And something [else] for testing. 



 [else]: /else 
 For testing, yo.
 `

var testText = ` This text [will]: be included. 
	
[密林]: /密/林
みつ‐りん【密林】
樹木などがすきまのないほど生い茂っている林。
	
	
This will be included
	
[link refs]: won't result in terms, not referenced
* This list is con sumed by the above
	
* Bullet points should render
* Like you expect them
* *This* is still mark down.
	
[c]: /え
伝えられ =>　伝える
	
[a]: /
ひ‐ぶん【碑文】 
石碑に彫りつけた文章。碑銘。
	
筑波大などの国際研究チームは、カンボジアの北西 部の[密林][]に、古代クメール王朝が９世紀ごろに築いた[最初]の首都
「[マヘンドラパルバタ][b]」の遺跡を 見つけた、と発表した。１２世紀前半に建設された[アンコールワット]より
３００年ほど古い。これまで、[碑文][a][伝えられ][c]、寺院の一部が見つかっていたが、都市の全容はわかっていなかった。
	
[アンコールワット]: /アンコー/ルワ/
《「寺院町」の意》アンコールにある石造寺院遺跡。12世紀初め、クメール王朝スールヤバルマン2世の治下に建立。1992年、アンコールの他の遺跡とともに世界遺産（文化遺産）に登録された。 
	
[b]: /
マヘンドラパルバタ
	
[unreferenced resource]: /
	
`

var expectedText = ` This text [will]: be included.

This will be included

[link refs]: won't result in terms, not referenced
* This list is con sumed by the above

* Bullet points should render
* Like you expect them
* *This* is still mark down.

筑波大などの国際研究チームは、カンボジアの北西 部の[密林][]に、古代クメール王朝が９世紀ごろに築いた[最初]の首都
「[マヘンドラパルバタ][b]」の遺跡を 見つけた、と発表した。１２世紀前半に建設された[アンコールワット]より
３００年ほど古い。これまで、[碑文][a][伝えられ][c]、寺院の一部が見つかっていたが、都市の全容はわかっていなかった。

[unreferenced resource]: /

`

var termParsingXml = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">
<en-note
    style="word-wrap: break-word; -webkit-nbsp-mode: space; -webkit-line-break: after-white-space;">
  This text [will]: be included.
  <div>
    <br/>
  </div>
  <div>[密林]: /密/林</div>
  <div>みつ‐りん【密林】</div>
  <div>樹木などがすきまのないほど生い茂っている林。ジャングル。</div>
  <div>
    <br/>
  </div>
  <div>
    <br/>
  </div>
  <div>This will be included</div>
  <div>
    <br/>
  </div>
  <div>[link refs]: won't result in terms, not referenced</div>
  <div>* This list is con
    <u>sumed by the above</u>
  </div>
  <div>
    <u>
      <br/>
    </u>
  </div>
  <div>
    <u>* Bullet points should render</u>
  </div>
  <div>
    <u>* Like you expect them</u>
  </div>
  <div>
    <u>* *This* is still mark down.</u>
  </div>
  <div>
    <u>
      <br/>
    </u>
  </div>
  <div>
    <u>[c]: /え</u>
  </div>
  <div>
    <u>伝えられ =&gt;　伝える</u>
  </div>
  <div>
    <u>
      <br/>
    </u>
  </div>
  <div>
    <u>[a]: /</u>
  </div>
  <div>
    <u>ひ‐ぶん【<b>碑文</b>】
    </u>
  </div>
  <div>
    <u>石碑に彫りつけた文章。碑銘。</u>
  </div>
  <div>
    <u>
      <br/>
    </u>
  </div>
  <div><u>筑波大など</u>の国際研究チームは、カンボジアの北西
    <i>部の[密林][]に、古代クメール王朝が９世紀ごろに築いた[最初]の首都</i>
  </div>
  <div>
    <i>「[マヘンドラパルバタ][b]」の遺跡を
      <b>見つけた、と発表した。</b>
    </i>
    <b>１２世紀前半に建設された[アンコールワット]より</b>
  </div>
  <div>
    <b>３００年ほど古い。これまで、[碑文][a][伝えられ][c]、寺院の一部が見つかっていたが、都市の全容はわかっていなかった。</b>
  </div>
  <div>
    <b>
      <br/>
    </b>
  </div>
  <div>
    <b>[アンコールワット]: /アンコー/ルワ/</b>
  </div>
  <div><b>
    《「寺院町」の意》アンコールにある石造寺院遺跡。12世紀初め、クメール王朝スールヤバルマン2世の治下に建立。1992年、アンコールの他の遺跡とともに世界遺産（文化遺産）に登録さ</b>れた。
  </div>
  <div>
    <br/>
  </div>
  <div>[b]: /</div>
  <div>マヘンドラパルバタ</div>
  <div>
    <br/>
  </div>
  <div>[unreferenced resource]: /</div>
  <div>
    <br/>
  </div>
</en-note>
`
