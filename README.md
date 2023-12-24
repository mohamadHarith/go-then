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
	"context"
	"fmt"
	"time"

	promise "github.com/mohamadHarith/go-then"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	p := promise.New(ctx, func(resolve promise.Resolver, reject promise.Rejector) {

		// simulate async task
		time.Sleep(time.Second * 5)

		resolve("world")

	}).Then(func(resp any) {

		fmt.Println(resp)

	}).Catch(func(err error) {
		fmt.Println("err: ", err)
	})

	// wait for the promise to finish executing
	// before exiting the main thread
	defer p.Wait()

	fmt.Println("hello")
}

// output:
// hello
// world
```

## Features
- Javascript promise like syntax.
- Non-blocking execution.

Try it on [playground](https://go.dev/play/p/GBG6AyJrZc4).
