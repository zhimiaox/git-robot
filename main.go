package main

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/jinzhu/configor"
	cron "github.com/robfig/cron/v3"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const (
	configPath = "config.toml"
)

var (
	repository  *git.Repository
	worktree    *git.Worktree
	userStorage UserStorage
)

func main() {
	var err error
	err = configor.Load(Config, configPath)
	if err != nil {
		panic("配置加载失败")
	}
	repository, err = git.PlainOpen(Config.Work.GitPath)
	if err != nil {
		panic(err)
	}
	worktree, err = repository.Worktree()
	if err != nil {
		panic(err)
	}

	userStorage = NewMapStorage()

	StartWork()
	// 启动api接口
	StartApi()
}

func StartWork() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("1 0 0,12 * * ?", func() {
		for _, user := range userStorage.List() {
			err := GitWork(user)
			if err != nil {
				user.ErrCount++
			}
			if user.ErrCount > Config.Work.MaxErrCount {
				userStorage.Del(user.Sign())
			} else {
				userStorage.Set(user)
			}
		}
	})
	c.Start()
}

func GitWork(user *User) error {
	now := time.Now()
	user.LastRunTime = now
	day := now.Format("2006-01-02")
	ioutil.WriteFile(filepath.Join(Config.Work.GitPath, day), []byte(now.String()), 0666)
	worktree.Add("./")
	_, err := worktree.Commit("打卡器", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  user.User,
			Email: user.Email,
			When:  now,
		},
	})
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(user.DeployKeys)
	if err != nil {
		return err
	}
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
	repository.DeleteRemote("origin")
	remote, err := repository.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{user.RemoteURL},
	})
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err = remote.PushContext(ctx, &git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	return err
}
