package types

type Credentials struct {
	Site     string `json:"site"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Metadata string `json:"metadata"`
}

type CardInfo struct {
	ID         string `json:"id,omitempty"`
	Number     string `json:"number"`
	Expiration string `json:"expiration"`
	CVV        string `json:"cvv"`
	Metadata   string `json:"metadata"`
}

type Note struct {
	Key      string `json:"key"`
	Text     string `json:"text"`
	Metadata string `json:"metadata"`
}

type Key struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}
