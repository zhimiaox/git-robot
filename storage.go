package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type UserStorage interface {
	Set(user *User) error
	Del(key string) error
	Get(key string) *User
	List() []*User
}

// ---- 存储区
type MapStorage struct {
	sync.RWMutex
	Map        map[string]*User
	SaveSignal chan int
}

func NewMapStorage() *MapStorage {
	st := new(MapStorage)
	st.Map = make(map[string]*User)
	st.SaveSignal = make(chan int, 3)
	if Config.MapStorage.FilePath != "" {
		if jsonData, err := ioutil.ReadFile(Config.MapStorage.FilePath); err == nil {
			json.Unmarshal(jsonData, &st.Map)
		}
		go func(m *MapStorage) {
			for {
				<-m.SaveSignal
				st.RLock()
				data, err := json.Marshal(&m.Map)
				st.RUnlock()
				if err != nil {
					continue
				}
				ioutil.WriteFile(Config.MapStorage.FilePath, data, 0666)
			}
		}(st)
	}
	return st
}

func (m *MapStorage) Set(user *User) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Map[user.Sign()]; ok {
		return fmt.Errorf("已存在")
	}
	m.Map[user.Sign()] = user
	// 触发持久化
	m.SaveSignal <- 1
	return nil
}

func (m *MapStorage) Del(key string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.Map, key)
	m.SaveSignal <- 1
	return nil
}

func (m *MapStorage) Get(key string) *User {
	m.RLock()
	defer m.RUnlock()
	return m.Map[key]
}

func (m *MapStorage) List() []*User {
	m.RLock()
	defer m.RUnlock()
	list := make([]*User, 0)
	for _, user := range m.Map {
		list = append(list, user)
	}
	return list
}
