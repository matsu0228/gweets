package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

func errorExit(err error) {
	log.Fatalf("[ERROR]", err)
}

func envLoad() error {
	if err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))); err != nil {
		return err
	}
	return nil
}

func generateURL(tw anaconda.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", tw.User.IdStr, tw.IdStr)
}

// ExtractText :
func ExtractText(org string) (string, error) {
	re, err := regexp.Compile(`#([A-Z]|[a-z]|[0-9])?`)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(org, ""), nil
}

func main() {
	if err := envLoad(); err != nil {
		errorExit(err)
	}
	at := os.Getenv("TW_ACCESS_TOKEN")
	as := os.Getenv("TW_ACCESS_TOKEN_SECRET")
	ck := os.Getenv("WT_CONSUMER_KEY")
	cs := os.Getenv("WT_CONSUMER_SECRET")

	log.Printf("connect api with %v  %v  %v  %v", at, as, ck, cs)
	api := anaconda.NewTwitterApiWithCredentials(at, as, ck, cs)

	v := url.Values{}
	v.Set("count", "100")
	v.Set("lang", "ja")
	result, err := api.GetSearch("#NZLvNAM -RT", v)

	if err != nil {
		errorExit(err)
	}
	for _, tw := range result.Statuses {
		// if tw.Lang != "ja" {
		// 	continue
		// }
		// if strings.Contains(tw.Text, "RT") {
		// 	continue
		// }

		body, err := ExtractText(tw.FullText)
		if err != nil {
			log.Printf("[WARN] cant extract text. %v", err)
		}
		fmt.Printf("--------------------- \ntweet: %v [%v] %v => %v (%v) \n\n", tw.IdStr, tw.User.Name, tw.FullText, body, generateURL(tw))
	}
}
