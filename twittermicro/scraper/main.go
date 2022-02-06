package scraper

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
	pb "github.com/jmoussa/crypto-dashboard/twittermicro/twitter_pb"
)

func TweetToProto(tweet *twitter.Tweet) (*pb.Tweet, error) {
	// Convert twitter.Tweet to protobuf.Tweet
	pbTweet := &pb.Tweet{
		TweetId:   tweet.IDStr,
		Title:     tweet.Text,
		Text:      tweet.FullText,
		Url:       tweet.Source,
		Date:      tweet.CreatedAt,
		Author:    tweet.User.ScreenName,
		AuthorUrl: tweet.User.URL,
	}
	return pbTweet, nil
}

func ScrapeTwitterData(req *pb.GetTwitterDataRequest) ([]*pb.Tweet, error) {
	// Scrape data from Twitter
	// max_entries := req.MaxEntries
	items := []*pb.Tweet{}
	log.Println("Scraping Twitter Data")
	gen := TwitterTweetGenerator(req.SearchPhrase, Cfg)
	var i int64 = 0
	if req.MaxEntries <= 0 {
		req.MaxEntries = 100
	}
	for i = 0; i < req.MaxEntries; i++ {
		// Get the next tweet from the generator
		tweet := <-gen
		// Convert to protobuf
		pbTweet, err := TweetToProto(tweet.(*twitter.Tweet))
		if err != nil {
			log.Printf("Error converting tweet to proto: %s\n", err)
			continue
		}
		// Append to slice
		items = append(items, pbTweet)
	}
	log.Printf("Scraped %d items", len(items))
	return items, nil
}
