package client

import (
	"fmt"
	"testing"

	"github.com/curtisnewbie/gocommon"
	log "github.com/sirupsen/logrus"
)

func PreTest() {
	gocommon.LoadConfigFromFile(fmt.Sprintf("../app-conf-%v.json", "dev"))
}

func TestDownloadFile(t *testing.T) {
	PreTest()
	fileKey := "4b10b75b-501e-4010-b6e7-e4724835e210"
	err := DownloadFile(fileKey, fmt.Sprintf("/tmp/%s.png", fileKey))
	if err != nil {
		t.Error(err)
	}
}

func TestValidateFileKey(t *testing.T) {
	PreTest()

	fileKey := "4b10b75b-501e-4010-b6e7-e4724835e210"
	userId := "30"
	hasAccess, err := ValidateFileKey(fileKey, userId)
	if err != nil {
		t.Error(err)
		return
	}
	if !hasAccess {
		t.Errorf("User %s should have access to file %s", userId, fileKey)
	}

}

func TestGetNonExistingFileInfo(t *testing.T) {
	PreTest()

	fileKey := "non-existing-file-key"
	resp, err := GetFileInfo(fileKey)
	if err == nil {
		t.Error("Should have received error because the fileKey doesn't exist, but there is none")
		return
	}
	log.Infof("TestGetFileInfo Resp: %+v", resp)
}

func TestGetFileInfo(t *testing.T) {
	PreTest()

	fileKey := "4b10b75b-501e-4010-b6e7-e4724835e210"
	resp, err := GetFileInfo(fileKey)
	if err != nil {
		t.Errorf("File doesn't exist but it should, err: %v", err)
		return
	}
	if resp.Data == nil {
		t.Error("Resp doesn't contain data, and there should be")
		return
	}
	log.Infof("Normal Resp.Data: %+v", resp.Data)
}

func TestListFilesInDir(t *testing.T) {
	PreTest()

	fileKey := "5ddf49ca-dec9-4ecf-962d-47b0f3eab90c"
	resp, err := ListFilesInDir(fileKey, 100, 1)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.Data == nil {
		t.Error("Resp doesn't contain data")
		return
	}
	log.Infof("TestListFilesInDir Resp.Data: %+v", resp.Data)
}
