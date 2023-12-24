package go_then

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestPromiseResolve(t *testing.T) {
	p := New(
		context.Background(),
		func(resolve Resolver, reject Rejector) {
			time.Sleep(time.Second * 5)
			resolve("world")
		}).Then(func(i any) {
		if i != "world" {
			t.Fatalf("data mismatch expected world but got %v ", i)
			return
		}
	})
	defer p.Wait()
}

func TestPromiseReject(t *testing.T) {
	p := New(
		context.Background(),
		func(resolve Resolver, reject Rejector) {
			reject(errors.New("test reject"))
		}).Catch(func(err error) {
		if err.Error() != "test reject" {
			t.Fatalf("err mismatch expected 'test reject' but got %v ", err)
			return
		}
	})
	defer p.Wait()
}

func TestPromiseWithContextTimeout(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	p := New(
		ctx,
		func(resolve Resolver, reject Rejector) {
			time.Sleep(time.Second * 6)
			resolve("world")
		}).Then(func(i any) {
		if i != nil {
			t.Fatal("i should be nil")
			return
		}
	}).Catch(func(err error) {
		if err == nil {
			t.Fatal("err should not be nil")
			return
		}
	})
	defer p.Wait()
}

func TestPromiseWithCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	p := New(
		ctx,
		func(resolve Resolver, reject Rejector) {
			time.Sleep(time.Second * 6)
			resolve("world")
		}).Then(func(i any) {
		if i != nil {
			t.Fatal("i should be nil")
			return
		}
	}).Catch(func(err error) {
		if err == nil {
			t.Fatal("err should not be nil")
			return
		}
	})
	defer p.Wait()

	time.Sleep(time.Second * 2)
	cancel()
}

func TestPromiseWithNetworkCall(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	p := New(ctx, func(resolve Resolver, reject Rejector) {
		resp, err := http.Get("https://catfact.ninja/fact")
		if err != nil {
			reject(err)
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			reject(err)
			return
		}
		resolve(body)
	}).Then(func(i any) {
		resp, ok := i.([]byte)
		if !ok {
			t.Fatal("type assertion error")
			return
		}
		t.Log(string(resp))
	}).Catch(func(err error) {
		t.Fatal(err)
	})
	defer p.Wait()

	t.Log("fetching cat data")
}
