package processors

import (
	"context"
	"log"
	"time"

	"github.com/cdipaolo/sentiment"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TweetWithScore struct {
	BaseTweet *twitter.Tweet
	Score     interface{} //map[string]float64
	Type      string
	//IMDBMLSentimentScores uint8
}

func LexiconSentimentAnalysis(s interface{}) (interface{}, error) {
	// Takes in an interface{} message, fetchest sentiment scores, and pushes updated message with scores
	tweet := s.(twitter.Tweet)
	parseText := sentitext.Parse(tweet.Text, lexicon.DefaultLexicon)
	results := sentitext.PolarityScore(parseText)
	//log.Println("Positive:", results.Positive)
	//log.Println("Negative:", results.Negative)
	//log.Println("Neutral:", results.Neutral)
	//log.Println("Compound:", results.Compound)
	scores := map[string]float64{
		"Positive": results.Positive,
		"Negative": results.Negative,
		"Neutral":  results.Neutral,
		"Compound": results.Compound,
	}
	var obj TweetWithScore
	obj.BaseTweet = &tweet
	obj.Score = scores
	obj.Type = "lexicon"
	return obj, nil
}

func FormatAndUpload(s interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			return
		}
	}()
	collection := client.Database("twitter-sentiment").Collection("tweets")

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"basetweet.id": s.(TweetWithScore).BaseTweet.ID}
	score := s.(TweetWithScore).Score
	t := s.(TweetWithScore).Type
	update := bson.M{"$set": bson.M{t: score, "basetweet": s.(TweetWithScore).BaseTweet}}
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf("Error Upserting")
	}
	log.Println(result)
	//collection.InsertOne(ctx, s)
	// Takes in a string message, alters it and pushes updated message
	//tweet := s.(TweetWithScore)
	//log.Printf("Process 2: Text - %s\n--------------------------------\n", tweet.BaseTweet.Text)
	return s, nil
}

func IMDBModelSentimentAnalysis(s interface{}) (interface{}, error) {
	log.Println(s)
	tweet := s.(twitter.Tweet)
	// Model : restore or train(project directory)
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		log.Printf("Error formatting model: %s", err)
		return 0, err
		//panic(err)
	}
	results := sentimentModel.SentimentAnalysis(tweet.Text, sentiment.English)
	score := results.Score

	var obj TweetWithScore
	obj.BaseTweet = &tweet
	obj.Score = score
	obj.Type = "imdb_ml_model"
	return obj, nil
}
