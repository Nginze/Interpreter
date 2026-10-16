package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type TokenType int

const (
	KEYWORD TokenType = iota
	IDENT
	STRING
	INT
	FLOAT
	BOOL
	OPERATOR
	PUNCT
)

var tokenMap = map[int]string{
	0: "KEYWORD",
	1: "IDENT",
	2: "STRING",
	3: "INT",
	4: "FLOAT",
	5: "BOOL",
	6: "OPERATOR",
	7: "PUNCT",
}

var operatorRe = regexp.MustCompile(`[+\-*/=%><&^|]`)
var numberRe = regexp.MustCompile(`^\d+$`)
var floatRe = regexp.MustCompile(`^\d+\.\d+$`)
var punctRe = regexp.MustCompile(`[(){}\[\],;:.]`)

// var stringRe = regexp.MustCompile(`^"[^"]*"$`)
var keywordsRe = regexp.MustCompile(`^(if|else|for|while|return|fn|let)$`)

type Token struct {
	t   TokenType
	val string
}

func (tok Token) String() string {
	return fmt.Sprintf("Token(Type: %v, Value: %q)", tokenMap[int(tok.t)], tok.val)
}

type Lexer struct {
	stream *bufio.Reader
	output []Token
}

func NewLexer(fpath *string) *Lexer {
	bytes, err := os.ReadFile(*fpath)

	if err != nil {
		fmt.Printf("Failed to read file: %v \n", err)
	}

	fileString := string(bytes)

	fmt.Printf("%q\n", fileString)

	stream := bufio.NewReader(strings.NewReader(fileString))

	return &Lexer{stream: stream}

}

func (l *Lexer) Read(nBytes int) (string, error) {
	buf := make([]byte, nBytes)
	n, err := l.stream.Read(buf)

	if err != nil && err != io.EOF {
		fmt.Printf("Error reading ch to buffer\n", err)
		return "", errors.New("Error reading ch to buffer")
	}

	return string(buf[:n]), nil
}

func (l *Lexer) Peek(nbytes int) (string, error) {
	buf, err := l.stream.Peek(nbytes)

	if err != nil {
		// fmt.Println("Error peeking ch ", err)
		return "", errors.New("Error peeking ch ")
	}

	return string(buf), nil
}

func (l *Lexer) NextChunk() string {
	n := 1

	for {
		b, err := l.Peek(n)

		if err != nil {
			if len(b) > 0 {
				return string(b)
			}
			return ""
		}

		lastChar := rune(b[len(b)-1])
		if unicode.IsSpace(lastChar) || operatorRe.MatchString(string(lastChar)) || punctRe.MatchString(string(lastChar)) {
			// Found a whitespace boundary or operator boundary or punctuation boundary.
			// (indentifiers don't have specidal chars)! Return everything up to (but excluding) it
			return string(b[:len(b)-1])
		}

		n++

	}

}

func (l *Lexer) NextString() string {
	n := 2

	for {
		b, err := l.stream.Peek(n)
		if err != nil {
			if len(b) > 0 {
				return string(b)
			}
			return ""
		}

		lastByte := b[len(b)-1]

		if lastByte == '"' {
			return string(b)
		}

		n++
	}
}

func (l *Lexer) EscapeWhiteSpace() {
	for {
		r, _, err := l.stream.ReadRune()
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}

		if !unicode.IsSpace(r) {
			_ = l.stream.UnreadRune()
			return
		}
	}
}

func (l *Lexer) Tokenize() {
	for {
		l.EscapeWhiteSpace()
		// peek current token
		ch, err := l.Peek(1)
		var token Token
		// peek needs to return an error when we reach the eof

		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		switch {
		case operatorRe.MatchString(ch):
			fmt.Printf("Matched operator:  %v in stream \n", ch)
			chunk, err := l.Peek(2)

			if err != nil {
				val, err := l.Read(1)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}
			} else if chunk == "++" {
				val, err := l.Read(2)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}

			} else if chunk == "==" {
				val, err := l.Read(2)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}
			} else if chunk == ">=" {
				val, err := l.Read(2)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}
			} else if chunk == "<=" {
				val, err := l.Read(2)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}
			} else if chunk == "&&" {
				val, err := l.Read(2)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}
			} else if chunk == "||" {
				val, err := l.Read(2)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}
			} else {

				val, err := l.Read(1)

				if err != nil {
					continue
				}

				token = Token{
					t:   OPERATOR,
					val: val,
				}
			}

		case punctRe.MatchString(ch):
			fmt.Printf("Matched punctuation :  %v in stream \n", ch)

			val, err := l.Read(1)

			if err != nil {
				l.Read(1)
				continue
			}

			token = Token{
				t:   PUNCT,
				val: val,
			}

		case numberRe.MatchString(ch):

			fmt.Printf("Matched number:  %v in stream as number\n", ch)
			chunk := l.NextChunk()

			fmt.Printf("Number chunk: %v \n", chunk)

			if !numberRe.MatchString(chunk) && !floatRe.MatchString(chunk) {
				fmt.Println("Invalid Token encountered")
			}

			if floatRe.MatchString(chunk) {
				val, err := l.Read(len(chunk))

				if err != nil {
					continue
				}

				token = Token{
					t:   FLOAT,
					val: val,
				}
			} else {
				val, err := l.Read(len(chunk))

				if err != nil {
					continue
				}

				token = Token{
					t:   INT,
					val: val,
				}

			}
		case ch == string('"'):

			fmt.Printf("Matched string : %v in stream as string", ch)
			chunk := l.NextString()

			val, err := l.Read(len(chunk))

			if err != nil {
				continue
			}

			token = Token{
				t:   STRING,
				val: val,
			}

		default:

			fmt.Printf("Matched Identifier or Keyword:  %v in stream\n", ch)
			chunk := l.NextChunk()

			fmt.Printf("IDENT/KEYWORD Chunk %v: \n", chunk)

			if chunk == "true" || chunk == "false" {
				val, err := l.Read(len(chunk))

				if err != nil {
					continue
				}

				token = Token{
					t:   BOOL,
					val: val,
				}
			} else if keywordsRe.MatchString(chunk) {

				val, err := l.Read(len(chunk))

				if err != nil {
					continue
				}

				token = Token{
					t:   KEYWORD,
					val: val,
				}
			} else {

				val, err := l.Read(len(chunk))

				if err != nil {
					continue
				}

				token = Token{
					t:   IDENT,
					val: val,
				}
			}

		}

		l.output = append(l.output, token)
	}

}

func (l *Lexer) PrintAll() {
	for i, tok := range l.output {
		fmt.Printf("Token %d: %s\n", i, tok)
	}
}
