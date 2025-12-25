package middleware

import (
	"encoding/json"

	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func GinI18nLocalize() gin.HandlerFunc {
	return i18n.Localize(
		i18n.WithBundle(&i18n.BundleCfg{
			RootPath:         "./locales",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.English,
			UnmarshalFunc:    json.Unmarshal,
			FormatBundleFile: "json",
		}),
		i18n.WithGetLngHandle(
			func(context *gin.Context, defaultLang string) string {
				lang := context.Query("lang")
				if lang == "" {
					if cookie, err := context.Cookie("lang"); err == nil {
						lang = cookie
					} else {
						return defaultLang
					}
				}
				return lang
			},
		),
	)
}
