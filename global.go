package errx

var Global = New()

func SetGlobal(e *Errx) {
	Global = e
}
