package sphinx

import (
	"net/http"

	//"github.com/guregu/kami"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

type langkey int

var (
	japaneseBase, _         = language.Japanese.Base()
	lkey            langkey = 0
)

// Detect Language from accept-language
func DetectLanguage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	lang := ""
	if tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language")); err == nil {
		if len(tags) > 0 {
			t := tags[0]
			base, _, _ := t.Raw() // base, sscript, region
			if base == japaneseBase {
				lang = "ja-JP"
			} else if t == language.AmericanEnglish {
				lang = "en-US"
			}
		}
	}
	ctx = LangWithContext(lang, ctx)
	return ctx
}

func LangWithContext(lang string, ctx context.Context) context.Context {
	return context.WithValue(ctx, lkey, lang)
}

func LangFromContext(ctx context.Context) string {
	lang, _ := ctx.Value(lkey).(string)
	return lang
}
