# People file converter

This is a technical trial written for Softhouse Neava. I made two solutions because halfway through the Go solution I realized a recursive-descent parser is a bit overkill for this task and not the best fit so after I completed it I made a simpler version in Python as well.

## Instructions

### Go

Run converter.exe with stdio as input, e.g. `converter.exe < input` or just run converter and manually type in the input, followed by EOF. Tests are also included: `go test`.

### Python

Reads from stdio as above, the command is `python convert.py < input` or just start `python convert.py` and type the input followed by EOF. To run tests, run `python -m unittest discover` in the python directory.

## Additional notes
The supplied example input in the problem description didn't match the specification of the input format since one of the addresses in the example input didn't include an area code. I assumed that this was an omission in the input rather than an optional part of the address. In the "real world" I'd ask to clarify but I didn't feel it was important enough for this test to warrant putting the task off in waiting for a reply.

I wanted to pick some language I have no experience in so I picked Go for my first solution to get a little bit of exposure to the language. I'm aware of some limitations, for example I expect my solution to break for multi-byte text formats like UTF-8 since I blindly loop through text byte-wise. Also I'm not entirely happy with the 'functional-style' treatment of the parser where I submit a parser and get a new parser in return; it's a little messy and maybe slightly slow - mutating the supplied parser might be cleaner. Also I didn't write the tests as exhaustively as one should e.g. with randomized tests like in the Python version.

Since this is a code test **I did not use any LLMs** when writing this code (which is probably obvious at a glance) except for a couple tangential questions e.g. involving .gitignore. In a real-world work situation my first instinct would be to basically copy/paste the instructions to an LLM and check its output plus make sure the tests are randomized and comprehensive.
