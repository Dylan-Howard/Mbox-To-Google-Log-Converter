package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

const EXPORT_PREFIX = "output_"
const LOG_HEADERS = "Message ID,Start date,End date,Sender,Message size,Subject,Direction,Attachments,Recipient address,Event target,Event date,Event status,Event target IP address,Has encryption,Event SMTP reply code,Event description,Client Type,Device User Session ID"

type Mbox struct {
	DeliveredTo								string
	Received									[]string
	XGoogleSmtpSource					string
	XReceived									string
	ARCSeal										string
	ARCMessageSignature				string
	ARCAuthenticationResults	string
	ReturnPath 								string
	ReceivedSPF 							string
	AuthenticationResults 		string
	DKIMSignature 						string
	Date 											string
	From 											string
	MessageID 								string
	Subject 									string
	MIMEVersion 							string
	XSGEID 										string
	To 												string
	XEntityID 								string
	ContentType 							string
	ContentTransferEncoding		string
	Content 									string
}

func (m Mbox) GetAsLogRow() string {
	return fmt.Sprintf(
		"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"",
		m.MessageID,
		m.Date,
		m.Date,
		m.ReturnPath,
		"",
		m.Subject,
		"Received",
		"0",
		m.To,
		"GMAIL_INBOX",
		m.Date,
		"DELIVERED",
		"",
		"Not encrypted",
		"",
		"No Error",
		"",
		"",
	)
}

type MboxGenerator struct {
	ImportDirectory		string
	ExportDirectory		string
	ExportBlockCount	int
}

func (g MboxGenerator) GetDataFiles() []string {
	files := []string{}

	entries, err := os.ReadDir(g.ImportDirectory)
  if err != nil {
    log.Fatal(err)
		return files
  }
 
  for _, e := range entries {
		files = append(files, e.Name())
  }

	return files
}

func (g MboxGenerator) ParseFileContents(fileName string) Mbox {
  data, err := os.ReadFile(path.Join(g.ImportDirectory, fileName))
  if err != nil {
    panic(err)
  }

	rows := strings.Split(string(data[:]), "\r\n")
	
	/* Parse data rows */
	mbox := Mbox{}

	for i := 0; i < len(rows); i++ {
		prefixIndex := strings.Index(rows[i], ": ")

		if prefixIndex < 0 {
			if rows[i] == "" {
				mbox.Content = strings.Join(rows[i:], "\n")
				i = len(rows) - 1
			}

			continue
		}

		switch rows[i][0:prefixIndex] {
		case "Delivered-To":
			mbox.DeliveredTo = rows[i][prefixIndex+2:len(rows[i])]

		case "Received":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.Received = append(mbox.Received, rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n"))

			i = j

		case "X-Google-Smtp-Source":
			mbox.XGoogleSmtpSource = rows[i][prefixIndex+2:len(rows[i])]
		case "X-Received":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.XReceived = rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n")

			i = j

		case "ARC-Seal":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.ARCSeal = rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n")

			i = j

		case "ARC-Message-Signature":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.ARCMessageSignature = rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n")

			i = j

		case "ARC-Authentication-Results":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.ARCAuthenticationResults = rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n")

			i = j

		case "Return-Path":
			mbox.ReturnPath = rows[i][prefixIndex+2:len(rows[i])]

		case "Received-SPF":
			mbox.ReceivedSPF = rows[i][prefixIndex+2:len(rows[i])]

		case "Authentication-Results":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.AuthenticationResults = rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n")

			i = j

		case "DKIM-Signature":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.DKIMSignature = rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n")

			i = j

		case "Date":
			mbox.Date = rows[i][prefixIndex+2:len(rows[i])]

		case "From":
			mbox.From = rows[i][prefixIndex+2:len(rows[i])]

		case "Message-ID":
			mbox.MessageID = rows[i][prefixIndex+2:len(rows[i])]

		case "Subject":
			mbox.Subject = rows[i][prefixIndex+2:len(rows[i])]

		case "MIME-Version":
			mbox.MIMEVersion = rows[i][prefixIndex+2:len(rows[i])]

		case "X-SG-EID":
			j := i + 1
			for j < len(rows) {
				if (rows[j][0] != ' ') {
					break
				}
				j++
			}

			mbox.XSGEID =  rows[i][prefixIndex+2:len(rows[i])] + strings.Join(rows[i+1:j], "\n")

		case "To":
			mbox.To = rows[i][prefixIndex+2:len(rows[i])]

		case "X-Entity-ID":
			mbox.XEntityID = rows[i][prefixIndex+2:len(rows[i])]

		case "Content-Type":
			mbox.ContentType = rows[i][prefixIndex+2:len(rows[i])]

		case "Content-Transfer-Encoding":
			mbox.ContentTransferEncoding = rows[i][prefixIndex+2:len(rows[i])]

		}
	}

	return mbox
}

func (g MboxGenerator) GetMboxes() []Mbox {
	files := g.GetDataFiles()
	mboxes := []Mbox{}
 
  for _, f := range files {
		mboxes = append(mboxes, g.ParseFileContents(f))
  }

	return mboxes
}

func writeFile(filepath string, content string) error {
	data := []byte(content)

  err := os.WriteFile(filepath, data, 0644)

  if err != nil {
		return err
	}

	return nil
}

func (g MboxGenerator) GenerateLogs(boxes []Mbox) error {

	rows := []string{LOG_HEADERS}
	for i := 0; i < len(boxes); i++ {
		rows = append(rows, boxes[i].GetAsLogRow())
	}

	data := strings.Join(rows, "\n")

	writeFile(path.Join(g.ExportDirectory, fmt.Sprintf("%s%s%s", EXPORT_PREFIX, strconv.Itoa(1), ".csv")), data)

	return nil
}