package scraper

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/arl/statsviz"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/jmoussa/crypto-dashboard/config"
	"github.com/jmoussa/crypto-dashboard/twittermicro/processors"
	"golang.org/x/sync/semaphore"
)

// Parse JSON config for use
var Cfg config.Config = config.ParseConfig()

func TwitterTweetGenerator(searchPhrase string, cfg config.Config) chan interface{} {
	// Starts up a generator stream of tweets into the outputted channel
	out := make(chan interface{})
	go func() {
		defer close(out)

		con := cfg.Twitter
		c := oauth1.NewConfig(con["consumerkey"], con["consumersecret"])
		token := oauth1.NewToken(con["accesstoken"], con["accesssecret"])
		httpClient := c.Client(oauth1.NoContext, token)

		// intialize stream
		client := twitter.NewClient(httpClient)
		params := &twitter.StreamFilterParams{
			Track:         []string{searchPhrase},
			StallWarnings: twitter.Bool(true),
		}
		stream, err := client.Streams.Filter(params)
		if err != nil {
			log.Fatalf("Error querying stream, %s\n", err)
		}
		defer stream.Stop()

		// Initialize demux for interface{} type processing to channel
		demux := twitter.NewSwitchDemux()
		demux.Tweet = func(tweet *twitter.Tweet) {
			out <- tweet
		}
		for message := range stream.Messages {
			demux.Handle(message)
		}
	}()
	return out
}

func mergeAtomic(outputChan chan interface{}, cs ...<-chan interface{}) <-chan interface{} {
	// Atomically dump each channel into the output channel and return output channel
	var i int32
	atomic.StoreInt32(&i, int32(len(cs)))
	for _, c := range cs {
		go func(c <-chan interface{}) {
			for v := range c {
				outputChan <- v
			}
			if atomic.AddInt32(&i, -1) == 0 {
				close(outputChan)
			}
		}(c)
	}
	return outputChan
}

func broadcast(waitTime time.Duration, inputChannel chan interface{}, chs ...chan interface{}) (sentNumber int) {
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()

	jobQueue := make(chan chan interface{}, len(chs))
	for _, c := range chs {
		jobQueue <- c
	}

queue:
	for c := range jobQueue {
		select {
		case c <- inputChannel:
			// sent success
			sentNumber += 1
			if sentNumber == len(chs) {
				cancel()
			}
		case <-ctx.Done():
			// timeout, break job queue
			break queue
		default:
			// if send failed, retry later
			jobQueue <- c
		}
	}

	return
}

func sink(ctx context.Context, cancelFunc context.CancelFunc, values <-chan interface{}, errors <-chan error) {
	var count int64 = 0
	for {
		select {
		case <-ctx.Done():
			log.Print(ctx.Err().Error())
			return
		case err := <-errors:
			if err != nil {
				log.Println("error: ", err.Error())
				cancelFunc()
			}
		case _, ok := <-values:
			if ok {
				count += 1
				if count%100 == 0 {
					log.Printf("Tweet count: %d", count)
				}
			} else {
				log.Print("done")
				return
			}
		}
	}
}

func step[In any, Out any](
	ctx context.Context,
	inputChannel <-chan In,
	outputChannel chan Out,
	errorChannel chan error,
	fn func(In) (Out, error),
) {
	defer close(outputChannel)

	// create a new semaphore with a limit (of the CPU count) for processes the semaphore can access at a time
	limit := runtime.NumCPU()
	sem1 := semaphore.NewWeighted(int64(limit))

	// parse through messages in input channel
	for s := range inputChannel {
		select {
		// if cancelled, abort operation otherwise run while there's values in inputChannel
		case <-ctx.Done():
			log.Println("1 abort")
			break
		default:
		}
		// use semaphores to keep data integrity
		if err := sem1.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		// start up go functions to parallelize processing to CPU Cound
		go func(s In) {
			// release the semaphore at the end of this concurrent process
			defer sem1.Release(1)
			// Take the result of the function and send to outputChannel
			result, err := fn(s)
			if err != nil {
				errorChannel <- err
			} else {
				outputChannel <- result
			}
		}(s)
	}

	// after everything's finished fetch and lock the semaphore
	if err := sem1.Acquire(ctx, int64(limit)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}
}

func main() {
	// Simulate pipeline running E2E

	statsviz.RegisterDefault()
	go func() {
		log.Println("Navigate to: http://localhost:6070/debug/statsviz/ for metrics")
		log.Println(http.ListenAndServe("localhost:6070", nil))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*
		readStream, err := producer(ctx, source)
		if err != nil {
			log.Fatal(err)
		}
	*/

	// using generator as initial producer (outputs an interface{} channel)
	// BUFSIZE := 50
	// b := broadcaster.New(BUFSIZE)
	// sourceChannel := generator("covid", cfg)

	inputChannel1 := make(chan interface{})
	inputChannel2 := make(chan interface{})
	errorChannel := make(chan error)

	/*
		// broadcasting source to both lexicon and ML layers
		go func() {
			broadcast(10*time.Second, sourceChannel, inputChannel1, inputChannel2)
		}()
	*/
	// Layer 1: Sentiment Analysis (Lexicon)
	layer1OutputChannel := make(chan interface{})
	go func() {
		step(ctx, inputChannel1, layer1OutputChannel, errorChannel, processors.LexiconSentimentAnalysis)
	}()

	// Layer 2: Sentiment Analysis (ML)
	layer2OutputChannel := make(chan interface{})
	go func() {
		step(ctx, inputChannel2, layer2OutputChannel, errorChannel, processors.IMDBModelSentimentAnalysis)
	}()
	// Merge Output Channels
	inputChannel := make(chan interface{})
	mergeAtomic(inputChannel, layer1OutputChannel, layer2OutputChannel)

	// DB Upload tweets
	layer3OutputChannel := make(chan interface{})
	go func() {
		step(ctx, inputChannel, layer3OutputChannel, errorChannel, processors.FormatAndUpload)
	}()
	// Sink to final stream (make available via API)
	sink(ctx, cancel, layer3OutputChannel, errorChannel)
}
