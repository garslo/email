# email

Package to fetch emails over POP3

# Installation

`go get github.com/garslo/email`

# Usage

```go
func main() {
	fetcher := NewGmailFetcher("example@gmail.com", "yourpassword")
    // Or
    fetcher = NewTlsEmailFetcher("example@gmail.com", "yourpassword", "pop.gmail.com", 995)
	msgs, err := fetcher.FetchEmails()
	DieIf(err)
	for _, msg := range msgs {
		body, err := ioutil.ReadAll(msg.Body)
		DieIf(err)
		fmt.Println(string(body))
	}
}

func DieIf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```
