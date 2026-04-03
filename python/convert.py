import sys
from xml.dom import minidom


class ParseError(Exception):

    def __init__(self, line, str):
        return super().__init__(f"Line {line}: {str}")


class Telephone:
    def __init__(self):
        self.mobile = None
        self.landline = None


class Address:
    def __init__(self):
        self.street = None
        self.city = None
        self.areacode = None


class Family:
    def __init__(self):
        self.name = None
        self.born = None
        self.telephone = None
        self.address = None


class Person:
    def __init__(self):
        self.family = []
        self.address = None
        self.first_name = None
        self.last_name = None
        self.telephone = None
        self.family = []


def create_xml_element(root, parent, name, dic):
    """ Creates an xml element under a parent with a name with attributes given
        in a supplied dictionary """

    element = root.createElement(name)
    parent.appendChild(element)

    for key, val in dic.items():
        if val is None:
            continue
        child_ele = root.createElement(key)
        element.appendChild(child_ele)
        child_ele.appendChild(root.createTextNode(val))

    return element


def convert_to_xml(input):
    """ The main function that converts a list of people to XML format """

    current_person = None
    current_family = None
    line = 1
    persons = []

    for s in input.split('\n'):

        args = s.strip().split('|')
        if args[0] == 'P': # Person
            if len(args) != 3:
                raise ParseError(line, "P: not enough arguments")
            
            _, first_name, last_name = args
            current_person = Person()
            current_person.first_name = first_name
            current_person.last_name = last_name
            current_family = None

            persons.append(current_person)

        if args[0] == 'F': # Family
            if not current_person:
                raise ParseError(line, "F: no current person")

            if len(args) != 3:
                raise ParseError(line, "F: not enough arguments")
            
            _, name, birth = args
            current_family = Family()
            current_family.name = name
            current_family.born = birth
            current_person.family.append(current_family)

        if args[0] == 'T': # Telephone, assigns to the current person or current family [member]
            if not current_person and not current_family:
                raise ParseError(line, "T: no current person or family")

            if len(args) != 3:
                raise ParseError(line, "T: not enough arguments")
            
            _, mobile, landline = args
            # Assign to the current family if they have no telephone, otherwise this is repeated and we assign to the person
            phone = (current_family if current_family and not current_family.telephone else current_person).telephone = Telephone()
            phone.mobile = mobile
            phone.landline = landline

        if args[0] == 'A': # Address, assigns to the current person or current family
            if not current_person and not current_family:
                raise ParseError(line, "A: no current person or family")

            if len(args) != 4:
                raise ParseError(line, "A: not enough arguments")
            
            _, street, city, areacode = args
            # Works the same as assigning telephone above
            address = (current_family if current_family and not current_family.address else current_person).address = Address()
            address.street = street
            address.city = city
            address.areacode = areacode

        line += 1

    root = minidom.Document()

    people = root.createElement('people') 
    root.appendChild(people)

    # Generate the actual xml
    for person in persons:

        person_element = create_xml_element(root, people, 'person', {
            "firstname": person.first_name,
            "lastname": person.last_name
        })

        if person.address:
            create_xml_element(root, person_element, 'address', {
                    "street": person.address.street,
                    "city": person.address.city,
                    "areacode": person.address.areacode
                }
            )

        if person.telephone:
            create_xml_element(root, person_element, 'phone', {
                "mobile": person.telephone.mobile,
                "landline": person.telephone.landline
            })

        
        for family in person.family:
            family_element = create_xml_element(root, person_element, 'family', {
                    "name": family.name,
                    "born": family.born
                }
            )

            if family.address:
                create_xml_element(root, family_element, 'address', {
                        "street": family.address.street,
                        "city": family.address.city,
                        "areacode": family.address.areacode
                    }
                )

            if family.telephone:
                create_xml_element(root, family_element, 'phone', {
                        "mobile": family.telephone.mobile,
                        "landline": family.telephone.landline
                    }
                )

    xml_str = people.toprettyxml(indent ="  ")

    return xml_str


def main():
    """ Reads the input format from stdin and prints out the corresponding XML to stdout """
    input = ""

    while True:
        s = sys.stdin.readline()
        if not s:
            break 
        input += s + '\n'

    xml_str = convert_to_xml(input)
    print(xml_str)


if __name__ == '__main__':
     main()
