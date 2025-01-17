package main

import "fmt"
import "sync"

type NumObj struct{
    counter int
    mutex sync.Mutex
}

var waitGroup sync.WaitGroup
var numObj = NumObj{}

func main(){
    for i:= 0; i< 1000; i++ {
        waitGroup.Add(1)
        go incrementCounter(i, &numObj)
    }
    waitGroup.Wait()
    fmt.Println("final counter is ",numObj.counter)
}

func incrementCounter(number int, numObj* NumObj){
    defer waitGroup.Done()
    numObj.mutex.Lock()
    defer numObj.mutex.Unlock()
    fmt.Println("adding up ",number)
    if number % 2 == 0 {
        fmt.Println("number is divided by 2 and the number is ",number)
        return
    }
    numObj.counter += number
    fmt.Println("counter now is ",numObj.counter)
}
