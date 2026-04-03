from convert import convert_to_xml
import time
import unittest
import xml.dom.minidom as minidom
import random

def get_attrs(node, dic):
    """ Populates the attributes given in the keys of a dictionary from the
        children of an xml node """

    for n in node.childNodes:
        if n.nodeType == node.ELEMENT_NODE:
            if n.nodeName in dic:
                dic[n.nodeName] = n.firstChild.data if n.firstChild else ""


def rnd_string():
    """ Generates a random (possibly empty) string """
    s = ""
    for i in range(0, random.randint(0, 10)):
        s += random.choice("abcdefghijklmnopqrstuvxyz")
    return s


def rnd_address():
    """ Generates a random address """
    return f"A|{rnd_string()}|{rnd_string()}|{rnd_string()}\n"


def rnd_phone():
    """ Generates a random phone number """
    # We generate characters instead of numbers but format checking is
    # not part of the converter so it doesn't matter
    return f"T|{rnd_string()}|{rnd_string()}\n"


def rnd_phone_address():
    """ Generates a string possibly containing a P and a T line """
    
    s = ""
    
    if random.randint(0, 1):
        s += rnd_address()

    if random.randint(0, 1):
        s += rnd_phone()

    return s


def xml_to_original(xml):
    """ Converts XML back to the original format """

    s = ""

    doc = minidom.parseString(xml)
    people = doc.firstChild

    for node in people.childNodes:

        if node.nodeType == node.ELEMENT_NODE and node.nodeName == 'person':
            person = node
            
            dic = { "firstname": None, "lastname": None }

            get_attrs(person, dic)

            s += f"P|{dic['firstname']}|{dic['lastname']}\n"

            for node in person.childNodes:
                if node.nodeType == node.ELEMENT_NODE and node.nodeName == 'address':
                    dic = { "street": None, "city": None, "areacode": None }
                    get_attrs(node, dic)
                    s += f"A|{dic['street']}|{dic['city']}|{dic['areacode']}\n"

                if node.nodeType == node.ELEMENT_NODE and node.nodeName == 'phone':
                    dic = { "mobile": None, "landline": None }
                    get_attrs(node, dic)
                    s += f"T|{dic['mobile']}|{dic['landline']}\n"

                if node.nodeType == node.ELEMENT_NODE and node.nodeName == 'family':
                    family = node
                    dic = { "name": None, "born": None }
                    get_attrs(family, dic)
                    s += f"F|{dic['name']}|{dic['born']}\n"
                    for rel in family.childNodes:

                        if rel.nodeType == node.ELEMENT_NODE and rel.nodeName == 'address':
                            dic = { "street": None, "city": None, "areacode": None }
                            get_attrs(rel, dic)
                            s += f"A|{dic['street']}|{dic['city']}|{dic['areacode']}\n"

                        if rel.nodeType == node.ELEMENT_NODE and rel.nodeName == 'phone':
                            dic = { "mobile": None, "landline": None }
                            get_attrs(rel, dic)
                            s += f"T|{dic['mobile']}|{dic['landline']}\n"
    return s


class TestConverter(unittest.TestCase):

    def test_convert(self):
        """ Tests the converter by generating random input, converting to xml, and then converting back
            to the input format and comparing to the original input. """

        random.seed()

        for n in range(0, 100):

            seed = random.randint(0, 1000000)
            random.seed(seed)

            n_persons = random.randint(0, 2)

            input = ""

            for i in range(0, n_persons):
                input += f"P|{rnd_string()}|{rnd_string()}\n"

                n_family = random.randint(0, 2)

                # Here we always generate any A's before any T's because the XML converter will always output
                # the XML in that order and our reconverter converts it back in that order as well. We don't
                # want them flipped so the tests fail. So technically this is not 100% exhaustive testing
                input += rnd_phone_address()

                for i in range(0, n_family):
                    input += f"F|{rnd_string()}|{rnd_string()}\n"
                    input += rnd_phone_address()

            xml = convert_to_xml(input)
            original = xml_to_original(xml)

            self.assertEqual(input, original, f'seed is {seed}')
