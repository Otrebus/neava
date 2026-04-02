package main

import (
	"errors"
	"fmt"
	"strings"
)

type Family struct {
	Name    string
	Year    string
	Address Address
	Phone   Phone
}

type Phone struct {
	MobileNumber   string
	StandardNumber string
}

type Address struct {
	Street string
	City   string
	Postal string
}

type Person struct {
	FirstName string
	LastName  string
	Phone     Phone
	Address   Address
	Family    []Family
}

type Parser struct {
	text string
	pos  int
}

func (r *Parser) peek() byte {
	if r.pos < len(r.text) {
		return r.text[r.pos]
	}
	return 0
}

func (r *Parser) adv() {
	r.pos++
}

func makeError(err string, parser Parser) error {
	return errors.New(fmt.Sprintf("%v: %v", parser.pos, err))
}

func acceptChar(byt byte, parser *Parser) bool {
	if parser.peek() != byt {
		return false
	}
	parser.adv()
	return true
}

func acceptWhiteSpace(parser Parser) Parser {
	for acceptChar('\t', &parser) || acceptChar('\r', &parser) || acceptChar('\n', &parser) || acceptChar(' ', &parser) {
	}
	return parser
}

func acceptPipeVal(parser Parser) (string, Parser, error) {
	if !acceptChar('|', &parser) {
		return "", parser, makeError("Expected '|'", parser)
	}

	var s strings.Builder
	for {
		c := parser.peek()
		if strings.Contains("|\r\n\000", string(c)) { // todo: slow, make better contains function
			return s.String(), parser, nil
		}
		s.WriteByte(c)
		parser.adv()
	}
}

func parseP(parser Parser) (string, string, Parser, error) {
	char := acceptChar('P', &parser)
	if !char {
		return "", "", parser, makeError("Expected P", parser)
	}
	firstName, newParser, err := acceptPipeVal(parser)
	if err != nil {
		return "", "", newParser, makeError("Expected first name", newParser)
	}
	lastName, newParser, err := acceptPipeVal(newParser)
	if err != nil {
		return "", "", newParser, makeError("Expected last name", newParser)
	}
	return firstName, lastName, newParser, nil
}

func parseA(parser Parser) (*Address, Parser, error) {
	char := acceptChar('A', &parser)
	if !char {
		return nil, parser, makeError("Expected A", parser)
	}
	street, newParser, err := acceptPipeVal(parser)
	if err != nil {
		return nil, newParser, makeError("Expected street name", newParser)
	}
	city, newParser, err := acceptPipeVal(newParser)
	if err != nil {
		return nil, newParser, makeError("Expected city", newParser)
	}
	postal, newParser, err := acceptPipeVal(newParser)
	if err != nil {
		return nil, newParser, makeError("Expected postal number", newParser)
	}
	return &Address{Street: street, City: city, Postal: postal}, newParser, nil
}

func parseF(parser Parser) (*Family, Parser, error) {
	char := acceptChar('F', &parser)
	if !char {
		return nil, parser, makeError("Expected F", parser)
	}
	name, newParser, err := acceptPipeVal(parser)
	if err != nil {
		return nil, newParser, makeError("Expected name", newParser)
	}
	year, newParser, err := acceptPipeVal(newParser)
	if err != nil {
		return nil, newParser, makeError("Expected year", newParser)
	}
	return &Family{Name: name, Year: year}, newParser, nil
}

func parseT(parser Parser) (*Phone, Parser, error) {
	char := acceptChar('T', &parser)
	if !char {
		return nil, parser, makeError("Expected T", parser)
	}
	mobile, newParser, err := acceptPipeVal(parser)
	if err != nil {
		return nil, newParser, makeError("Expected mobile number", newParser)
	}
	stationary, newParser, err := acceptPipeVal(newParser)
	if err != nil {
		return nil, newParser, makeError("Expected stationary number", newParser)
	}
	return &Phone{MobileNumber: mobile, StandardNumber: stationary}, newParser, nil
}

func parseFamily(parser Parser) (*Family, Parser, error) {

	family, newParser, err := parseF(parser)
	if err != nil {
		return nil, parser, err
	}

	for true {
		newParser = acceptWhiteSpace(newParser)
		address, parserA, err1 := parseA(newParser)
		if err1 == nil {
			family.Address = *address
		}
		phone, parserT, err2 := parseT(parserA)
		if err2 == nil {
			family.Phone = *phone
		}
		if err1 != nil && err2 != nil {
			return family, newParser, nil
		}
		newParser = parserT
	}
	return nil, Parser{}, nil
}

func parsePerson(parser Parser) (*Person, Parser, error) {

	var person Person

	first, last, newParser, err := parseP(parser)
	if err != nil {
		return nil, parser, err
	}
	person.FirstName = first
	person.LastName = last

	for true {
		newParser = acceptWhiteSpace(newParser)
		address, newParserA, err1 := parseA(newParser)
		if err1 == nil {
			person.Address = *address
		}
		family, newParserF, err2 := parseFamily(newParserA)
		if err2 == nil {
			person.Family = append(person.Family, *family)
		}
		phone, newParserT, err3 := parseT(newParserF)
		if err3 == nil {
			person.Phone = *phone
		}

		if err1 != nil && err2 != nil && err3 != nil {
			return &person, newParserT, nil
		}

		newParser = newParserT
	}
	return nil, Parser{}, nil
}

func parsePeople(parser Parser) []*Person {
	var people []*Person
	for true {
		parser = acceptWhiteSpace(parser)
		person, newParser, err := parsePerson(parser)
		parser = newParser
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		people = append(people, person)

		c := parser.peek()
		if c == '\000' {
			return people
		}
	}
	return people
}

func printPeople(people []*Person) {
	for _, person := range people {
		fmt.Println(person.FirstName)
		for _, family := range person.Family {
			fmt.Println(family.Name)
		}
	}
}

func main() {
	input :=
		`	P|Carl Gustaf|Bernadotte
			T|0768-101801|08-101801
			A|Drottningholms slott|Stockholm|10001
			F|Victoria|1977
			A|Haga Slott|Stockholm|10002
			F|Carl Philip|1979
			T|0768-101802|08-101802
			P|Barack|Obama
			A|1600 Pennsylvania Avenue|Washington, D.C`

	parser := Parser{input, 0}

	people := parsePeople(parser)

	printPeople(people)
}
