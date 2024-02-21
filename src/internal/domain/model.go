package domain

type AdminInfo struct {
	ActiveParsers   []*ParserSettings
	AutoGrantLimit  int
	DefaultInterval int
}
