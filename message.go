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

func Message(code InternalCode, language string) string {
	if code == "" {
		return DefaultMessage
	}

	if _, exist := globalMessageMap[code]; !exist {
		return DefaultMessage
	}

	return globalMessageMap[code][language]
}
