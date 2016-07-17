package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"

	"pault.ag/go/wtf"
)

func ohshit(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) <= 2 {
		log.Fatalf("Usage: leveldb ...csv")
	}
	db, err := wtf.NewDatabase(os.Args[1])
	ohshit(err)

	for _, path := range os.Args[2:] {
		fd, err := os.Open(path)
		ohshit(err)

		r := csv.NewReader(fd)

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			ohshit(err)

			acronym := strings.ToLower(record[0])
			acronyms := wtf.Acronyms{}

			if err := db.Unpack(acronym, &acronyms); err != nil {
				if err != leveldb.ErrNotFound {
					panic(err)
				}
			}

			acronyms = append(acronyms, wtf.Acronym{
				Acronym:  record[0],
				Meaning:  record[1],
				Location: record[2],
				Note:     record[3],
			})
			ohshit(db.Set(acronym, acronyms))
		}
	}
}
