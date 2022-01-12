package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/kokizzu/gotro/L"
)

const (
	username = "misryarrazy@gmail.com"
	passwd   = "#"
	sender   = "vandyahmad2404@gmail.com"
)

func main() {
	cli, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = cli.Logout()
	}()

	if err = cli.Login(username, passwd); err != nil {
		log.Fatal(err)
	}

	// mailboxes := make(chan *imap.MailboxInfo, 10)
	// done := make(chan error, 1)
	// go func () {
	// 	done <- cli.List("", "*", mailboxes)
	// }()
	//
	// log.Println("Mailboxes:")
	// for m := range mailboxes {
	// 	log.Println("* " + m.Name)
	// }

	_, err = cli.Select(imap.InboxName, true)
	if err != nil {
		log.Fatal(err)
	}

	ids, err := cli.Search(&imap.SearchCriteria{
		Since: time.Now().Add(-1 * time.Hour),
		Body:  []string{"Done"},
		Header: textproto.MIMEHeader{
			"FROM": []string{sender},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(ids) == 0 {
		log.Println("not found")
		return
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- cli.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	// messageIds := make([]string, len(messages))
	var messageIds []string
	for msg := range messages {
		// messageIds = append(messageIds, msg.Envelope.InReplyTo)
		log.Println(fmt.Sprintf("%s -> %s", msg.Envelope.Subject, msg.Envelope.InReplyTo))
		messageIds = append(messageIds, msg.Envelope.InReplyTo)
	}

	// log.Println(len(messageIds))

	// <CAHeXJ8_nf3tL7BufFyZoSB8vMtuGY6Ym8AA9uC84NC0JxJtwQA@mail.gmail.com>

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	_, err = cli.Select("[Gmail]/Surat Terkirim", true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INI ..................")
	for _, v := range messageIds {

		ids, err = cli.Search(&imap.SearchCriteria{
			// Since: time.Now().Add(-24 * time.Hour),
			// Body: []string{"done"},
			Header: textproto.MIMEHeader{
				"Message-Id": []string{v},
			},
		})
		// fmt.Println(ids)
		if err != nil {
			log.Fatal(err)
		}

		if len(ids) == 0 {
			log.Println("not found")
			return
		}

		seqset = new(imap.SeqSet)
		seqset.AddNum(ids...)

		messages = make(chan *imap.Message, 10)
		done = make(chan error, 1)
		var section imap.BodySectionName
		go func() {
			done <- cli.Fetch(seqset, []imap.FetchItem{imap.FetchBody, imap.FetchBodyStructure, imap.FetchEnvelope, section.FetchItem()}, messages)
		}()
		// fmt.Println(messages)

		for msg := range messages {
			r := msg.GetBody(&section)
			if r == nil {
				log.Fatal("Server didn't returned message body")
			}

			// Create a new mail reader
			mr, err := mail.CreateReader(r)
			if err != nil {
				log.Fatal(err)
			}
			// i := 0
			// for {
			log.Println("Proses ambil pesan tiap part")
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				// log.Println("Got text:", string(b))
				text := string(b)
				sp := strings.Split(text, ":")
				// L.Describe(sp[0])
				// re := regexp.MustCompile("[0-9]+")
				re := regexp.MustCompile(`08\d{8,12}`)
				number := re.FindAllString(sp[1], -1)
				L.Describe(number)
				// L.Describe(split)
			}
			// i++
			// }
			// 	log.Println(fmt.Sprintf("%s -> %s", msg.Envelope.Subject, msg.Envelope.InReplyTo))
		}
		// log.Println(v)
	}
	log.Println("END ..................")
}
