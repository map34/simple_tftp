package main

import "bitbucket.org/adrian_prananda/simple_tftp/tftputils"

func main() {
	err := tftputils.SpawnServeSession()
	if err != nil {
		panic(err)
	}
}
