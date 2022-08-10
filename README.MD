# cab
## Media *cab*inet management

## Build
```
go build -o ./cab ./app
```

## Usage
```
SUBCOMMANDS

list
Print a list of media that's been marked as done, doing, or to do
  -json
    	Output lists in JSON, default is table
  -status string
    	Status of the media to list (default "reading")
  -type string
    	Book or movie (default "book")

add
Add books and movies from Google Books and Letterboxd
  -title string
    	Name of the book or movie to add
  -type string
    	Book or movie (default "book")

backfill
Use Google Books to fill in missing author and ISBN data
```
