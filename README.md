# Email Reply Parser for Go

[![Build Status][travis-image]][travis-url] [![Coverage][coveralls-image]][coveralls-url] [![GoDoc][godoc-image]][godoc-url]

A Go port of GitHub's [Email Reply Parser][email_reply_parser] library. The 
library is used to strip away non essential content from email bodies. An 
example use case is to allow email replies to comments without including 
signatures or extra noise.

## Installation

To install `emailreplyparser` run the standard installation:

```sh
go get github.com/recapco/emailreplyparser
```

## Usage

The library can be used to get the essential text of body with following example:

```go
reply, err := emailreplyparser.ParseReply(emailBody)
```

The library can also be used to retrieve the signature. For example:

```go
func Signature(text string) (string, error) {
	email, err := emailreplyparser.Read(text)
	if err != nil {
		return "", err
	}
	for _, fragment := range email.Fragments {
		if fragment.Signature {
			return fragment.String(), nil
		}
	}
	return "", nil
}
```

The library can also help discover quotes in an email. For example:

```go
func Quotes(text string) ([]string, error) {
	email, err := emailreplyparser.Read(text)
	if err != nil {
		return [], err
	}
	var quotes []string
	for _, fragment := range email.Fragments {
		if fragment.Quoted {
			quotes = append(quotes, fragment.String())
		}
	}
	return quotes, nil
}
```

## Building and Testing

Building and testing follow the normal Go conventions of `go build` and 
`go test`.

## Contributing

Please feel to submit Pull Requests and Issues.

## License

[MIT][license]

[email_reply_parser]: https://github.com/github/email_reply_parser
[license]: https://github.com/recapco/emailreplyparser/blob/master/LICENSE

[travis-url]: http://travis-ci.org/recapco/emailreplyparser
[travis-image]: http://img.shields.io/travis/recapco/emailreplyparser/master.svg?style=flat-square

[coveralls-url]: https://coveralls.io/r/recapco/emailreplyparser
[coveralls-image]: https://img.shields.io/coveralls/recapco/emailreplyparser/master.svg?style=flat-square

[godoc-url]: https://godoc.org/github.com/recapco/emailreplyparser
[godoc-image]: https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square