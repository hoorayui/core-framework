// generated by cap at 2019-12-18 19:24:24.4527213 +0800 CST m=+0.009989801, DO NOT EDIT.
package test

import (
	"fmt"

	"framework/pkg/cap/msg/i18n"
	"golang.org/x/text/language"
)

var test_lr test_langRegistry

func initTrRegistry() {
	if test_lr == nil {
		test_lr = test_langRegistry{}
	}
}

func (l test_langRegistry) addTranslation(id i18n.TrID, lang language.Tag, translation string) error {
	if ts, ok := test_lr[id]; ok {
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
	i18n.AddTrCandidates(&test_tr{})
}

// map[langID]map[langCode]translation
type test_langRegistry map[i18n.TrID]map[language.Tag]string

type test_tr struct{}

func (t *test_tr) Translate(lang language.Tag, id i18n.TrID, args ...interface{}) (string, error) {
	if ts, ok := test_lr[id]; ok {
		if t, ok := ts[lang]; ok {
			return fmt.Sprintf(t, args...), nil
		}
	}
	return "", i18n.ErrNoTranslateCandidate
}
