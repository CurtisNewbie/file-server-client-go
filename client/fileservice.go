package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/curtisnewbie/gocommon"
	"github.com/sirupsen/logrus"
)

const (
	DIR               FileType = "DIR"
	FILE              FileType = "FILE"
	FILE_SERVICE_NAME string   = "file-service"

	PROP_FILE_SERVICE_BASE_URL = "client.fileServiceUrl"
)

type FileType string

type ValidateFileKeyResp struct {
	gocommon.Resp
	Data bool `json:"data"`
}

type FileInfoResp struct {

	/** name of the file */
	Name string `json:"name"`

	/** file's uuid */
	Uuid string `json:"uuid"`

	/** size of file in bytes */
	SizeInBytes int64 `json:"sizeInBytes"`

	/** uploader id, i.e., user.id */
	UploaderId int `json:"uploaderId"`

	/** uploader name */
	UploaderName string `json:"uploaderName"`

	/** when the file is deleted */
	IsDeleted bool `json:"isDeleted"`

	/** file type: FILE, DIR */
	FileType FileType `json:"fileType"`

	/** parent file's uuid */
	ParentFile string `json:"parentFile"`
}

type GetFileInfoResp struct {
	gocommon.Resp
	Data *FileInfoResp `json:"data"`
}

type ListFilesInDirResp struct {
	gocommon.Resp
	// list of file key
	Data []string `json:"data"`
}

// List files in dir from file-service
func ListFilesInDir(fileKey string, limit int, page int) (*ListFilesInDirResp, error) {
	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/indir/list?fileKey=%s&limit=%d&page=%d", fileKey, limit, page))
	logrus.Infof("List files in dir, url: %s", url)
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	logrus.Infof("List files in dir, resp: %v", string(body))

	var resp ListFilesInDirResp
	if e = json.Unmarshal(body, &resp); e != nil {
		return nil, e
	}

	if resp.Resp.Error {
		return nil, gocommon.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Get file info from file-service
func GetFileInfo(fileKey string) (*GetFileInfoResp, error) {
	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/info?fileKey=%s", fileKey))
	logrus.Infof("Get file info, url: %s", url)
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	logrus.Infof("Get file info, resp: %v", string(body))

	var resp GetFileInfoResp
	if e = json.Unmarshal(body, &resp); e != nil {
		return nil, e
	}

	if resp.Resp.Error {
		return nil, gocommon.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Download file from file-service
func DownloadFile(fileKey string, absPath string) error {
	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/download?fileKey=%s", fileKey))
	logrus.Infof("Download file, url: %s, absPath: %s", url, absPath)

	out, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	logrus.Infof("Finished downloading file, url: %s", url)
	return nil
}

// Validate the file key, return true if it's valid else false
func ValidateFileKey(fileKey string, userId string) (bool, error) {

	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/owner/validation?fileKey=%s&userId=%s", fileKey, userId))
	logrus.Infof("Validate file key, url: %s", url)

	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		return false, e
	}
	logrus.Infof("Validate file key, key: %v resp: %v", fileKey, string(body))

	var resp ValidateFileKeyResp
	if e := json.Unmarshal(body, &resp); e != nil {
		return false, e
	}

	if resp.Resp.Error {
		return false, gocommon.NewWebErr(resp.Resp.Msg)
	}

	return resp.Data, nil
}

/* 
	Concatenate given relative url to base url, the relUrl may or may not start with "/"

	This func looks for prop:

		PROP_FILE_SERVICE_BASE_URL

*/
func BuildFileServiceUrl(relUrl string) string {
	if !strings.HasPrefix(relUrl, "/") {
		relUrl = "/" + relUrl
	}

	// if consul is enabled, try to look it up in the server list first
	if gocommon.IsConsulClientInitialized() {
		address, err := gocommon.ResolveServiceAddress(FILE_SERVICE_NAME)
		if err == nil && address != "" {
			return "http://" + address + relUrl
		}
		logrus.Infof("Unable to find service address from consul for '%s', trying to use the one in config json file", FILE_SERVICE_NAME)
	}

	baseUrl := gocommon.GetPropStr(PROP_FILE_SERVICE_BASE_URL)
	if baseUrl == "" {
		panic("Missing client.fileServiceUrl configuration, unable to resolve base url for file-service")
	}
	return baseUrl + relUrl
}
