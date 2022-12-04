# go-then

Inspired by Javascript's [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise), go-then is a library that attempts to allow non-blocking execution of codes using goroutines, channels and wait groups.

Sample Javascript Promise

```
let promise = new Promise((resolve, reject) => {
  setTimeout(() => {
    resolve("world");
  }, 5000);
});

promise.then((v) => {
  console.log(v);
}).catch((e) => {
  console.log(v);
});

console.log("hello");
```

go-then's equivalent of the above Javascript Promise:

```
package main

import (
	"log"
	"time"

	promise "github.com/mohamadHarith/go-then"
)

func main() {
	promise := promise.New(&promise.Config{TimeOutInSecs: 60})
	defer promise.Wait() // wait for the promise to execute

	promise.Execute(func() {
		// wait for 5 mins before resolving
		time.Sleep(time.Second * 5)
		promise.Resolve("world")

	}).Then(func(i any) {
		log.Println(i)

	}).Catch(func(err error) {
		log.Println(err)
	})

	log.Println("hello")
}
```
