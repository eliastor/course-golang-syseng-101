package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"

	"github.com/brianvoe/gofakeit"
)

// Document stores text created by bureaucrat and have field for signature
type Document struct {
	text string
	sign []byte
}

// Let's define three roles in our office: bureaucrat, executor and verifier

// bureaucrat generates documents and take a rest after every document creation.
// as soon as bureaucrat gets signal (done channel is closed) to finish its work it stops
func bureaucrat(done <-chan struct{}, out chan<- *Document) int {
	docs_total := 0
	for {
		select {
		case _ = <-done:
			return docs_total
		case out <- &Document{text: gofakeit.Sentence(9)}:
			docs_total++
			//

			time.Sleep(time.Millisecond * time.Duration(mathrand.Intn(100)))
		}
	}
}

// executor signs documents it receives and send the documents further
func executor(priv ed25519.PrivateKey, in <-chan *Document, out chan<- *Document) int {
	docs_total := 0
	for doc := range in {
		doc.sign = ed25519.Sign(priv, []byte(doc.text))
		docs_total++
		out <- doc
	}
	return docs_total
}

// verifier checks signature in every document it receives and send the documents further
func verifier(pub ed25519.PublicKey, in <-chan *Document, out chan<- *Document) int {
	docs_total := 0
	for doc := range in {
		if ed25519.Verify(pub, []byte(doc.text), doc.sign) {
			out <- doc
		}
		docs_total++
	}
	return docs_total
}

func fanoutProxy(in <-chan *Document, fakeIn chan *Document, out chan *Document) {
	for msg := range in {
		fakeIn <- msg
	}
	close(fakeIn)
	close(out)
}

// SpawnBureaucrat is generator function that spawns one bureaucrat,
// and returns channel with generated documents. Also handles channel lifecycle
func SpawnBureaucrat(done <-chan struct{}) <-chan *Document {
	out := make(chan *Document)
	go func() {
		total := bureaucrat(done, out)
		fmt.Println(total, "documents were created by Bureaucrat")
		close(out)
	}()
	return out
}

// SpawnExecutors is generator function that create n executors,
// fan out documents to executors from input channel,
// and returns channel with signed documents. Also handles channel lifecycle
func SpawnExecutors(n int, priv ed25519.PrivateKey, in <-chan *Document) <-chan *Document {
	out := make(chan *Document)
	fakeIn := make(chan *Document)
	totals := make([]int, n)
	// we must propagate closing of channels from INs to OUTs
	// close() must be called only once
	// we spawn multiple workers
	// so we need one goroutine from that we can catch channel closing and close out channel
	// for this purpose fanoutProxy
	// fanoutProxy reads documents from IN channel, sends them to fakeIn
	// If IN channel is closed it will close fakeIN and OUT,
	// so executors will be closed gracefully, because of close(fakeIN)
	// and close will be propagated further to Out channel
	go fanoutProxy(in, fakeIn, out)
	for i := 0; i < n; i++ {
		go func(i int) {
			totals[i] = executor(priv, fakeIn, out)
			fmt.Println(totals[i], "documents were signed by Executor", i) // why it is not shown in output?
		}(i)
	}

	return out
}

// SpawnExecutors is generator function that create one verifier,
// and returns channel with verified documents. Also handles channel lifecycle
func SpawnVerifier(pub ed25519.PublicKey, in <-chan *Document) <-chan *Document {
	out := make(chan *Document)
	go func() {
		total := verifier(pub, in, out)
		fmt.Println(total, "documents were verified by verifier")
		close(out)
	}()
	return out
}

func main() {
	// to get executors and verifier work we need to provide private and public keys
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	done := make(chan struct{})
	docsNew := SpawnBureaucrat(done)
	docsSigned := SpawnExecutors(2, priv, docsNew)
	docsVerified := SpawnVerifier(pub, docsSigned)

	// we don't want this to run infinitely, so we'll stop this by signaling bureaucrat to stop
	go func(done chan struct{}) {
		time.Sleep(3 * time.Second)
		close(done)
		// channel closing flow is:
		// close(done) --> close(docsNew)--> [fanoutProxy] --> close(docsSigned) --> close(docsVerified)
	}(done)

	// the whole structure will be like that:
	//                       / [executor 1] \
	// [bureaucrat]---------|                |--------[verifier]-------[range in main]
	//                       \ [executor 2] /

	sentToNowhere := 0
	// nowhere emulation
	for range docsVerified {
		sentToNowhere++
	}

	fmt.Println(sentToNowhere, "documents were sent to nowhere")
}
