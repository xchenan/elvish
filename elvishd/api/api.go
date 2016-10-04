package api

type Request struct {
	Ping    bool
	ListDir bool
	// Quit    bool
}

type Response struct {
	Fatal  string `json:",omitempty"`
	Error  string `json:",omitempty"`
	Number int
}
