package main

import "github.com/map34/simple_tftp/tftputils"

func main() {
	err := tftputils.SpawnServeSession()
	if err != nil {
		panic(err)
	}
}
