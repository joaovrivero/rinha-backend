package models

type Pessoa struct {
	ID         string   `json:"id"`
	Apelido    string   `json:"apelido"`
	Nome       string   `json:"nome"`
	Nascimento string   `json:"nascimento"`
	Stack      []string `json:"stack"`
}

type CreatePessoaRequest struct {
	Apelido    string   `json:"apelido" binding:"required,max=32"`
	Nome       string   `json:"nome" binding:"required,max=100"`
	Nascimento string   `json:"nascimento" binding:"required,datetime=2006-01-02"`
	Stack      []string `json:"stack"`
}
