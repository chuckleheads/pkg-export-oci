package account

import (
	"html/template"
	"os"
	"path/filepath"
)

const PASSWD = "defaults/etc/passwd"
const GROUP = "defaults/etc/group"

type Account struct {
	User  string
	Group string
}

func New(user string, group string) Account {
	return Account{
		User:  user,
		Group: group,
	}
}

func (a *Account) Write(fsroot string) {
	tmpl, err := template.ParseFiles(PASSWD)
	f, err := os.Create(filepath.Join(fsroot, "etc", "passwd"))
	check(err)
	tmpl.Execute(f, a)
	f.Close()
	tmpl, err = template.ParseFiles(GROUP)
	f, err = os.Create(filepath.Join(fsroot, "etc", "group"))
	check(err)
	tmpl.Execute(f, a)
	f.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
