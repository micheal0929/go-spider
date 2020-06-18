package common

import "log"

func SafeGo(f func()) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("err :%v", err)
		}
	}()
	f()
}
