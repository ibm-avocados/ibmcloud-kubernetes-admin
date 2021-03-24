package notifier

import "encoding/json"

type Github struct {
}

func (g *Github) Send() error {
	return nil
}

func NewGithub(content []byte) (*Github, error) {
	var github Github
	if err := json.Unmarshal(content, &github); err != nil {
		return nil, err
	}
	return &github, nil
}
