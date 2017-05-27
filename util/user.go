package util

import (
	"log"
	"os/user"
	"path"
)

func GetUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func GetUserDefaultSSHPublicKeyPath() string {
	home := GetUserHomeDir()
	return path.Join(home, ".ssh", "id_rsa.pub")
}
