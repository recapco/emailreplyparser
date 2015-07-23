# Email Reply Parser for Go

[![Build Status][travis-image]][travis-url] [![GoDoc][godoc-image]][godoc-url]

A Go port of GitHub's [Email Reply Parser][email_reply_parser] library. The
library parses an email body into fragments, marking the fragments as quoted or
as part of a signature as appropriate.

The most common use case is to get the text of a reply while ignoring signatures
and the quoted original email. It properly parses top, bottom and inline replies.

## Installation

To install `emailreplyparser` run:

```sh
go get github.com/recapco/emailreplyparser
```

## Usage

The library can be used to get the reply text from an email body as such:

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

The library can also help discover quoted segments in an email. For example:

```go
func Quotes(text string) ([]string, error) {
	email, err := emailreplyparser.Read(text)
	if err != nil {
		return nil, err
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

Please feel free to submit pull requests and issues.

## License

[MIT][license]

[email_reply_parser]: https://github.com/github/email_reply_parser
[license]: https://github.com/recapco/emailreplyparser/blob/master/LICENSE

[travis-url]: http://travis-ci.org/recapco/emailreplyparser
[travis-image]: http://img.shields.io/travis/recapco/emailreplyparser/master.svg?style=flat-square

[godoc-url]: https://godoc.org/github.com/recapco/emailreplyparser
[godoc-image]: https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square