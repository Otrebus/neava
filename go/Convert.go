package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Family struct {
	XMLName xml.Name `xml:"family"`
	Name    string   `xml:"name"`
	Year    string   `xml:"born"`
	Address *Address `xml:"address"`
	Phone   *Phone   `xml:"phone"`
}

type Phone struct {
	XMLName        xml.Name `xml:"phone"`
	MobileNumber   string   `xml:"mobile"`
	LandlineNumber string   `xml:"landline"`
}

type Address struct {
	XMLName xml.Name `xml:"address"`
	Street  string   `xml:"street"`
	City    string   `xml:"city"`
	Postal  string   `xml:"postal"`
}

type People struct {
	XMLName xml.Name `xml:"people"`
	Persons []Person
}

type Person struct {
	XMLName   xml.Name `xml:"person"`
	FirstName string   `xml:"firstname"`
	LastName  string   `xml:"lastname"`
	Phone     *Phone   `xml:"phone"`
	Address   *Address `xml:"address"`
	Family    []Family
}

type Parser struct {
	text   *string
	pos    int
	line   int
	column int
}

func (r *Parser) peek() byte {
	// Returns the byte at the current position in the parser
	if r.pos < len(*r.text) {
		return (*r.text)[r.pos]
	}
	return 0
}

func (r *Parser) adv() {
	// Advances the position of the parser
	r.column++
	if (*r.text)[r.pos] == '\n' {
		r.line++
		r.column = 1
	}
	r.pos++
}

func makeError(err string, parser Parser) error {
	// Augments an error message with the current line number
	return fmt.Errorf("Line %v, col %v: %v", parser.line, parser.column, err)
}

func acceptChar(byt byte, parser *Parser) bool {
	// Advances the parser position if the given character is at the current position
	if parser.peek() != byt {
		return false
	}
	parser.adv()
	return true
}

func acceptWhiteSpace(parser Parser) Parser {
	// Advances the parser position until the next non-whitechar character
	for acceptChar('\t', &parser) || acceptChar('\r', &parser) || acceptChar('\n', &parser) || acceptChar(' ', &parser) {
	}
	return parser
}

func acceptPipeVal(parser Parser) (string, Parser, error) {
	// Accepts a string of type "|<string>" until endline or EOF and returns the string
	if !acceptChar('|', &parser) {
		return "", parser, makeError("Expected '|'", parser)
	}

	var s strings.Builder
	for {
		c := parser.peek()
		if strings.Contains("|\r\n\000", string(c)) { // TODO: slow, make better contains function
			return s.String(), parser, nil
		}
		s.WriteByte(c)
		parser.adv()
	}
}

func parseP(parser Parser) (string, string, Parser, error) {
	// Parses a string P|<string>|<string>
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
	// Parses a string A|<string>|<string>|<string>
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
	// Parses a string F|<string>|<string>
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
	// Parses a string T|<string>|<string>
	char := acceptChar('T', &parser)
	if !char {
		return nil, parser, makeError("Expected T", parser)
	}
	mobile, newParser, err := acceptPipeVal(parser)
	if err != nil {
		return nil, newParser, makeError("Expected mobile number", newParser)
	}
	landline, newParser, err := acceptPipeVal(newParser)
	if err != nil {
		return nil, newParser, makeError("Expected landline number", newParser)
	}
	return &Phone{MobileNumber: mobile, LandlineNumber: landline}, newParser, nil
}

func parseFamily(parser Parser) (*Family, Parser, error) {
	// Parses the F of a P
	family, newParser, err := parseF(parser)
	if err != nil {
		return nil, parser, err
	}

	for true {
		newParser = acceptWhiteSpace(newParser)

		if newParser.peek() == 'A' {
			address, newParserA, err := parseA(newParser)
			if err != nil {
				return nil, parser, err
			}
			if family.Address != nil {
				// Here we read something that was already assigned which means we have a P
				// that assigns his own A after that of his F
				return family, newParser, nil
			}
			family.Address = address
			newParser = newParserA
		} else if newParser.peek() == 'T' {
			// This case is similar to A
			phone, newParserT, err := parseT(newParser)
			if err != nil {
				return nil, parser, err
			}
			if family.Phone != nil {
				return family, newParser, nil
			}
			family.Phone = phone
			newParser = newParserT
		} else {
			return family, newParser, nil
		}
	}
	return nil, Parser{}, nil
}

func parsePerson(parser Parser) (*Person, Parser, error) {
	// Parse information belonging to a Person - A 'P' and lines below it
	var person Person

	first, last, newParser, err := parseP(parser)
	if err != nil {
		return nil, parser, err
	}
	person.FirstName = first
	person.LastName = last

	for true {
		newParser = acceptWhiteSpace(newParser)
		if newParser.peek() == 'A' {
			address, newParserA, err := parseA(newParser)
			newParser = newParserA
			if err != nil {
				return nil, parser, err
			}
			person.Address = address
		} else if newParser.peek() == 'F' {
			family, newParserF, err := parseFamily(newParser)
			newParser = newParserF
			if err != nil {
				return nil, parser, err
			}
			person.Family = append(person.Family, *family)
		} else if newParser.peek() == 'T' {
			phone, newParserT, err := parseT(newParser)
			newParser = newParserT
			if err != nil {
				return nil, parser, err
			}
			person.Phone = phone
		} else {
			return &person, newParser, nil
		}
	}
	return nil, Parser{}, nil
}

func parsePeople(parser Parser) (*People, error) {
	// Parses a list of people
	var people People
	for true {
		parser = acceptWhiteSpace(parser)

		c := parser.peek()
		if c == '\000' {
			break
		}

		person, newParser, err := parsePerson(parser)
		parser = newParser
		if err != nil {
			return nil, err
		}
		people.Persons = append(people.Persons, *person)

	}
	return &people, nil
}

func stringToPeople(input string) (*People, error) {
	// Turns the input format into our intermediary data structure
	parser := Parser{&input, 0, 1, 1}
	return parsePeople(parser)
}

func peopleToXml(people People) []byte {
	// Turns our intermediary data structure to XML (we separated these into
	// functions so we can easily test them)
	out, _ := xml.MarshalIndent(people, "", "  ")
	return out
}

func main() {
	// Reads the proprietary people data format and turns it into XML
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString(0)

	people, err := stringToPeople(input)
	if err == nil {
		xml := peopleToXml(*people)
		fmt.Print(string(xml))
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}
