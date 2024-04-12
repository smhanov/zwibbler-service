package main

import (
	"log"
	"math/rand"

	"github.com/smhanov/zwibserve"
)

// simulate a database that randomly looses writes, but silently

type UnreliableDB struct {
	zwibserve.DocumentDB
	failRate int
}

func createUnreliableDB(db zwibserve.DocumentDB, failRate int) *UnreliableDB {
	return &UnreliableDB{db, failRate}
}

func (db *UnreliableDB) AppendDocument(docID string, oldLength uint64, newData []byte) (uint64, error) {
	if db.failRate > 0 && rand.Intn(db.failRate) == 0 {
		log.Printf("Silently dropping write to %s", docID)
		return oldLength + uint64(len(newData)), nil
	}
	log.Printf("Not silently dropping write")
	return db.DocumentDB.AppendDocument(docID, oldLength, newData)
}
