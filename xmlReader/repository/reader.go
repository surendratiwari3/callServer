package repository

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/beevik/etree"
)

func GetDocument(url string) *etree.Document{
	// Open our xmlFile
	xmlFile, err := http.Get(url)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened " + url)

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile.Body)

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(byteValue); err != nil {
		panic(err)
	}

	return doc
}
