package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/admpub/log"
)

func TestManagerLoop(t *testing.T) {
	for i := 0; i < 3; i++ {
		TestManager(t)
	}
}

func TestManager(t *testing.T) {
	dirs := []string{
		`./group_a/a`, `./group_b/b`,
	}
	changed := make(chan string)
	Default.AddCallback(`./`, func(name, typ, event string) {
		changed <- fmt.Sprintln(name, typ, event)
	})
	go func() {
		for {
			select {
			case t := <-changed:
				println(`---------->`, t)
			}
		}
	}()
	defer log.Sync()
	log.GetLogger(`watcher`).SetLevel(`Debug`)
	var err error
	for _, dir := range dirs {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = Default.AddWatchDir(dir)
		if err != nil {
			panic(err)
		}
	}

	Default.Start()
	defer Default.Close()

	err = ioutil.WriteFile(`./group_a/a/test.log`, []byte(time.Now().String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
	//log.GetLogger(`watcher`).SetLevel(`Debug`)
	b, err := Default.GetTemplate(`./group_a/a/test.log`)
	if err != nil {
		panic(err)
	}
	println(`./group_a/a/test.log:`, string(b))
	/*
		println(`cancel ./group_a`)
		Default.CancelWatchDir(`./group_a`)
		println(`cancel ./group_b`)
		Default.CancelWatchDir(`./group_b`)
	*/
	for file := range Default.(*Manager).caches {
		println(`===>:`, file)
	}
	time.Sleep(2 * time.Second)

	err = ioutil.WriteFile(`./group_b/b/test.log`, []byte(time.Now().String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	err = ioutil.WriteFile(`./group_a/a/test.log`, []byte(time.Now().String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	err = ioutil.WriteFile(`./group_a/test.log`, []byte(time.Now().String()), os.ModePerm)
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)
	for _, dir := range dirs {
		err = os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}
	os.RemoveAll(`./group_a`)
	os.RemoveAll(`./group_b`)
}
