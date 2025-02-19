// cSpell:disable
package forwarder

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
)

type testFrdConfig struct {
	Port       int
	Forwarding map[string]string
}

type requestPool struct {
	requestPooltruth  []string
	requestPoolunreal []string
	requestPoolmulti  []string
	truth             bool
	unreal            bool
	multi             bool
}

func sendRequest(pool []string, t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()
	numTest := 50
	num := 0
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		str := pool[rand.Intn(3)]
		t.Logf("str ----> %v", str)
		i := rand.Intn(2)
		t.Logf("i ----> %d", i)
		switch {
		case num == numTest:
			return
		case i == 0:
			result, err := forwardToDNSServer(ctx, []byte(str), "8.8.8.8:53")
			if err != nil {
				t.Logf("err ----> %s", err)
			}
			t.Logf("respons %s\n", string(result))
			if string(result) == "" && err == nil {
				t.Error("[WARN] Became empty response")
			}
		case i != 0:
			result, err := forwardToDNSServer(ctx, []byte(str), "1.1.1.1:53")
			if string(result) == "" && err == nil {
				t.Error("[WARN] Became empty response")
			}
		}
		num++
	}
}

func TestForwardToDNSServer(t *testing.T) {
	wg := &sync.WaitGroup{}
	var request = requestPool{
		requestPooltruth:  []string{"mail.ru", "google.ru", "yandex.com"},
		requestPoolunreal: []string{"i", "wnat", "to", "break", "you", "!"},
		requestPoolmulti:  []string{"mail.ru", "google.ru", "yandex.com", "i", "wnat", "to", "break", "you", "!"},
		truth:             false,
		unreal:            false,
		multi:             false,
	}

	for i := 0; i < 2; i++ {
		switch {
		case request.truth:
			request.truth = true
			wg.Add(1)
			go sendRequest(request.requestPooltruth, t, wg)
		case request.unreal:
			request.unreal = true
			wg.Add(1)
			go sendRequest(request.requestPoolunreal, t, wg)
		default:
			request.multi = true
			wg.Add(1)
			go sendRequest(request.requestPoolmulti, t, wg)
		}
	}
	wg.Wait()
}
