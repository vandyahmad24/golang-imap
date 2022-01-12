package main

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
)

func entah() {
	// log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("misryarrazy@gmail.com", "#"); err != nil {
		log.Fatal(err)
	}
	// log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	to := mbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(to, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	// log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("=====================================")
		tes := msg.Envelope.ReplyTo
		// fmt.Printf("%#v\n", msg.Envelope.From)

		for _, v := range tes {
			// log.Println("#########################")
			// L.Describe(v.PersonalName)
			log.Println(v.PersonalName)
			email := fmt.Sprintf("Email : %s@%s", v.MailboxName, v.HostName)
			log.Println(email)
			log.Println(v.HostName)
			// log.Println("#########################")
		}
		// fmt.Println(tes)
		log.Println("Subject: " + msg.Envelope.Subject)
		// log.Println("In REply" + msg.Envelope.InReplyTo)
		log.Println("=====================================")
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}

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

	// log.Println("Mailboxes:")
	// for m := range mailboxes {
	// 	log.Println("*== " + m.Name)
	// }

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Flags for INBOX:", mbox.Flags)
	// getMail2(c)

	idleClient := idle.NewClient(c)
	// Create a channel to receive mailbox updates
	updates := make(chan client.Update)
	c.Updates = updates

	// Check support for the IDLE extension
	if ok, err := idleClient.SupportIdle(); err == nil && ok {
		// Start idling
		// stopped := false
		stop := make(chan struct{})
		done := make(chan error, 1)
		go func() {
			done <- idleClient.Idle(stop)
		}()

		// Listen for updates
		for {
			select {
			case update := <-updates:
				//https://github.com/emersion/go-imap-idle/issues/11
				// have
				// entah()
				switch update.(type) {
				case *client.MessageUpdate:
					log.Println("Found message update type")
					// msg, _ := update.(*client.MessageUpdate) //This prints what you want here
					// // log.Println("ini dari message update", msg.Message)
					// L.Describe(msg)
				case *client.MailboxUpdate:
					log.Println("Found mailbox update type")
					// mbx, _ := update.(*client.MailboxUpdate)
					// // log.Println(mbx.Mailbox)
					// L.Describe(mbx.Mailbox.Messages)
					entah()
				// log.Println("UnseenSeqNum:", mbx.Mailbox.UnseenSeqNum) //邮箱中第一封未读邮件的序列号。
				// log.Println("Messages:", mbx.Mailbox.Messages)         //此邮箱中的邮件数。
				// log.Println("Recent:", mbx.Mailbox.Recent)             //自上次打开邮箱以来未看到的邮件数。
				// log.Println("UidNext:", mbx.Mailbox.UidNext)           //

				// log.Println("UidValidity:", mbx.Mailbox.UidValidity) //与UID一起，它是消息的唯一标识符。必须大于或等于1。
				// getMail(c)
				// go sendWXworkMsg()

				// lib.FetchLast(c, mbx.Mailbox)
				//etc....
				default:
					log.Println("skipping update")
				}

				// if !stopped {
				// 	close(stop)
				// 	stopped = true
				// }
			case err := <-done:
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Not idling anymore?")
				return
			}
		}
	}

	log.Println("Done!")
}
