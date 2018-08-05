package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	TID = "0xb1ed364e4333aae1da4a901d5231244ba6a35f9421d4607f7cb90d60bf45578a"
	URL = "https://mainnet.infura.io"
)

var (
	sess            session.Session
	svc             *translate.Translate
	articleTextInPt string
)

func confAWS() {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", "bleu-hackathon"),
	})
	if err != nil {
		log.Fatalf("could not create AWS sessioin: %s\n", err)
	}
	svc = translate.New(sess)
}

func downloadChineseArticle() (string, error) {
	var textInChinese string

	log.Println("retrieving chinese report from ethereum blockchain")

	rpcCli, errRPCClient := rpc.Dial(URL)
	if errRPCClient != nil {
		return textInChinese, fmt.Errorf("RPCC dial error: %v", errRPCClient)
	}

	var cli = ethclient.NewClient(rpcCli)
	var ctx = context.Background()

	tx, isPending, err := cli.TransactionByHash(ctx, common.HexToHash(TID))

	if err != nil {
		log.Fatalf("TransactionByHash error: %v\n", err)
	} else if isPending == false {
		textInChinese = string(tx.Data())
	}

	log.Println("successfully downloaded text from chinese report")
	return textInChinese, nil
}

func chinese2English(chineseArticleLines []string) ([]string, error) {
	var reportLinesEn []string
	log.Println("using AWS to translate report from chinese to english")
	for _, line := range chineseArticleLines {

		if len(line) == 0 {
			continue
		}

		txtInput := translate.TextInput{
			SourceLanguageCode: aws.String("zh"),
			TargetLanguageCode: aws.String("en"),
			Text:               aws.String(line),
		}

		chinese2English, err := svc.Text(&txtInput)
		if err != nil {
			return reportLinesEn, fmt.Errorf("error translating report: %v", err)
		}

		reportLinesEn = append(reportLinesEn, *chinese2English.TranslatedText)

	}
	log.Println("successfully translated chinese report to english")
	return reportLinesEn, nil
}

func english2Portuguese(englishArticleLines []string) ([]string, error) {
	var reportLinesPt []string
	log.Println("using AWS to translate english translation to portguese")
	for _, line := range englishArticleLines {

		if len(line) == 0 {
			continue
		}

		txtInput := translate.TextInput{
			SourceLanguageCode: aws.String("en"),
			TargetLanguageCode: aws.String("pt"),
			Text:               aws.String(line),
		}

		portuguese2English, err := svc.Text(&txtInput)
		if err != nil {
			return reportLinesPt, fmt.Errorf("error translating report: %v", err)
		}

		reportLinesPt = append(reportLinesPt, *portuguese2English.TranslatedText)

	}
	log.Println("successfully translated english lines to portuguese")
	return reportLinesPt, nil
}

type donaMariaPage struct {
	Title      string
	SubTitle   string
	ReportText string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := donaMariaPage{
		Title:      "Olá, Dona Maria!",
		SubTitle:   "Aqui está a notícia censurada pelo governo chinês :)",
		ReportText: articleTextInPt,
	}
	t, err := template.ParseFiles("templates/donamaria.gohtml")
	if err != nil {
		log.Fatalf("error parsing template: %v\n", err)
	}
	t.Execute(w, p)
}

func init() {
	confAWS()
}

func main() {

	zh, err := downloadChineseArticle()
	if err != nil {
		log.Fatalf("error downloading chinese article: %v\n", err)
	}

	reportLinesZh := strings.Split(zh, "\n")

	reportLinesEn, err := chinese2English(reportLinesZh)
	if err != nil {
		log.Fatalf("could not translate chinese article to english: %v\n", err)
	}

	reportLinesPt, err := english2Portuguese(reportLinesEn)
	if err != nil {
		log.Fatalf("could not translate english lines to chinese: %v\n", err)
	}

	articleTextInPt = strings.Join(reportLinesPt, "\n")

	log.Println("serving application at port 8080")
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)

}
