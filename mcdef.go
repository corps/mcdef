/*
	mcdef provides a mechanism for calling out reference terms in a text and providing 'splits'
	on that text, most useful for defining study clozes.   This syntax is similar to markdown style
	links with references.  Terms can be split through manual definition in the text itself, but when
	it is not, a WordSplitter interface is used to provide a default split for the remaining parts of a
	term.  See the test for an example of this module's usage.
*/
package mcdef

import (
	"strings"
	"regexp"
)

var (
	termReferenceDefinitionPattern = regexp.MustCompile(`\[([^\]]+)\](\[([^\]]*)\])?`)
	paragraphSplittingPattern = regexp.MustCompile(`\s*\n\s*\n`)
	termDefinitionPattern = regexp.MustCompile(`(?s)\[(\S+)\]:\s*\/([^\s]*)\s*\n?(.*)`)
	JapaneseWordSplitter = NewRegexWordSplitter(regexp.MustCompile(`[一-龥]|[ぁ-ゖ゙゙゙ゝ-ヾ･-ﾝ]{1,2}|[^\s一-龥ぁ-ヾ　]{1,5}`))
)

type Term struct {
	Text string
	Reference string
	Definition string
	Splits []string
}

type WordSplitter interface {
	SplitWord(word string) []string
}

type WordSplitterFunc func(word string) []string

func (f WordSplitterFunc) SplitWord(word string) []string {
	return f(word)
}

func NewRegexWordSplitter(regex *regexp.Regexp) WordSplitterFunc {
	return func(word string) []string {
		result := regex.FindAllString(word, -1)
		return result
	}
}

func FindTerms(text string, splitter WordSplitter) (terms []Term, nonDefinitionText string) {
	paragraphs := paragraphSplittingPattern.Split(text, -1)

	termDefinitions := make([][]string, 0, 10)
	nonDefinitionParagraphs := make([]string, 0, 10)
	if paragraphs != nil {
		for _, paragraph := range paragraphs {
			match := termDefinitionPattern.FindStringSubmatch(paragraph)
			if match != nil {
				termDefinitions = append(termDefinitions, match)
			} else {
				nonDefinitionParagraphs = append(nonDefinitionParagraphs, paragraph)
			}
		}
	}

	nonDefinitionText = strings.Join(nonDefinitionParagraphs, "\n\n")
	return matchedTermDefinitions(termTexts(nonDefinitionText), termDefinitions, splitter), nonDefinitionText
}

func matchedTermDefinitions(termTexts map[string]string, termDefinitions [][]string, splitter WordSplitter) (terms []Term) {
	terms = make([]Term, 0, 10)
	seenReferences := make(map[string]bool) // Map only the first definition of a reference.

	for _, definition := range termDefinitions {
		reference := definition[1]
		manualSplits := splitBySeparator(definition[2])
		definition := strings.TrimSpace(definition[3])

		text, ok := termTexts[reference]
		if !ok {
			continue // Did not reference to it in the text.
		}

		splits := make(map[string]bool, 10)
		for k, v := range manualSplits {
			splits[k] = v
		}	

		for _, unsplitPiece := range unsplitParts(text, manualSplits) {
			if len(unsplitPiece) > 0 {
				for _, split := range splitter.SplitWord(unsplitPiece) {
					if len(split) > 0 {
						splits[split] = true
					}
				}
			}
		}

		_, exist := seenReferences[reference]
		if !exist {
			term := Term{Text: text, Definition: definition, Reference: reference}
			for splitText, _ := range splits {
				term.Splits = append(term.Splits, splitText)
			}
			seenReferences[reference] = true
			terms = append(terms, term)
		}
	}

	return terms
}

func unsplitParts(t string, manualSplits map[string]bool) []string {
	definedSplits := make([]string, 0, len(manualSplits))
	for manualSplit, _ := range manualSplits {
		definedSplits = append(definedSplits, regexp.QuoteMeta(manualSplit))
	}


	if len(definedSplits) > 0 {
		splitRegex := regexp.MustCompile(strings.Join(definedSplits, "|"))
		return splitRegex.Split(t, -1)
	} else {
		return []string {t}
	}
}

func splitBySeparator(t string) (result map[string]bool) {
	result = make(map[string]bool, 10)

	pieces := strings.Split(t, "/")
	for _, piece := range pieces {
		piece = strings.TrimSpace(piece)
		if len(piece) > 0 {
			result[piece] = true
		}
	}

	return result
}

func termTexts(text string) (references map[string]string) {
	references = make(map[string]string)
	referenceDefinitions := termReferenceDefinitionPattern.FindAllStringSubmatch(text, -1)
	if referenceDefinitions != nil {
		for _, referenceDefinition := range referenceDefinitions {
			termText := referenceDefinition[1]
			termName := referenceDefinition[3]
			if len(termName) == 0 {
				termName = termText
			}
			termName = strings.TrimSpace(termName)
			references[termName] = termText
		}
	}

	return references
}

