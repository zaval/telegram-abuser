package main

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	const (
		apiId   = <appId>
		apiHash = <apiHash>
	)

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  apiId,
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "TGAbuser",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}
	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Fatalf("SetLogVerbosityLevel error: %s", err)
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	fileData, err := os.ReadFile("publics.txt")
	if err != nil {
		log.Fatal(err)
	}

	for _, tgLink := range strings.Split(string(fileData), "\n") {

		tgLink = strings.TrimSpace(tgLink)
		var chat *client.Chat

		if strings.HasPrefix(tgLink, "https://t.me") {
			chat, err = tdlibClient.JoinChatByInviteLink(&client.JoinChatByInviteLinkRequest{
				InviteLink: tgLink,
			})
			if err != nil {
				log.Printf("Can't join the group %s: %s\n", tgLink, err)
				continue
			}
		} else {
			chat, err = tdlibClient.SearchPublicChat(&client.SearchPublicChatRequest{
				Username: tgLink,
			})
		}

		log.Printf("chat: %v", chat.Title)

		history, err := tdlibClient.GetChatHistory(&client.GetChatHistoryRequest{
			ChatId: chat.Id,
			Limit:  10000,
		})

		if err != nil {
			log.Printf("Can't get hostory of chat %s: %s", chat.Title, err)
			continue
		}

		var messageIds []int64
		for _, message := range history.Messages {
			messageIds = append(messageIds, message.Id)
		}

		_, err = tdlibClient.ReportChat(&client.ReportChatRequest{
			ChatId:     chat.Id,
			MessageIds: messageIds,
			Reason:     &client.ChatReportReasonFake{},
			Text:       "",
		})
		if err != nil {
			log.Printf("Can't report Fake: %s\n", err)
		} else {
			log.Println("Reported fake!")
			time.Sleep(time.Second * 5)
		}

		_, err = tdlibClient.ReportChat(&client.ReportChatRequest{
			ChatId:     chat.Id,
			MessageIds: messageIds,
			Reason:     &client.ChatReportReasonViolence{},
			Text:       "",
		})
		if err != nil {
			log.Printf("Can't report Violence: %s\n", err)
		} else {
			log.Println("reported violence!")
		}

		_, err = tdlibClient.LeaveChat(&client.LeaveChatRequest{ChatId: chat.Id})
		if err != nil {
			log.Printf("Can't leave chat %s, do it manually", tgLink)
		}
	}

}
