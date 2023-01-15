// Bible - Holy bible CLI
// https://github.com/luk4z7/bible for the canonical source repository
//
// Copyright 2017 The Lucas Alves Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main bootstrap the program
package main

import (
	"database/sql"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/agtorre/gocolorize"
	_ "rsc.io/sqlite"
)

// @todo include other languages config
var (
	// bookConversion provide the conversion of names of books
	bookConversion = map[string][]string{
		"genesis": {
			"01O",
			"Gênesis",
		},
		"exodo": {
			"02O",
			"Êxodo",
		},
		"levitico": {
			"03O",
			"Levítico",
		},
		"numeros": {
			"04O",
			"Números",
		},
		"deuteronomio": {
			"05O",
			"Deuteronômio",
		},
		"josue": {
			"06O",
			"Josué",
		},
		"juizes": {
			"07O",
			"Juízes",
		},
		"rute": {
			"08O",
			"Rute",
		},
		"1samuel": {
			"09O",
			"1 Samuel",
		},
		"2samuel": {
			"10O",
			"2 Samuel",
		},
		"1reis": {
			"11O",
			"1 Reis",
		},
		"2reis": {
			"12O",
			"2 reis",
		},
		"1cronicas": {
			"13O",
			"1 Crônicas",
		},
		"2cronicas": {
			"14O",
			"2 Crônicas",
		},
		"esdras": {
			"15O",
			"Esdras",
		},
		"neemias": {
			"16O",
			"Neemias",
		},
		"ester": {
			"17O",
			"Ester",
		},
		"jo": {
			"18O",
			"Jó",
		},
		"salmos": {
			"19O",
			"Salmos",
		},
		"proverbios": {
			"20O",
			"Provérbios",
		},
		"eclesiastes": {
			"21O",
			"Eclesiastes",
		},
		"canticosalomao": {
			"22O",
			"Cântico de Salomão",
		},
		"isaias": {
			"23O",
			"Isaías",
		},
		"jeremias": {
			"24O",
			"Jeremias",
		},
		"lamentacoes": {
			"25O",
			"Lamentações",
		},
		"ezequiel": {
			"26O",
			"Ezequiel",
		},
		"daniel": {
			"27O",
			"Daniel",
		},
		"oseias": {
			"28O",
			"Oseias",
		},
		"joel": {
			"29O",
			"Joel",
		},
		"amos": {
			"30O",
			"Amós",
		},
		"obadias": {
			"31O",
			"Obadias",
		},
		"jonas": {
			"32O",
			"Jonas",
		},
		"miqueias": {
			"33O",
			"Miqueias",
		},
		"naum": {
			"34O",
			"Naum",
		},
		"habacuque": {
			"35O",
			"Habacuque",
		},
		"sofonias": {
			"36O",
			"Sofonias",
		},
		"ageu": {
			"37O",
			"Ageu",
		},
		"zacarias": {
			"38O",
			"Zacarias",
		},
		"malaquias": {
			"39O",
			"Malaquias",
		},
		"mateus": {
			"40N",
			"Mateus",
		},
		"marcos": {
			"41N",
			"Marcos",
		},
		"lucas": {
			"42N",
			"Lucas",
		},
		"joao": {
			"43N",
			"João",
		},
		"atos": {
			"44N",
			"Atos",
		},
		"romanos": {
			"45N",
			"Romanos",
		},
		"1corintios": {
			"46N",
			"1 Coríntios",
		},
		"2corintios": {
			"47N",
			"2 Coríntios",
		},
		"galatas": {
			"48N",
			"Gálatas",
		},
		"efesios": {
			"49N",
			"Efésios",
		},
		"filipenses": {
			"50N",
			"Filipenses",
		},
		"colossenses": {
			"51N",
			"Colossenses",
		},
		"1tessalonicenses": {
			"52N",
			"1 Tessalonicenses",
		},
		"2tessalonicenses": {
			"53N",
			"2 Tessalonicenses",
		},
		"1timoteo": {
			"54N",
			"1 Timóteo",
		},
		"2timoteo": {
			"55N",
			"2 Timóteo",
		},
		"tito": {
			"56N",
			"Tito",
		},
		"filemon": {
			"57N",
			"Filêmon",
		},
		"hebreus": {
			"58N",
			"Hebreus",
		},
		"tiago": {
			"59N",
			"Tiago",
		},
		"1pedro": {
			"60N",
			"1 Pedro",
		},
		"2pedro": {
			"61N",
			"2 Pedro",
		},
		"1joao": {
			"62N",
			"1 João",
		},
		"2joao": {
			"63N",
			"2 João",
		},
		"3joao": {
			"64N",
			"3 João",
		},
		"judas": {
			"65N",
			"Judas",
		},
		"apocalipse": {
			"66N",
			"Apocalipse",
		},
	}

	p           = fmt.Println
	logger      *log.Logger
	colorOutput func(v ...interface{}) string
	winsize     int
)

