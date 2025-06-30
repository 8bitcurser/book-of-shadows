watch:
	templ generate --watch
bs:
	templ generate && go build book-of-shadows && ./book-of-shadows