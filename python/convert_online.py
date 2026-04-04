

class ParseError(Exception):

    def __init__(self, line, str):
        return super().__init__(f"Line {line}: {str}")


def convert_to_xml(input):
    """ Alternative formulation of the converter with no intermediate representation """

    # I realized the parsing is so simple you can do it online. This is a quick-and-dirty
    # version and isn't really up to standard - a cleaner way to handle indentation and
    # closing of tags etc is to do this recursively instead of iteratively

    current_family = None
    current_person = None
    line = 1

    output = "<people>\n"

    for s in input.split('\n'):

        args = s.strip().split('|')
        if args[0] == 'P': # Person
            if len(args) != 3:
                raise ParseError(line, "P: not enough arguments")
            
            if current_family:
                output += "    </family>\n"

            if current_person:
                output += "  </person>\n"

            output += "  <person>\n"
            
            _, first_name, last_name = args
            current_family = None
            current_person = True

            output += f"    <firstname>{first_name}</firstname>\n"
            output += f"    <lastname>{last_name}</lastname>\n"

        if args[0] == 'F': # Family
            if not current_person:
                raise ParseError(line, "F: no current person")

            if len(args) != 3:
                raise ParseError(line, "F: not enough arguments")
            
            
            if current_family:
                output += "    </family>\n"
            
            _, name, birth = args
            current_family = True

            output += "    <family>\n"

            output += f"      <name>{name}</name>\n"
            output += f"      <born>{birth}</born>\n"

        if args[0] == 'T': # Telephone, assigns to the current person or current family [member]
            if not current_person and not current_family:
                raise ParseError(line, "T: no current person or family")

            if len(args) != 3:
                raise ParseError(line, "T: not enough arguments")
            
            _, mobile, landline = args
            # Assign to the current family if they have no telephone, otherwise this is repeated and we assign to the person
            indent = 6 if current_family else 4

            output += " "*indent + "<phone>\n"

            if mobile:
                output += " "*(indent+2) + f"<mobile>{mobile}</mobile>\n"

            if landline:
                output += " "*(indent+2) + f"<landline>{landline}</landline>\n"

            output += " "*indent + "</phone>\n"

        if args[0] == 'A': # Address, assigns to the current person or current family
            if not current_person and not current_family:
                raise ParseError(line, "A: no current person or family")

            if len(args) != 4:
                raise ParseError(line, "A: not enough arguments")
            
            _, street, city, areacode = args
            # Works the same as assigning telephone above
            indent = 6 if current_family else 4

            output += " "*indent + "<address>\n"

            if street:
                output += " "*(indent+2) + f"<street>{street}</street>\n"

            if city:
                output += " "*(indent+2) + f"<city>{city}</city>\n"

            if areacode:
                output += " "*(indent+2) + f"<areacode>{areacode}</areacode>\n"

            output += " "*indent + "</address>\n"

        line += 1

    if current_family:
        output += "    </family>\n"

    if current_person:
        output += "  </person>\n"

    output += "</people>"

    return output
