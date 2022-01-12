package main

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("misryarrazy@gmail.com", "#"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	// from := uint32(1)
	to := mbox.Messages
	fmt.Println(to)
	// if mbox.Messages > 1 {
	// 	// We're using unsigned integers here, only subtract if the result is > 0
	// 	from = mbox.Messages - 1
	// }
	seqset := new(imap.SeqSet)
	seqset.AddRange(to, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("=====================================")
		tes := msg.Envelope.ReplyTo
		// fmt.Printf("%#v\n", msg.Envelope.From)

		for _, v := range tes {
			// log.Println("#########################")
			// L.Describe(v.PersonalName)
			log.Println(v.PersonalName)
			log.Println(v.MailboxName)
			log.Println(v.HostName)
			// log.Println("#########################")
		}
		// fmt.Println(tes)
		log.Println("Subject " + msg.Envelope.Subject)
		// log.Println("In REply" + msg.Envelope.InReplyTo)
		log.Println("=====================================")
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}
