// generated by cap at 2019-12-19 17:34:05.6260512 +0800 CST m=+0.004976601, DO NOT EDIT.
package ue

import (
	"fmt"

	"github.com/hoorayui/core-framework/pkg/cap/msg/i18n"
	"golang.org/x/text/language"
)

var ue_lr ue_langRegistry

func initTrRegistry() {
	if ue_lr == nil {
		ue_lr = ue_langRegistry{}
	}
}

func (l ue_langRegistry) addTranslation(id i18n.TrID, lang language.Tag, translation string) error {
	if ts, ok := ue_lr[id]; ok {
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
	i18n.AddTrCandidates(&ue_tr{})
}

// map[langID]map[langCode]translation
type ue_langRegistry map[i18n.TrID]map[language.Tag]string

type ue_tr struct{}

func (t *ue_tr) Translate(lang language.Tag, id i18n.TrID, args ...interface{}) (string, error) {
	if ts, ok := ue_lr[id]; ok {
		if t, ok := ts[lang]; ok {
			return fmt.Sprintf(t, args...), nil
		}
	}
	return "", i18n.ErrNoTranslateCandidate
}
