package main

type Msg interface{}

type Cmd func() Msg

type Program struct {
	ch       chan Msg
	done     chan bool
	finished chan struct{}
}

type (
	initMsg struct{}
	secMsg  struct{}
	doneMsg struct{}
)

func (p *Program) eventloop(cmd chan Msg) {
	for {
		select {
		case <-p.done:
			println("done")
			cmd <- doneMsg{}
			return
		case msg := <-p.ch:
			switch msg.(type) {
			case initMsg:
				println("initCmd")
				go func() {
					p.ch <- secMsg{}
				}()
			case secMsg:
				println("secMsg")
				go func() {
					p.done <- true
				}()
			}
		}
	}
}

type channelHandlers []chan struct{}

func main() {
	c := make(chan Msg)
	done := make(chan bool)
	p := Program{ch: c, done: done}
	cmd := make(chan Msg)
	go func() {
		select {
		case <-cmd:
			println("first cmd")
		}
	}()
	go func() {
		c <- initMsg{}
	}()
	p.eventloop(cmd)
	println("exit")
}
