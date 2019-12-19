# golang tour
Go files and solutions to Russ Cox 2010 course. Good course to learn various aspects of Go.

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
Now load the url localhost:4000 in your browser and start typing in the text-area. It will print the words and their ocurrences interactively.

#### Slices
-slices options runs the Slices exercise which accepts some drawing functions. Thos functions produces some art. Run the binary with following option:
```sh
tour -slices="X+Y"
```
Now load the url localhost:4000 in your browser and you will see some art.







tour has following cmd line options :

       
        

