package main

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func main() {

	// useStartTLS := true
	host := "imap.gmail.com"
	port := "993"
	uname := "misryarrazy@gmail.com"
	upass := "#"

	hostport := host + ":" + port

	log.Println("Using host and port:-", hostport)
	log.Println("With User:-", uname)
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)

	//c.SetDebug(os.Stdout)

	if err != nil {
		log.Fatal("APA INI", err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// if useStartTLS == true {

	// 	ret, err := c.SupportStartTLS()
	// 	if err != nil {
	// 		log.Println("Error trying to determine whether StartTLS is supported.")
	// 		log.Fatal(err)
	// 	}

	// 	if ret == false {
	// 		log.Println("StartTLS is not supported.")
	// 		log.Fatal(err)
	// 	}

	// 	log.Println("Good, StartTLS is supported.")

	// 	// Start a TLS session
	// 	tlsConfig := &tls.Config{ServerName: host}
	// 	if err := c.StartTLS(tlsConfig); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println("TLS started")
	// }

	// Login
	if err := c.Login(uname, upass); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	log.Println("Getting Last Message")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}

	// Get the last message
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)
	//since mbox.Messages is total number, this gets last

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	msg := <-messages
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	// Print some info about the message
	// header := mr.Header
	// if date, err := header.Date(); err == nil {
	// 	log.Println("Date:", date)
	// }
	// if from, err := header.AddressList("From"); err == nil {
	// 	log.Println("From:", from)
	// }
	// if to, err := header.AddressList("To"); err == nil {
	// 	log.Println("To:", to)
	// }
	// if subject, err := header.Subject(); err == nil {
	// 	log.Println("Subject:", subject)
	// }

	// Process each message's part
	for {
		log.Println("Proses ambil pesan tiap part")
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			log.Println("Got text:", string(b))
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Println("Got attachment: %v", filename)
		}
	}
}
