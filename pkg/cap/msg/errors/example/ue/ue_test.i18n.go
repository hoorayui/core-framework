// generated by cap from ue_test.toml at 2019-12-19 17:34:05.6629308 +0800 CST m=+0.041856201, DO NOT EDIT.
package ue

import "golang.org/x/text/language"

func init() {
	initTrRegistry()
	ue_lr.addTranslation(ERR_TEST_1, language.Make("zh-CN"), "错误一")
	ue_lr.addTranslation(ERR_TEST_1, language.Make("en-US"), "error one")
	ue_lr.addTranslation(ERR_TEST_2, language.Make("zh-CN"), "错误二 %d, %s")
	ue_lr.addTranslation(ERR_TEST_2, language.Make("en-US"), "Error two %d, %s")
}
