package gdutils

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"os"
)

func NewConfig(Credentials []byte) (*Config, error) {
	// If modifying these scopes, delete your previously saved token-old.json.
	config, err := google.ConfigFromJSON(Credentials, drive.DriveScope)
	return (*Config)(config), err
}

func NewConfigFromFile(path string) (*Config, error) {
	d, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}
	return NewConfig(d)
}

type Config oauth2.Config

func (a *Config) raw() *oauth2.Config {
	return (*oauth2.Config)(a)
}

// GetTokenFromBytes 从bytes获取Token
func (a *Config) GetTokenFromBytes(data []byte) (*Token, error) {
	var t oauth2.Token
	return &Token{
		c: a.raw(),
		t: &t,
	}, json.Unmarshal(data, &t)
}

// GetTokenFromFile 从文件中获取Token
func (a *Config) GetTokenFromFile(path string) (*Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tok oauth2.Token
	return &Token{
		c: a.raw(),
		t: &tok,
	}, json.NewDecoder(f).Decode(&tok)
}

// GetTokenFromWeb 网络请求token
// Request a token from the web, then returns the retrieved token.
func (a *Config) GetTokenFromWeb(ProcessAuthURl func(AuthUrl string) (AuthCode string)) (*Token, error) {
	tok, err := a.raw().Exchange(context.TODO(), ProcessAuthURl(a.raw().AuthCodeURL("state-token", oauth2.AccessTypeOffline)))

	return &Token{
		c: a.raw(),
		t: tok,
	}, err
}
