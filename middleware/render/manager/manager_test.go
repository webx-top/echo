package manager

import (
	"os"
	"testing"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo/middleware/render/driver"
)

func TestManagerLoop(t *testing.T) {
	mgr := New()
	for i := 0; i < 3; i++ {
		testManager(t, mgr)
	}
}

func TestManager(t *testing.T) {
	mgr := New()
	testManager(t, mgr)
}

func testManager(t *testing.T, mgr driver.Manager) {
	dirs := []string{
		`./group_a/a`, `./group_b/b`,
	}
	changed := make(chan string)
	go func() {
		for t := range <-changed {
			println(`---------->`, t)
		}
	}()
	println(`~~~~~~~~~~~~>AddCallback`)
	// mgr.AddCallback(`./`, func(name, typ, event string) {
	// 	changed <- fmt.Sprintln(name, typ, event)
	// })
	defer log.Sync()
	log.GetLogger(`watcher`).SetLevel(`Debug`)
	var err error
	for _, dir := range dirs {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = mgr.AddWatchDir(dir)
		if err != nil {
			panic(err)
		}
	}

	println(`~~~~~~~~~~~~>Start`)
	mgr.Start()
	defer mgr.Close()

	err = os.WriteFile(`./group_a/a/test.log`, []byte(time.Now().String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
	//log.GetLogger(`watcher`).SetLevel(`Debug`)
	b, err := mgr.GetTemplate(`./group_a/a/test.log`)
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
	// for file := range mgr.(*Manager).caches {
	// 	println(`===>:`, file)
	// }
	// mgr.(*Manager).caches.Range(func(file string, _ []byte) bool {
	// 	println(`===>:`, file)
	// 	return true
	// })
	time.Sleep(2 * time.Second)

	err = os.WriteFile(`./group_b/b/test.log`, []byte(time.Now().String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	err = os.WriteFile(`./group_a/a/test.log`, []byte(time.Now().String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	err = os.WriteFile(`./group_a/test.log`, []byte(time.Now().String()), os.ModePerm)
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
