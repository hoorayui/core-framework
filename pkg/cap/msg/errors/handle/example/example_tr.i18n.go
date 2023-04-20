// generated by cap at 2020-02-05 16:31:12.5613654 +0800 CST m=+0.012004101, DO NOT EDIT.
package main

import (
	"fmt"

	"framework/pkg/cap/msg/i18n"
	"golang.org/x/text/language"
)

var example_lr example_langRegistry

func initTrRegistry() {
	if example_lr == nil {
		example_lr = example_langRegistry{}
	}
}

func (l example_langRegistry) addTranslation(id i18n.TrID, lang language.Tag, translation string) error {
	if ts, ok := example_lr[id]; ok {
		if _, ok := ts[lang]; ok {
			return fmt.Errorf("dupplicate translation for [%s#%s]", id, lang.String())
		}
		l[id][lang] = translation
	} else {
		l[id] = map[language.Tag]string{
			lang: translation,
		}
	}
	return nil
}

func init() {
	i18n.AddTrCandidates(&example_tr{})
}

// map[langID]map[langCode]translation
type example_langRegistry map[i18n.TrID]map[language.Tag]string

type example_tr struct{}

func (t *example_tr) Translate(lang language.Tag, id i18n.TrID, args ...interface{}) (string, error) {
	if ts, ok := example_lr[id]; ok {
		if t, ok := ts[lang]; ok {
			return fmt.Sprintf(t, args...), nil
		}
	}
	return "", i18n.ErrNoTranslateCandidate
}
