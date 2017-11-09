package main

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

type store struct {
	db       *badger.DB
	isClosed bool
	keys     storeKeys
	secret   []byte
}

type storeKeys struct {
	appSecret  []byte
	plexToken  []byte
	plexServer []byte
}

func initDataStore(dirName string) (store, error) {
	var db store

	if isVerbose {
		fmt.Println("checking if our datastore exists in the home directory at:", dirName)
	}

	// create a directory for our database
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if isVerbose {
			fmt.Println("creating directory because it doesn't exist")
		}

		if err := os.Mkdir(dirName, os.ModePerm); err != nil {
			return db, err
		}
	} else if !os.IsNotExist(err) && isVerbose {
		fmt.Println("datastore already exits")
	}

	options := badger.DefaultOptions

	options.Dir = dirName
	options.ValueDir = dirName

	kvStore, err := badger.Open(options)

	if err != nil {
		return db, err
	}

	if isVerbose {
		fmt.Println("successfully opened data store")
	}

	db.db = kvStore
	db.keys = storeKeys{
		appSecret:  []byte("app-secret"),
		plexToken:  []byte("plex-token"),
		plexServer: []byte("plex-server"),
	}

	return db, nil
}

func (s store) Close() {
	if s.isClosed {
		fmt.Println("datastore already closed")
		return
	}

	if err := s.db.Close(); err != nil {
		fmt.Printf("datastore failed to closed: %v\n", err)
	}

	if isVerbose {
		fmt.Println("datastore is closed")
	}

	s.isClosed = true
}

func (s store) getSecret() []byte {
	var secret []byte

	// an error is returned when the key is not found
	// so just return an empty secret
	s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.appSecret)

		if err != nil {
			return err
		}

		_secret, err := item.Value()

		if err != nil {
			return err
		}

		secret = _secret

		return nil
	})

	return secret
}

func (s store) saveSecret(secret []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.appSecret, secret)
	})
}

func (s store) getPlexToken() (string, error) {
	var plexToken string

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.plexToken)

		if err != nil {
			return err
		}

		tokenHash, err := item.Value()

		if err != nil {
			return err
		}

		_plexToken, err := decrypt(s.secret, string(tokenHash))

		if err != nil {
			if isVerbose {
				fmt.Println("token decryption failed")
			}
			return err
		}

		plexToken = _plexToken

		return nil
	}); err != nil {
		return plexToken, err
	}

	if isVerbose {
		fmt.Printf("Your plex token is %s\n", plexToken)
	}

	return plexToken, nil
}

func (s store) savePlexToken(token string) error {
	tokenHash, err := encrypt(s.secret, token)

	if err != nil {
		return err
	}

	if isVerbose {
		fmt.Printf("your plex token hash: %s\n", string(tokenHash))
	}

	if err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.plexToken, []byte(tokenHash))
	}); err != nil {
		return err
	}

	if isVerbose {
		fmt.Println("saved token hash to store")
	}

	return nil
}

func (s store) getPlexServer() (server, error) {
	var plexServer server

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.plexServer)

		if err != nil {
			return err
		}

		serializedServer, err := item.Value()

		if err != nil {
			return err
		}

		_plexServer, err := unserializeServer(serializedServer)

		if err != nil {
			return err
		}

		plexServer = _plexServer

		return nil
	})

	return plexServer, err
}

func (s store) savePlexServer(plexServer server) error {
	serializedServer, err := plexServer.Serialize()
	if err != nil {
		return err
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.plexServer, serializedServer)
	})
}
