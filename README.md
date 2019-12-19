# golang tour
Go files and solutions to Russ Cox 2010 course. Good course to learn various aspects of Go.

TODO : Not all the code is mine. , try urself first..

## Build instructions
Pre-requisites: You should have Go installed and GOPATH environment variable set.
For linux

```sh
$ cd $GOPATH/src
$ git clone https://github.com/yogisinha/tour.git
$ go install
```

It will put the "tour" binary in your $GOPATH/bin directory

## Things to Note
Most of the following exercises runs a server at port 4000. If you have to change that on your machine replace that 
port number with your port number in following instructions. You can end the server by pressing Ctrl+C to end each exercise.

Following is the list of exercises and how to run them:

### Nuts and Bolts

#### Square root
-sqrt option accepts the number for which square root to be printed
```sh
tour -sqrt=66
```
#### WordCount
-wcf option runs the WordCount exercise. Run the binary with following option:
```sh
tour -wcf
```
Now go to localhost:4000 in your browser and start typing in the text-area. It will print the words and their ocurrences interactively.

#### Slices
-slices options runs the Slices exercise which accepts some drawing functions. It displays the slice data interpreted as grayscale pixel values. Those functions produces some art. Run the binary with following option:
```sh
tour -slices="X+Y"
```
Now go to localhost:4000 in your browser and you will see some art. Possible values for this option is 
1. "X+Y"
2. "X*Y"
3. "(X+Y)/2"
4. "X^Y"

### Interfaces and Web Servers

#### Hello World 2.0
-helloworld2 option runs the Hello World 2 exercise. Run the binary with following option:
```sh
tour -helloworld2
```
Now go to localhost:4000 with request parameter "name" and it will respond with "Hello, name" 

#### Methods
-methods option runs the Methods exercise. Run the binary with following option:
```sh
tour -methods
```
It just shows how to expose a type as a server by implementing ServeHTTP method on those types.
Go to localhost:4000/string and localhost:4000/struct in your browser.

#### Image Viewer
-imgver option runs the Image Viewer exercise. It is the same concept as Slices exercise but instead of 
returning [][]unit8 slices, it returns image.Image type which then displayed on browser. Run the binary with following option:
```sh
tour -imgvwer="X+Y"
```
Go to localhost:4000 in your browser and you will see some art. Possible values for this option is
same as in Slices exercise.

#### PNG Encoding 
-pngencode option runs the PNG Encoding exercise. Run the binary with following option:
```sh
tour -pngencode
```
It shows how to define youur own image type and expose it as a http end point by implementing ServeHTTP method on it.
You can specify width and height parameter as request parameters. for e.g. localhost:4000?x=200&y=200 and it will produce
png image of that size. Currently, it just implments one function to produce the pixel colors. More functions as in Slices
exercise can be implemented.

#### Mandelbrot Set Viewer
-mandelbrot option runs the Mandelbrot Set Viewer exercise. Run the binary with following option:
```sh
tour -mandelbrot
```
Go to localhost:4000 and it will generate Mandelbrot image on the browser. You can pan in the viewer
by dragging with the mouse and zoom by using the scroll wheel or by clicking + and −

#### Julia Set Viewer
-julia option runs the Julia Set Viewer exercise. Run the binary with following option:
```sh
tour -julia
```
Go to http://localhost:4000/juliaviewer and it will generate Julia set image on the browser. You can pan in the viewer
by dragging with the mouse and zoom by using the scroll wheel or by clicking + and −

### Concurrency and Channels

#### Equivalent Binary Trees
It is implemented in "func Same(t1, t2 *tree.Tree) bool" function to check whether 2 Binary search Trees are equivalent or not.
This exercise demonstrates the use of goroutines and channels.

#### Web Logging
-weblog option runs the Web Logging exercise. Run the binary with following option:
```sh
tour -weblog
```
It sends an infinite sequence (say, 1, 2, 3, 4, 5, ...; ) of Log messages on browser, sleeping for a second
after each send. Load http://localhost:4000/ and you should see the log messages. 

#### Chat Room
-chatex option runs all the Chat related exercises (Chat Room, User Lists, Julia Set Bot). Run the binary with following option:
```sh
tour -chatex
```
Open two or more browser windows on http://localhost:4000 and you can 
1. Join the chat and send messages from one browser window to another. Messages will be broadcasted to all the windows.
2. When a user joins (or exits), a message will be displayed with a list of all the users, formatted like:
[+newUser, oldUser1, oldUser2]
[−userWhoLeft, remainingUser1, remainingUser2]
3. For the julia bot - This feature respond with a julia image for messages like julia: 0.285+0.01i

### Networking with RPCs

#### File downloads
Exercises under this topic was mainly about building a distributed file download like system. 
This exercise will stream the images in small fragments to your browser either in sequential 
or distributed fashion.

This exercise mainly starts with -rpc option:
```sh
tour -rpc -image=<some image path> -x=<number> -y=<number> -random 
```
Below is the description of options to run this exercise:
1. -image : accepts the any image path on your system
2. -x and -y : accepts width and height of block which the program will take and break the image in block of those 
dimension and send it to browser. (default for x and y is 50 and 50)
3. -random : whether to send the image blocks in random order or not (default is false)
4. -mode : sequential or distributed. In sequential, image blocks are sent by only one goroutine, In distributed, image 
blocks are distributed among multiple goroutines listening on tcp connections and each goroutine is responsible for 
certain no of blocks. Possible values for this flag is s and d. (default is s (sequential))


