package types

import "errors"

type CreateNoteRequest struct {
	Key      string `json:"key"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
}

func (req *CreateNoteRequest) Validate() (*Note, error) {
	if req.Key == "" {
		return nil, errors.New("incorrect key")
	}
	return &Note{
		Key:      req.Key,
		Text:     req.Data,
		Metadata: req.Metadata,
	}, nil
}

type UpdateNoteRequest struct {
	ID       string `json:"id"`
	Key      string `json:"key"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
}

func (req *UpdateNoteRequest) Validate() (*Note, error) {
	if req.Key == "" || req.ID == "" {
		return nil, errors.New("incorrect request")
	}
	return &Note{
		Key:      req.Key,
		Text:     req.Data,
		Metadata: req.Metadata,
	}, nil
}

type CreateCardRequest struct {
	ID         string `json:"id,omitempty"`
	Number     string `json:"number"`
	Expiration string `json:"expiration"`
	CVV        string `json:"cvv"`
	Metadata   string `json:"metadata"`
}

func (req *CreateCardRequest) Validate() (*CardInfo, error) {
	if req.Number == "" || req.Expiration == "" || req.CVV == "" {
		return nil, errors.New("incorrect data")
	}

	return &CardInfo{
		ID:         req.ID,
		Number:     req.Number,
		Expiration: req.Expiration,
		CVV:        req.CVV,
		Metadata:   req.Metadata,
	}, nil
}

type CreateCredentialsRequest struct {
	Site     string `json:"site"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Metadata string `json:"metadata"`
}

func (req *CreateCredentialsRequest) Validate() (*Credentials, error) {
	if req.Site == "" || req.Login == "" {
		return nil, errors.New("incorrect key")
	}
	return &Credentials{
		Site:     req.Site,
		Login:    req.Login,
		Password: req.Password,
		Metadata: req.Metadata,
	}, nil
}

type UpdateCredentialsRequest struct {
	ID       string `json:"id"`
	Site     string `json:"site"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Metadata string `json:"metadata"`
}

func (req *UpdateCredentialsRequest) Validate() (*Credentials, error) {
	if req.Site == "" || req.Login == "" || req.ID == "" {
		return nil, errors.New("incorrect request")
	}
	return &Credentials{
		Site:     req.Site,
		Login:    req.Login,
		Password: req.Password,
		Metadata: req.Metadata,
	}, nil
}
