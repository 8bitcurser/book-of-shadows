## Setup

1. `go get`
2. `cd scripts`
3. `python3 -m venv venv`
4. `source ./scripts/venv/bin/activate`
5. `pip3 install -r requirements.txt`
6. `go build book-of-shadows`

## Current Features

- Generate random pulp cthulhu investigator
- Export to official PDF
- CRUD investigators with CookieStorage
- Cookie export through QR code or code line for another browser
- Investigator Wizard


## Why no centralized Storage?

I had two main constraints:
- Make it cheap
- Privacy First

I don't want to store any personal information from the users, this is a tool
that must stay free forever. Chaosium does not allow any profit from this projects
therefor I didn't want to incur in costs I would not be able to cost.

## But there is a data folder and I see some sqlite dependencies on the Docker

That's correct, I'm using the sqlite database as a temporary storage for exporting/importing 
the investigators to different browsers, this way the user can remain unknown. 
On the export we create a record on the DB that will live at most 24 hours that has the content of
the cookies from that user, they get a UUID that they can later use on another browser to import
the investigators.


## Cookie Challenge

Cookie size limit, some compression was needed to create the cookies in order to capture all
the relevant investigator data needed.


## Why I need python for this project if it runs on Golang

For the PDF export, golang libraries needed either a license or did not make the cut to manipulate
PDF forms effectively, python has a great library for it so running a script for it made sense. 
I could go a bit extra and run python from within Go but it was extra work and time I didn't have at the time.