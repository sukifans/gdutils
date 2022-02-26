package gdutils

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"os"
)

type Token struct {
	c *oauth2.Config
	t *oauth2.Token
}

func (a *Token) SaveTo(path string) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(a)
}

func (a *Token) NewService() (*ServerClient, error) {
	srv, e := drive.NewService(
		context.Background(),
		option.WithHTTPClient(a.c.Client(context.Background(), a.t)),
	)
	return (*ServerClient)(srv), e
}
