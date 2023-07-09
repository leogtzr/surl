package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type memoryDB struct {
	db            map[int]string
	autoIncrement int
}

// InMemoryURLDAOImpl ...
type InMemoryURLDAOImpl struct {
	DB *memoryDB
}

// InMemoryUserDAOImpl ...
type InMemoryUserDAOImpl struct {
	db       map[string]UserInMemory
	rndIDGen randGenSrc
}

// StatsDAOMemoryImpl ...
type StatsDAOMemoryImpl struct {
	db map[int][]StatsInMemory
}

func (im InMemoryURLDAOImpl) save(url URL) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	im.DB.autoIncrement++
	id := im.DB.autoIncrement
	im.DB.db[id] = url.URL

	return id, nil
}

func (im InMemoryURLDAOImpl) findAllByUser() ([]URLStat, error) {
	// shortID:int, url:string
	var urls []URLStat

	// dummy impl...
	for shortID, url := range im.DB.db {
		urls = append(urls, URLStat{
			ShortID: shortID,
			Url:     url,
		})
	}

	return urls, nil
}

func (im InMemoryURLDAOImpl) findByID(id int) (URL, error) {
	u, found := im.DB.db[id]
	if found {
		url := URL{
			URL: u,
		}

		return url, nil
	}

	return URL{}, errorURLNotFound(id)
}

func (im InMemoryURLDAOImpl) update(id int, oldURL, newURL URL) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := im.DB.db[id]; !ok {
		return id, errorURLNotFound(id)
	}

	newID := shortURLToID(newURL.URL, chars)
	url := im.DB.db[id]

	im.DB.db[newID] = url
	delete(im.DB.db, id)

	return newID, nil
}

func (dao InMemoryUserDAOImpl) addUser(username, password string) (interface{}, error) {
	hashPassword := password

	id := dao.rndIDGen.Uint64()

	newUser := UserInMemory{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		User:      username,
		Password:  hashPassword,
	}

	dao.db[username] = newUser

	return newUser, nil
}

func (dao InMemoryUserDAOImpl) userExists(username string) (bool, error) {
	_, exists := dao.db[username]
	if !exists {
		return false, nil
	}

	return true, nil
}

func (dao InMemoryUserDAOImpl) findByUsername(username string) (interface{}, error) {
	user, exists := dao.db[username]
	if !exists {
		return UserInMemory{}, errorUserNotFound(username)
	}

	return user, nil
}

func (dao InMemoryUserDAOImpl) validateUserAndPassword(username, password string) (bool, error) {
	user, err := dao.findByUsername(username)
	if err != nil {
		return false, err
	}

	u, ok := user.(UserInMemory)
	if !ok {
		return false, errorIncompatibleTypes()
	}

	hashFromDatabase := []byte(u.Password)
	if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func (dao InMemoryUserDAOImpl) findAll() ([]interface{}, error) {
	var users []interface{}

	for _, v := range dao.db {
		users = append(users, v)
	}

	return users, nil
}

func (dao StatsDAOMemoryImpl) save(shortURL string, headers *map[string][]string) (int, error) {
	urlShortID := shortURLToID(shortURL, chars)

	stat := StatsInMemory{
		CreatedAt: time.Now(),
		ShortID:   urlShortID,
		Headers:   *headers,
	}

	dao.db[urlShortID] = append(dao.db[urlShortID], stat)

	fmt.Println("debug save")
	fmt.Println(dao.db[urlShortID])

	return 0, nil
}

func (dao StatsDAOMemoryImpl) findByShortID(shortID int) ([]interface{}, error) {
	return []interface{}{}, nil
}

func (dao StatsDAOMemoryImpl) findAll() (map[int][]StatsInMemory, error) {
	fmt.Println("debug -")
	fmt.Println(dao.db)

	return dao.db, nil
}