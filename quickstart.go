package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"

	mesg, err := srv.Users.Messages.Get(user, "166063e1fa906699").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve message %v: %v", "166063e1fa906699", err)
	}
	for _, msgLabel := range mesg.LabelIds {
		fmt.Printf("- Labels:: %s\n", msgLabel)
	}

	/*
		r, err := srv.Users.Labels.List(user).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve labels: %v", err)
		}
		if len(r.Labels) == 0 {
			fmt.Println("No labels found.")
			return
		}
		fmt.Println("Labels:")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n", l.Name)
		}

		ms, err := srv.Users.Messages.List(user).Q("is:UNREAD").Do()

		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}
		if len(ms.Messages) == 0 {
			fmt.Println("No messages found.")
			return
		}
		fmt.Println("Messages:")
		for _, l := range ms.Messages {
			msg, err := srv.Users.Messages.Get("me", l.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", l.Id, err)
			}
			for _, mailLabel := range msg.LabelIds {

				//fmt.Printf("- %s\n", mailLabel)
				if mailLabel == "UNREAD" {
					fmt.Printf("- MailID: %s\n", l.Id)

					for _, mailHeader := range msg.Payload.Headers {

						if mailHeader.Name == "Subject" {
							fmt.Printf("- %s\n", mailHeader.Value)
						}

						if mailHeader.Name == "Date" {
							fmt.Printf("- %s\n", mailHeader.Value)
						}
					}
				}
			}

			for _, mailHeader := range msg.Payload.Headers {

				fmt.Printf("- %s :: %s\n", mailHeader.Name, mailHeader.Value)

				if mailHeader.Name == "Subject" {
					fmt.Printf("- %s\n", mailHeader.Value)
				}

			}

			fmt.Printf("=======================================================\n")
		}
	*/
}
