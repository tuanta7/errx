package errx

type MultilingualMessage = map[string]string
type InternalCode = string

const DefaultMessage = "internal server error"

var globalMessageMap = make(map[InternalCode]MultilingualMessage)

func RegisterMessage(code InternalCode, language, message string) {
	if _, exist := globalMessageMap[code]; !exist {
		globalMessageMap[code] = make(MultilingualMessage)
	}

	globalMessageMap[code][language] = message
}

func Message(err *Error, language string) string {
	if err == nil {
		return DefaultMessage
	}

	errorCode := err.Code()

	messages, exists := globalMessageMap[errorCode]
	if !exists {
		return getErrorMessage(err)
	}

	if localized, exists := messages[language]; exists {
		return localized
	}

	return getErrorMessage(err)
}

func getErrorMessage(err *Error) string {
	if m := err.Message(); m != "" {
		return m
	}

	return DefaultMessage
}
