# go-then

Inspired by Javascript's [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise), go-then is a library that attempts to allow non-blocking execution of codes using goroutines, channels and wait groups.

Sample Javascript Promise

```javascript
const p = new Promise((resolve, reject)=>{
	setTimeOut(()=>{
		resolve("hola")
	}, 5000)

p.then((v)=>{
	console.log("resolved value=> ", v)
}).catch((e)=>{
	console.log("caught error=> ", v)
})
})

console.log("hi")

// prints in the below order:
// hi
// resolved value=>  hola
```

go-then's equivalent of the above Javascript Promise:

```golang
	promise1 := promise.New()
	defer promise1.Wait() // wait for the promise to execute

	promise1.Execute(func() {
		// wait for 5 mins before resolving
		time.Sleep(time.Second * 5)
		promise1.Resolve("hola")

	}).Then(func(i any) {
		log.Println("resolved value=> ", i)

	}).Catch(func(err error) {
		log.Println("caught err=> ", err)
	})

	log.Println("hi")

	// prints in the below order:
	// hi
	// resolved value=>  hola
```
