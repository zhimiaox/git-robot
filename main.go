package main

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron/v3"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const (
	MaxErrCount = 5
	GitPath     = "D://git/dk"
)

var (
	repository  *git.Repository
	worktree    *git.Worktree
	userStorage UserStorage
)

func init() {
	var err error
	repository, err = git.PlainOpen(GitPath)
	if err != nil {
		panic(err)
	}
	worktree, err = repository.Worktree()
	if err != nil {
		panic(err)
	}
	userStorage = NewMapStorage()
}

func main() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 * * * * ?", func() {
		for _, user := range userStorage.List() {
			err := GitWork(user)
			if err != nil {
				user.ErrCount++
			}
			if user.ErrCount > MaxErrCount {
				userStorage.Del(user.Sign())
			} else {
				userStorage.Set(user)
			}
		}
	})
	c.Start()
	r := gin.Default()
	gr := r.Group("/api")
	gr.POST("/user", Set)
	gr.GET("/user/search", Search)
	gr.DELETE("/user", Del)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
func Search(c *gin.Context) {

}

func Set(c *gin.Context) {

}

func Del(c *gin.Context) {

}
func GitWork(user *User) error {
	now := time.Now()
	user.LastRunTime = now
	day := now.Format("2006-01-02")
	ioutil.WriteFile(filepath.Join(GitPath, day), []byte(now.String()), 0666)
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
