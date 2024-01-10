package utils

type ScreenType string
type Direction string

const (
	RAW_INPUT_SCREEN        ScreenType = "RIS"
	LIST_SCREEN             ScreenType = "LS"
	ROUTE_SCREEN            ScreenType = "RS"
	EXTERNAL_SCREEN         ScreenType = "ES"
	INBOX                   Direction  = "INBOX"
	OUTBOX                  Direction  = "OUTBOX"
	DEFAULT_SCREEN_LOCATION            = 0
)
