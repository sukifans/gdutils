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

type ServerClient struct {
	Server *drive.Service

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
func GetTokenFromBytes(data []byte) (c *oauth2.Token, e error) {
	e = json.Unmarshal(data, c)
	return
}

// GetTokenFromFile 从文件中获取Token
func GetTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// GetTokenFromWeb 网络请求token
// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config, ProcessAuthURl func(AuthUrl string), ReturnAuthCode func() string) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	//此处为输出验证文件回调函数
	ProcessAuthURl(authURL)

	var authCode string

	//此处为获取信息回调函数
	authCode = ReturnAuthCode()

	tok, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return tok
}

func SaveToken(path string, token *oauth2.Token) {

	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

}

// getClient 获取客户端
func getClient(config *oauth2.Config, tok *oauth2.Token) *http.Client {

	return config.Client(context.Background(), tok)

}

// GetService 获取服务器
// 进一步封装,优化结构
func GetService(config *oauth2.Config, tok *oauth2.Token) *drive.Service {

	client := getClient(config, tok)

	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return srv

}

//获取为共享drive的id
//不包含自己的drive

func (s *ServerClient) GetDriveList(PageSize int64) (*drive.DriveList, error) {

	List, err := s.Server.Drives.List().PageSize(PageSize).Do()

	return List, err

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

	FileList, err := s.Server.Files.List().Corpora("drive").Q("mimeType='application/vnd.google-apps.folder'").OrderBy("createdTime desc").IncludeItemsFromAllDrives(true).SupportsAllDrives(true).DriveId(DriveId).Do()

	var FolderList []*drive.File

	if err != nil {

		return nil, err

	}

	FolderList = FileList.Files

	return FolderList, err

}

// Upload 上传文件
func (s *ServerClient) Upload(FileName string, FolderId string, Reader io.Reader) (*drive.File, error) {

	FileMeta := &drive.File{Name: FileName, Parents: []string{FolderId}}

	FileInfo, err := s.Server.Files.Create(FileMeta).Media(Reader).SupportsAllDrives(true).Do()

	if err != nil {

		return nil, err

	}

	return FileInfo, err

}

func (s ServerClient) CreateFolder(FolderName string, DriveId string) (*drive.File, error) {

	FileMeta := &drive.File{Name: FolderName, DriveId: DriveId, Parents: []string{DriveId}, MimeType: s.FoldType}

	FolderInfo, err := s.Server.Files.Create(FileMeta).SupportsAllDrives(true).SupportsTeamDrives(true).Do()

	if err != nil {
		return nil, err
	}

	return FolderInfo, err

}
