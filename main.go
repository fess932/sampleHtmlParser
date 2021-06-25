package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("html.html")
	if err != nil {
		log.Println(err)
	}

	p := NewParser(f)

	log.Println(p.GetTitleAndDescription())
}

type Parser struct {
	r        io.ByteReader
	buf      bytes.Buffer
	nextByte int

	title       string
	description string

	currentTag int

	currentOpenedTag string
	tagOpened        bool

	inTitle bool
}

func NewParser(r io.Reader) *Parser {
	return &Parser{r: bufio.NewReader(r), nextByte: -1}
}

func (p *Parser) GetTitleAndDescription() (title, description string) {
	p.parse2()

	return p.title, p.description
}

const (
	tagOpen = iota
	tagClose
	tagSelfClosed
)

const (
	startTagC = byte('<')
	endTagC   = byte('>')
	closeTagC = byte('/')
)

const (
	tagStarted = '<'
	tagEnded   = '>'
	closeTag   = '/'
)

func (p *Parser) parse2() {
	var err error
	var s byte
	defer log.Println("title: ", p.title)
	defer log.Println("description: ", p.description)

	for {
		s, err = p.r.ReadByte()
		if err != nil {
			log.Println(err)
			return
		}

		//log.Println("in title", p.inTitle)

		//if p.inTitle {
		//	p.writeTitle(s)
		//}

		// ищем открывающий тег если тег не открыт
		if !p.tagOpened {
			p.findOpenTag(s)
			continue
		}
		// если тег открыт, если текущий символ закрывающий тег то закрываем тег, ичищаем буфер переходим на следующую итерацию
		if p.findCloseTag(s) {
			if p.inTitle {
				p.inTitle = false // при закрывающем title выходим из title
				p.title = p.buf.String()
			}
			log.Printf("{%s}", p.buf.String())
			tag := p.buf.String()

			if !strings.HasPrefix(tag, "/") {
				p.currentOpenedTag = tag
				log.Println("current opened tag", p.currentOpenedTag)
			}

			if strings.Contains(tag, "/head") {
				log.Println("------- end head -----------")
				return
			}

			p.buf.Reset()
			continue
		}

		if strings.Contains(p.buf.String(), "title") {
			p.inTitle = true
		}

		p.buf.WriteByte(s)
	}
}

func (p *Parser) writeTitle(b byte) {
	p.buf.WriteByte(b)
}

func (p *Parser) findOpenTag(b byte) {
	if b == startTagC {
		p.tagOpened = true
	}
}

func (p *Parser) findCloseTag(b byte) bool {
	if b == endTagC {
		p.tagOpened = false
		return true
	}
	return false
}

func (p *Parser) findHead() {

}

func (p *Parser) findEndHead() {

}
