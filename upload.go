package gdutils

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
)

type Token oauth2.Token

func (a *Token) SaveTo(path string) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(a)
	if err != nil {
		log.Fatalf("Unable to decode oauth token: %v", err)
	}
}

type ServerClient struct {
	Server   *drive.Service
	FoldType string
}

func NewServerClient(Server *drive.Service) ServerClient {
	ServerClient := ServerClient{
		Server:   Server,
		FoldType: "application/vnd.google-apps.folder",
	}

	return ServerClient
}

func GetConfig(Credentials []byte) *oauth2.Config {
	// If modifying these scopes, delete your previously saved token-old.json.
	config, err := google.ConfigFromJSON(Credentials, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config
}

// GetTokenFromBytes 从bytes获取Token
func GetTokenFromBytes(data []byte) (c *Token, e error) {
	e = json.Unmarshal(data, c)
	return
}

// GetTokenFromFile 从文件中获取Token
func GetTokenFromFile(file string) (*Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return (*Token)(tok), err
}

// GetTokenFromWeb 网络请求token
// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config, ProcessAuthURl func(AuthUrl string) (AuthCode string)) *Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	//此处为输出验证文件回调函数与获取信息回调函数
	authCode := ProcessAuthURl(authURL)

	tok, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return (*Token)(tok)
}

// getClient 获取客户端
func getClient(config *oauth2.Config, tok *Token) *http.Client {
	return config.Client(context.Background(), (*oauth2.Token)(tok))
}

// GetService 获取服务器
// 进一步封装,优化结构
func GetService(config *oauth2.Config, tok *Token) *drive.Service {
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(getClient(config, tok)))

	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return srv
}

//获取为共享drive的id
//不包含自己的drive

func (s *ServerClient) GetDriveList(PageSize int64) (*drive.DriveList, error) {
	return s.Server.Drives.List().PageSize(PageSize).Do()
}

// GetFiles 获取文件列表
func (s *ServerClient) GetFiles(DriveId string) ([]*drive.File, error) {
	FileList, err := s.Server.Files.List().Corpora("drive").IncludeItemsFromAllDrives(true).SupportsAllDrives(true).DriveId(DriveId).Do()
	if err != nil {

		return nil, err

	}

	return FileList.Files, err
}

// GetFolders 获取文件夹列表
func (s *ServerClient) GetFolders(DriveId string) ([]*drive.File, error) {
	FileList, err := s.Server.Files.List().
		Corpora("drive").Q("mimeType='application/vnd.google-apps.folder'").
		OrderBy("createdTime desc").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		DriveId(DriveId).Do()
	if err != nil {
		return nil, err
	}

	return FileList.Files, err
}

// Upload 上传文件
func (s *ServerClient) Upload(FileName string, FolderId string, Reader io.Reader) (*drive.File, error) {
	return s.Server.Files.Create(&drive.File{
		Name:    FileName,
		Parents: []string{FolderId}},
	).Media(Reader).SupportsAllDrives(true).Do()
}

func (s ServerClient) CreateFolder(FolderName string, DriveId string) (*drive.File, error) {
	return s.Server.Files.Create(&drive.File{
		Name:    FolderName,
		DriveId: DriveId, Parents: []string{DriveId},
		MimeType: s.FoldType,
	}).SupportsAllDrives(true).SupportsTeamDrives(true).Do()
}
