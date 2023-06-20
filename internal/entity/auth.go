package entity

type Token string

func (t Token) String() string {
	return string(t)
}

func NewToken(value string) *Token {
	token := Token(value)
	return &token
}
