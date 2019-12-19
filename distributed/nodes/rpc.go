package nodes

import (
	"image"
	"log"
	"math"
	"net/rpc"
)

type StartResult struct {
	ID string
}

type Address struct {
	Addr string
}

func Call(connectAddress, method string, args, reply interface{}) {
	//fmt.Println("In call...", connectAddress, method, args)

	client, err := rpc.Dial("tcp", connectAddress)

	if err != nil {
		log.Fatal("In Call dialing:", err)
	}
	defer client.Close()

	// Synchronous call
	err = client.Call(method, args, reply)
	if err != nil {
		log.Printf("Error in RPC call : method %s  %v:", method, err)
	}

}

// GetBlocks gets the no of blocks which we will get by breaking an image
// of X x Y dimension (Width: X, Height: Y) where each block is of dimension mxn
func GetBlocks(X, Y int, m, n int) []image.Rectangle {
	var xsteps, ysteps int
	if X%m == 0 {
		xsteps = X / m
	} else {
		xsteps = X/m + 1
	}

	if Y%n == 0 {
		ysteps = Y / n
	} else {
		ysteps = Y/n + 1
	}

	if m >= X {
		xsteps = 1
	}
	if n >= Y {
		ysteps = 1
	}

	var x0, y0 int
	var blocks []image.Rectangle
	for y := 0; y < ysteps; y++ {
		for x := 0; x < xsteps; x++ {
			frstRect := image.Rect(x0, y0, x0+m, y0+n)
			secondRect := image.Rect(x0, y0,
				int(math.Min(float64(x0+m), float64(X))),
				int(math.Min(float64(y0+n), float64(Y))))

			blocks = append(blocks, frstRect.Intersect(secondRect))
			x0 += m
		}

		x0 = 0
		y0 += n
	}

	return blocks

}
