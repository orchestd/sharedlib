package stringHelpers

func BackQuote(object string) string {
	return "`" + object + "`"
}

func GetValueFromMapByLang(nameUseCases map[string]string, lang string) string {
	return nameUseCases[lang]
}
