# go-then

API

```
getIp := promise.New(func(resolve func(string), reject func(error)) {
		ip, err := fetchIP()
		if err != nil {
			reject(err)
			return
		}
		resolve(ip)
})

getIp().then(func(){

})
```
