package server

import (
	"errors"
	"fmt"
	"sync"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

type Node struct {
	sync.Mutex
	id int
}

type Arg struct{}
type Result struct {
	ID int
}

func (n *Node) Ping(args *Arg, r *Result) error {
	fmt.Println("in Ping")
	*r = Result{}
	return nil
}

func (n *Node) Start(args *Arg, res *Result) error {
	n.Lock()
	defer n.Unlock()

	n.id++
	res.ID = n.id
	return nil
}

func (n *Node) GetID(args *Arg, res *Result) error {
	res.ID = n.id
	return nil
}
