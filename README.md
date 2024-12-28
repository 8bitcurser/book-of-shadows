## Setup

1. `go get`
2. `source ./scripts/venv/bin/activate`
3. `pip3 install -r requirements.txt`
3. `go build book-of-shadows`

## Current Features

- Generate random pulp cthulhu investigator
- Export to official PDF
- CRUD investigators with CookieStorage

## Pending Features

- Cookie export through QR code or code line for another browser
- Investigator Wizard

## Why no centralized Storage?

I had two main constraints:
- Make it cheap
- Privacy First

I don't want to store any personal information from the users, this is a tool
that must stay free forever. Chaosium does not allow any profit from this projects
therefor I didn't want to incur in costs I would not be able to cost.

## Cookie Challenge

Cookie size limit, some compression was needed to create the cookies in order to capture all
the relevant investigator data needed.


## Why I need python for this project if it runs on Golang

For the PDF export, golang libraries needed either a license or did not make the cut to manipulate
PDF forms effectively, python has a great library for it so running a script for it made sense. 
I could go a bit extra and run python from within Go but it was extra work and time I didn't have at the time.