package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/go-humble/locstor"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jQuery = jquery.NewJQuery
var chance = 10
var injectees []string
var injecteesSave []string

func start(inputInjects []string) {
	injectees = make([]string, len(inputInjects))
	copy(injectees, inputInjects)
	injecteesSave = make([]string, len(injectees))
	copy(injecteesSave, injectees)
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Save("original", injecteesSave); err != nil {
		println("Couldn't save original!")
	}
	load()
	if determine() {
		setMessage(injectees[0])
		injectees = append(injectees[:0], injectees[0+1:]...)
	}
	save()
}

func main() {
	js.Global.Set("ti", map[string]interface{}{
		"start": start,
	})
}

func determine() (returnVal bool) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randNum := r1.Intn(chance)
	if randNum == 0 {
		println("Showing \"injectee\"!")
		returnVal = true
	} else {
		println("Chances dictated " + strconv.Itoa(randNum) + " where 0 was needed. Chance is " + strconv.Itoa(chance))
		returnVal = false
		chance--
	}
	return
}

func setMessage(message string) {
	jQuery(".mainContent").SetText(message)
}

func save() {
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Delete("chance"); err != nil {
		println("Couldn't delete chance!")
	}
	if err := store.Save("chance", chance); err != nil {
		println("Couldn't save chance!")
	}
	if err := store.Save("injectees", injectees); err != nil {
		println("Couldn't save injectees!")
	}
}

func load() {
	var original []string
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Find("chance", &chance); err != nil {
		println("Couldn't load chance!", err)
	}
	if err := store.Find("injectees", &injectees); err != nil {
		println("Couldn't load injectees!", err)
	}
	if err := store.Find("original", &original); err != nil {
		println("Couldn't load original!", err)
	} else {
		if !testEq(original, injecteesSave) {
			println(original, injecteesSave)
			println("Injectees has changed! Resetting original and resetting all values to defaults")
			//reset()
		}
	}
}

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func reset() {
	store := locstor.NewDataStore(locstor.JSONEncoding)
	if err := store.Delete("original"); err != nil {
		println("Couldn't delete original!")
	}
	if err := store.Save("original", injecteesSave); err != nil {
		println("Couldn't save original!")
	}
	chance = 10
	save()
}
