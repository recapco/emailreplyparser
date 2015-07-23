# Email Reply Parser for Go
A Go port of GitHub's [Email Reply Parser](https://github.com/github/email_reply_parser) library.

## Usage
To parse out the reply body:

```go
reply, err := emailreplyparser.ParseReply(emailBody)
```

## Installation
`go get github.com/recapco/emailreplyparser`