# go-then

Inspired by Javascript's [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise), go-then is a library that attempts to allow non-blocking execution of codes using goroutines, channels and wait groups.

This repository is under active development. Any contribution is highly appreciated. ^_^

Sample Javascript Promise

```javascript
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

// output:
// hello
// world
```

go-then's equivalent of the above Javascript Promise:

```golang
package main

import (
	"log"
	"time"

	promise "github.com/mohamadHarith/go-then"
)

func main() {
	ctx := context.Background()
	promise.New(ctx, func(resolve promise.Resolver, reject promise.Rejector){
		// some work
		time.Sleep(time.Second*5)
		resolve("world")
	}).Then(func(i any) {
		log.Println(i)

	}).Catch(func(err error) {
		log.Println(err)
	})

	log.Println("hello")
}

// output:
// hello
// world
```

Try it on [playground](https://go.dev/play/p/GevioARAp-S).
