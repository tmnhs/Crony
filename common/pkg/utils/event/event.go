package event

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

const (
	EXIT = "exit"
	WAIT = "wait"
)

var (
	Events = make(map[string][]func(interface{}), 2)
)

func OnEvent(name string, fs ...func(interface{})) error {
	evs, ok := Events[name]
	if !ok {
		evs = make([]func(interface{}), 0, len(fs))
	}

	for _, f := range fs {
		if f == nil {
			continue
		}

		fp := reflect.ValueOf(f).Pointer()
		for i := 0; i < len(evs); i++ {
			if reflect.ValueOf(evs[i]).Pointer() == fp {
				return fmt.Errorf("func[%v] already exists in event[%s]", fp, name)
			}
		}
		evs = append(evs, f)
	}
	Events[name] = evs
	return nil
}

func EmitEvent(name string, arg interface{}) {
	evs, ok := Events[name]
	if !ok {
		return
	}

	for _, f := range evs {
		f(arg)
	}
}

func EmitAllEvent(arg interface{}) {
	for _, fs := range Events {
		for _, f := range fs {
			f(arg)
		}
	}
	return
}

func OffEvent(name string, f func(interface{})) error {
	evs, ok := Events[name]
	if !ok || len(evs) == 0 {
		return fmt.Errorf("envet[%s] doesn't have any funcs", name)
	}

	fp := reflect.ValueOf(f).Pointer()
	for i := 0; i < len(evs); i++ {
		if reflect.ValueOf(evs[i]).Pointer() == fp {
			evs = append(evs[:i], evs[i+1:]...)
			Events[name] = evs
			return nil
		}
	}

	return fmt.Errorf("%v func dones't exist in event[%s]", fp, name)
}

func OffAllEvent(name string) error {
	Events[name] = nil
	return nil
}

// 等待信号
// 如果信号参数为空，则会等待常见的终止信号
// SIGINT 2 A 键盘中断（如break键被按下）
// SIGTERM 15 A 终止信号
func WaitEvent(sig ...os.Signal) os.Signal {
	c := make(chan os.Signal, 1)
	if len(sig) == 0 {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	} else {
		signal.Notify(c, sig...)
	}
	return <-c
}
