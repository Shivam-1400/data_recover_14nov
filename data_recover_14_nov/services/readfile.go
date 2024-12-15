package services

import (
	"bufio"
	"context"
	"data_recover_14_nov/globals"
	"data_recover_14_nov/model"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func FileRead(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	folderPath := globals.ApplicationConfig.Application.ReadPath

	files, err := filepath.Glob(filepath.Join(folderPath, "*"))
	if err != nil {
		fmt.Printf("Error reading files: %v\n", err)
		return
	}

	for _, file := range files {
		select {
		case <-ctx.Done():
			fmt.Println("File processing stopped due to shutdown signal.")
			return
		default:
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				processFile(file)
			}(file)
		}
	}
}

func processFile(filePath string) {
	ctx := context.Background()
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("File processing stopped due to shutdown signal.")
			return
		default:
			line := scanner.Text()
			parseLogLine(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
	}
}

func parseLogLine(line string) *model.BulksmsLogData {
	if !strings.Contains(line, "DM REQUEST IS") {
		return nil
	}

	parts := strings.Split(line, "&")
	logData := &model.BulksmsLogData{}

	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			continue
		}
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		a, err := url.QueryUnescape(value)
		if err != nil {
			log.Fatal(err)
		} else {
			value = a
		}

		switch key {
		case "appid":
			logData.AppID = value
		case "feedid":
			logData.FeedID = value
		case "entid":
			logData.EntID = value
		case "keyword":
			logData.Keyword = value
		case "orgTempalte":
			logData.OrgTemplate = value
		case "dotStarCnt":
			logData.DotStarCnt = value
		case "spaceFlag":
			logData.SpaceFlag = value
		case "specialCharFlag":
			logData.SpecialCharFlag = value
		case "customdomain":
			logData.CustomDomain = value
		case "token":
			logData.Token = value
		case "dmCheckStatus":
			logData.DMCheckStatus = value
		case "breqid":
			logData.BReqID = value
		case "btid":
			logData.BTID = value
		case "traicategoryid":
			logData.TraiCategoryID = value
		case "traimessagetype":
			logData.TraiMessageType = value
		case "traimessagemode":
			logData.TraiMessageMode = value
		case "bsms_intime":
			logData.BSMSInTime = value
		case "template_id":
			logData.TemplateID = value
		case "bmsgtag":
			logData.BMsgTag = value
		case "text":
			logData.Text = value
		case "to":
			logData.To = value
		case "from":
			logData.From = value
		case "dltentityid":
			logData.DLTEntityID = value
		case "bsms_outtime":
			logData.BSMSOutTime = value
		case "tid":
			logData.TID = value
		}
	}

	phoneNumbers := strings.Split(logData.To, ",")
	fmt.Println(len(phoneNumbers))
	if len(phoneNumbers) == 1 {
		globals.DataMap.Store(time.Now().Local(), *logData)
		fmt.Printf("Parsed Log Data: %+v\n", logData)
	}
	return nil
}
