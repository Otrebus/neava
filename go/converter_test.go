package main

import (
	"testing"
)

// No fancy testing here, just a few hard-coded examples. Ideally one should add some sort
// of randomized testing akin to that of the Python version

func TestConverter(t *testing.T) {
	input :=
		`	P|Carl Gustaf|Bernadotte
			T|0768-101801|08-101801
			A|Drottningholms slott|Stockholm|10001
			F|Victoria|1977
			A|Haga Slott|Stockholm|10002
			F|Carl Philip|1979
			T|0768-101802|08-101802
			P|Barack|Obama
			A|1600 Pennsylvania Avenue|Washington, D.C|90210`

	p, _ := stringToPeople(input)
	a := string(peopleToXml(*p))

	b :=
		`<people>
  <person>
    <firstname>Carl Gustaf</firstname>
    <lastname>Bernadotte</lastname>
    <phone>
      <mobile>0768-101801</mobile>
      <landline>08-101801</landline>
    </phone>
    <address>
      <street>Drottningholms slott</street>
      <city>Stockholm</city>
      <postal>10001</postal>
    </address>
    <family>
      <name>Victoria</name>
      <born>1977</born>
      <address>
        <street>Haga Slott</street>
        <city>Stockholm</city>
        <postal>10002</postal>
      </address>
    </family>
    <family>
      <name>Carl Philip</name>
      <born>1979</born>
      <phone>
        <mobile>0768-101802</mobile>
        <landline>08-101802</landline>
      </phone>
    </family>
  </person>
  <person>
    <firstname>Barack</firstname>
    <lastname>Obama</lastname>
    <address>
      <street>1600 Pennsylvania Avenue</street>
      <city>Washington, D.C</city>
      <postal>90210</postal>
    </address>
  </person>
</people>`

	if a != b {
		t.Errorf("Mismatch")
	}
}

func TestConverter2(t *testing.T) {
	input :=
		`	P|Carl Gustaf|Bernadotte

			P|Barack|Obama

			A|1600 Pennsylvania Avenue|Washington, D.C|90210`

	p, _ := stringToPeople(input)
	a := string(peopleToXml(*p))

	b :=
		`<people>
  <person>
    <firstname>Carl Gustaf</firstname>
    <lastname>Bernadotte</lastname>
  </person>
  <person>
    <firstname>Barack</firstname>
    <lastname>Obama</lastname>
    <address>
      <street>1600 Pennsylvania Avenue</street>
      <city>Washington, D.C</city>
      <postal>90210</postal>
    </address>
  </person>
</people>`

	if a != b {
		t.Errorf("Mismatch")
	}
}

func TestConverter3(t *testing.T) {
	input := ``

	p, _ := stringToPeople(input)
	a := string(peopleToXml(*p))

	b := `<people></people>`

	if a != b {
		t.Errorf("Mismatch")
	}
}

func TestErrors(t *testing.T) {
	input :=
		`	T|0768-101801|08-101801
			A|Drottningholms slott|Stockholm|10001
		`
	_, err := stringToPeople(input)
	if err == nil {
		t.Errorf("Expected error")
	}

	input = `	|| `
	_, err = stringToPeople(input)
	if err == nil {
		t.Errorf("Expected error")
	}

	input = `	P|Carl Gustaf|Bernadotte||`
	_, err = stringToPeople(input)
	if err == nil {
		t.Errorf("Expected error")
	}

	input = `	P|Carl Gustaf`
	_, err = stringToPeople(input)
	if err == nil {
		t.Errorf("Expected error")
	}

	input = `	P|Carl Gustaf | | |`
	_, err = stringToPeople(input)
	if err == nil {
		t.Errorf("Expected error")
	}

	input = `P`
	_, err = stringToPeople(input)
	if err == nil {
		t.Errorf("Expected error")
	}

	input =
		`	P|Carl Gustaf|Bernadotte
			F|Victoria|1977
			A|
      

      `
	_, err = stringToPeople(input)
	if err == nil {
		t.Errorf("Expected error")
	}
}