const (
	drive = "sqlite3"
)

func modulePath(name, version string) string {
	cache, ok := os.LookupEnv("GOMODCACHE")
	if !ok {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = build.Default.GOPATH
		}

		cache = path.Join(gopath, "pkg", "mod")
	}

	return path.Join(cache, name+"@"+version)
}

// Instance receive the instance of database sqlite
type Instance struct {
	db *sql.DB
}

// getDB return instance of database sqlite
func getDB() *sql.DB {
	db, err := sql.Open(drive, modulePath("github.com/luk4z7/bible", "v1.0.2")+"/bible.db")
	if err != nil {
		log.Fatalln("could not open database:", err)
	}

	return db
}

// Winsize receive the information about the terminal when executing
type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--index" {
		p("Index")
		putDisplay()
		return
	}
	if len(os.Args) < 3 {
		p("Name of book, chapter and verse are required")
		os.Exit(2)
	}
	nameBook := os.Args[1]
	if len(bookConversion[nameBook]) < 1 {
		p("There is no book with this name")
		os.Exit(2)
	}
	splitData := strings.Split(os.Args[2:][0], ":")

	if len(splitData) == 1 {
		p("Invalid input, correct pattern chapter:verse")
		os.Exit(2)
	}
	if ok := splitData[0]; ok == "" {
		p("Chapter number is required")
		os.Exit(2)
	}
	if ok := splitData[1]; ok == "" {
		p("Verse number is required")
		os.Exit(2)
	}
	chapter := splitData[0]
	splitVerse := strings.Split(splitData[1], "-")

	var verseBegin string
	var verseEnd string
	verseBegin = splitVerse[0]
	verseEnd = verseBegin

	if len(splitVerse) > 1 {
		if ok := splitVerse[1]; ok == "" {
			p("Invalid input, correct pattern chapter:versebegin-verseend")
			os.Exit(2)
		}
		verseEnd = splitVerse[1]
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	winsize = getWidth()
	color := gocolorize.NewColor("yellow")
	colorOutput = color.Paint
	logger = log.New(os.Stdout, colorOutput(" Holy bible  :: "), 0)

	p()
	title := "Book: " + bookConversion[nameBook][1] + "  Chapter: " + chapter + "  Verse: " + verseBegin + "-" + verseEnd
	newText := generateSpaces(" "+title, 18)
	logger.Println(colorOutput(newText))
	p()

	instance := Instance{}
	instance.db = getDB()
	defer instance.db.Close()

	rows, err := get(instance.db, bookConversion[nameBook][0], chapter, verseBegin, verseEnd)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var a, b int
		var c string
		err = rows.Scan(
			&a,
			&b,
			&c,
		)
		if err != nil {
			log.Fatalln("Could not scan row:", err)
			break
		}
		result := " " + colorOutput(strconv.Itoa(b)) + " " + c
		p(breakLine(result))
	}
	rows.Close()
}

type sortedMap struct {
	m map[string][]string
	s [][]string
	b []string
}

func sortedKeys(m map[string][]string) *sortedMap {
	sm := new(sortedMap)
	sm.m = m
	i := 0
	for key, j := range m {
		sm.b = append(sm.b, key)
		sm.s = append(sm.s, j)
		i++
	}
	return sm
}

// @todo create index for viewing all books
// putDisplay print the index of bible
func putDisplay() {}

// get return information about parameters received
func get(db *sql.DB, id, chapter, verseBegin, verseEnd string) (*sql.Rows, error) {
	query := `
		SELECT
			capitulo,
			versiculo,
			text
		FROM bible
		WHERE id = '` + id + `' AND
		capitulo = ` + chapter + ` AND
		versiculo >= ` + verseBegin + ` AND
		versiculo <= ` + verseEnd + `;
	`
	rows, err := db.Query(query)
	return rows, err
}

func getDistinctID(db *sql.DB) (*sql.Rows, error) {
	query := `SELECT DISTINCT id FROM bible ORDER BY id;`
	rows, err := db.Query(query)
	return rows, err
}

// generateSpaces add spaces for format output the text
func generateSpaces(str string, param int) string {
	length := (winsize - param) - len(str)
	data := []byte(str)
	for i := 0; i < length; i++ {
		data = append(data, '\u0020')
	}
	return string(data)
}

// @todo Improvement the broken line when the line is required break one more line
// @todo include hyphen in line breaks where the word will be separated
// breakLine breaks the rows for column generation
func breakLine(str string) string {
	half := int(winsize / 2)
	if len(str) > half {
		str = str[:half] + "\n\u0020" + str[half:]
	}
	return str
}

// getWidth obtain width about terminal on executing
func getWidth() int {
	ws := &Winsize{}
	retCode, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 {
		panic(err)
	}
	return int(ws.Col)
}
