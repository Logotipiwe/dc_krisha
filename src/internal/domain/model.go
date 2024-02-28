package domain

type AdminInfo struct {
	ActiveParsers      []*ParserSettings
	AllParsersSettings []*ParserSettings
	AutoGrantLimit     int
	DefaultInterval    int
}
