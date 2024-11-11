package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/bellh14/DesignManager/pkg/utils/log"
)

type PayloadJson struct {
	UserName string `json:"username"`
	Content  string `json:"content"`
}

type FileInfo struct { // FileNum=@FileName
	FieldName string // file1, file2, etc
	FileName  string
	FilePath  string
}

type DiscordHook struct {
	PayloadJson PayloadJson
	Files       []FileInfo
	WebhookURL  string
	ThreadID    string
	Logger      log.Logger
}

func NewDiscordHook(
	payloadJson PayloadJson,
	fileInfo []FileInfo,
	webHookURL string,
	threadID string,
	logger log.Logger,
) *DiscordHook {
	return &DiscordHook{
		PayloadJson: payloadJson,
		Files:       fileInfo,
		WebhookURL:  webHookURL,
		ThreadID:    threadID,
		Logger:      logger,
	}
}

func (dis *DiscordHook) CallWebHook() {
	url, err := url.Parse(dis.WebhookURL)
	if err != nil {
		dis.Logger.Error("Error Parsing url:", err)
		return
	}

	jsonData, err := json.Marshal(dis.PayloadJson)
	if err != nil {
		dis.Logger.Error("Failed to marshal json payload:", err)
		return
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	err = writer.WriteField("payload_json", string(jsonData))
	if err != nil {
		dis.Logger.Error("Error adding json field:", err)
		return
	}

	for _, fileInfo := range dis.Files {
		fullFileName := dis.Files[0].FilePath + "/" + fileInfo.FileName
		file, err := os.Open(fullFileName)
		if err != nil {
			dis.Logger.Error(fmt.Sprintf("Error opening file %s:", fileInfo.FileName), err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fileInfo.FieldName, filepath.Base(fileInfo.FileName))
		if err != nil {
			dis.Logger.Error("Error creating form file:", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			dis.Logger.Error("Error copying file content:", err)
		}
	}

	err = writer.Close()
	if err != nil {
		dis.Logger.Error("Error closing writer:", err)
	}

	if dis.ThreadID != "" {
		query := url.Query()
		query.Set("thread_id", dis.ThreadID)
		url.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("Post", url.String(), &requestBody)
	if err != nil {
		dis.Logger.Error("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "multipart/form-data")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		dis.Logger.Error("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	dis.Logger.Log(fmt.Sprintf("Response status: %v", resp.Status))
}
