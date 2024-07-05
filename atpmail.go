package main

import (
	"sync"

	"github.com/emersion/go-imap/client"
)

func testEmailAuthentication(emails []UsersSctruct, Source SourceAddressDTO, results chan<- EmailAuthResult) {
	var wg sync.WaitGroup

	for _, email := range emails {
		wg.Add(1)
		go func(e UsersSctruct) {
			defer wg.Done()
			c, err := client.DialTLS(Source.Address, nil)
			if err != nil {
				results <- EmailAuthResult{Email: e.Username + "@" + Source.Domain, AuthError: err}
				return
			}

			defer c.Logout()

			err = c.Login(e.Username+"@"+Source.Domain, e.Password)
			if err != nil {
				results <- EmailAuthResult{Email: e.Username + "@" + Source.Domain, AuthError: err}
				return
			}
			results <- EmailAuthResult{Email: e.Username + "@" + Source.Domain, AuthError: err}
		}(email)
	}
	wg.Wait()
	close(results)
}
