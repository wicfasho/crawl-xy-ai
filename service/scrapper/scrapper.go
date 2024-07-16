package scrapper

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/wicfasho/crawl-xy/database"
	"github.com/wicfasho/crawl-xy/server"
	"github.com/wicfasho/crawl-xy/sqlc"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type Page struct {
	url                          string
	title, description, keywords pgtype.Text
	lastVisited                  time.Time
}

type ScrapperModel interface {
	llms.Model
	embeddings.EmbedderClient
}

type AITool struct {
	llm           ScrapperModel
	embedder      *embeddings.EmbedderImpl
	pgVectorStore *pgvector.Store
}

const OllamaServerURL = "http://127.0.0.1:11434/"

var converter = md.NewConverter("", true, nil)
var parentURLs = []string{
	"conestogac.on.ca",
}

var Tool = AITool{}
var nodesToBeRemoved = []string{
	"script",
	"noscript",
}
var collyMaxDepth = 1
var similaritySearchNumDocuments = 3 // Number of docs to be return from similarity search

func Start() {
	var wg sync.WaitGroup

	// Init pgvector store
	tool, err := createPgVectorStore("openai")
	if err != nil {
		log.Fatal().Msgf("couldn't create store: %v\n", err)
	}
	Tool = *tool

	// Scrap Pages
	// Scrap(&wg) //Enable this to scrap webpages

	wg.Wait()
}

func Scrap(wg *sync.WaitGroup) {
	for _, pURL := range parentURLs {
		wg.Add(1)
		go func(inputURL string) {
			defer wg.Done()

			parsedURL, err := url.Parse("https://" + inputURL)
			if err != nil {
				log.Fatal().Msgf("invalid url: %s", parsedURL.RequestURI())
			}
			fetchContent(parsedURL.String())
		}(pURL)
	}
}

func fetchContent(url string) {
	var fcWg sync.WaitGroup
	ctx := context.Background()

	c := colly.NewCollector(
		colly.MaxDepth(collyMaxDepth),
		// colly.AllowedDomains(parentURLs...),
		colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnError(func(r *colly.Response, err error) {
		log.Err(err).Msg("An Error Occurred")
	})

	c.OnScraped(func(r *colly.Response) {
		saveVisitedURL(r.Request.URL.String())
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if allowURL(link) && !urlIsVisited(link) {
			e.Request.Visit(link) // Continue visiting valid URLs found on the page
		}
	})

	c.OnHTML("html", func(h *colly.HTMLElement) {
		page := Page{
			url: h.Request.URL.String(),
		}

		if !urlIsVisited(page.url) {
			titleTag := h.ChildTexts("title")
			if len(titleTag) > 0 {
				page.title = pgtype.Text{
					String: titleTag[0],
					Valid:  true,
				}
			}

			metaTagDescription := h.ChildAttrs("meta[name='description']", "content")
			if len(metaTagDescription) > 0 {
				page.description = pgtype.Text{
					String: metaTagDescription[0],
					Valid:  true,
				}
			}

			metaTagKeywords := h.ChildAttrs("meta[name='keywords']", "keywords")
			if len(metaTagKeywords) > 0 {
				page.keywords = pgtype.Text{
					String: metaTagKeywords[0],
					Valid:  true,
				}
			}

			page.lastVisited = time.Now()

			// Save Page data to DB
			server.GetServer().GetDBStore().InsertPage(ctx, sqlc.InsertPageParams{
				Url:             page.url,
				Title:           page.title,
				MetaDescription: page.description,
				MetaKeywords:    page.keywords,
			})

			// Squeeze out body from html and remove unncessary nodes
			body := h.DOM.Find("body")
			body.Children().Each(func(i int, s *goquery.Selection) {
				for _, node := range nodesToBeRemoved {
					if node == s.Nodes[0].Data {
						s.Remove()
					}
				}
			})
			pageContent, err := body.Html()
			if err != nil {
				log.Err(err).Msg("error creating html from page content")
				return
			}

			// Convert html to markdown and trim all unneccessary spaces
			markdown, err := converter.ConvertString(pageContent)
			if err != nil {
				log.Err(err).Msg("could not convert html to markdown")
				return
			}
			markdown = strings.ReplaceAll(markdown, "\n\n", "\n")
			// fmt.Println("md ->", markdown)

			// Save Embeddings to pgvector store
			if len(markdown) > 0 {
				_, err = Tool.pgVectorStore.AddDocuments(ctx, []schema.Document{
					{
						PageContent: fmt.Sprintf(`Webpage in Markdown: '%v'`, markdown),
						Metadata: map[string]any{
							"title": page.title,
							"url":   page.url,
						},
					},
				})
				if err != nil {
					log.Err(err).Msg("could not save embeddings to store")
					return
				}
			}
		}

		// log.Info().Msgf("%v", page)
	})

	c.Visit(url)

	fcWg.Wait()
}

func Ask(question string, sessionID string) (*string, error) {
	ctx := context.Background()

	fSimSearchQuestion := fmt.Sprintf(`Search for documents that can give granular information on this question. Question: '%s'`, question)


	fPromptQuestion := fmt.Sprintf(`%v. Question: '%s'`, os.Getenv("FPROMPT"), question) // Change prompt to suit requirement

	docs, err := Tool.pgVectorStore.SimilaritySearch(ctx, fSimSearchQuestion, similaritySearchNumDocuments,
		vectorstores.WithEmbedder(Tool.embedder),
	)
	if err != nil {
		return nil, fmt.Errorf("error running similarity serach: %v", err)
	}

	// fmt.Println("->", docs) // Uncomment to see filtered docs

	// Create a chain to ask question
	stuffQAChain := chains.LoadStuffQA(Tool.llm)
	answer, err := chains.Call(
		ctx,
		stuffQAChain,
		map[string]any{
			"input_documents": docs,
			"question":        fPromptQuestion,
		},
		chains.WithModel("gpt-3.5-turbo-0125"),
		chains.WithTemperature(0.8),
	)
	if err != nil {
		return nil, fmt.Errorf("error making chains call: %v", err)
	}

	answerText, ok := answer["text"]
	if !ok {
		return nil, fmt.Errorf("unable to get result for question")
	}

	result, ok := answerText.(string)
	if !ok {
		return nil, fmt.Errorf("unable to get result for question: error converting result to string")
	}

	return &result, nil
}

func createPgVectorStore(selectLLM string) (*AITool, error) {
	ctx := context.Background()

	Tool, err := initLLM(selectLLM)
	if err != nil {
		return nil, err
	}

	store, err := pgvector.New(
		ctx,
		// pgvector.WithConnectionURL(os.Getenv("DATABASE_URL")),
		pgvector.WithConn(database.GetDB()),
		pgvector.WithEmbedder(Tool.embedder),
	)
	if err != nil {
		return nil, err
	}

	Tool.pgVectorStore = &store

	return Tool, nil
}

func initLLM(selectLLM string) (*AITool, error) {
	ctx := context.Background()

	switch selectLLM {
	case "openai":
		llm, err := openai.New(
			openai.WithToken(os.Getenv("OPEN_AI_KEY")),
			openai.WithModel("gpt-3.5-turbo-0125"),
			openai.WithEmbeddingModel("text-embedding-3-small"),
		)
		if err != nil {
			return nil, err
		}

		embedding, err := embeddings.NewEmbedder(llm)
		if err != nil {
			return nil, err
		}

		return &AITool{
			llm,
			embedding,
			&pgvector.Store{},
		}, nil

	case "googleai":
		llm, err := googleai.New(
			ctx,
			googleai.WithAPIKey(os.Getenv("GEN_AI_KEY")),
			googleai.WithDefaultModel("gemini-pro"),
			googleai.WithDefaultEmbeddingModel("models/text-embedding-004"),
		)
		if err != nil {
			return nil, err
		}

		embedding, err := embeddings.NewEmbedder(llm)
		if err != nil {
			return nil, err
		}

		return &AITool{
			llm,
			embedding,
			&pgvector.Store{},
		}, nil

	case "ollama":
		llm, err := ollama.New(
			ollama.WithModel("[enter model name here]"),
			ollama.WithServerURL(OllamaServerURL),
		)
		if err != nil {
			return nil, err
		}

		embedding, err := embeddings.NewEmbedder(llm)
		if err != nil {
			return nil, err
		}

		return &AITool{
			llm,
			embedding,
			&pgvector.Store{},
		}, nil
	}

	return &AITool{}, nil
}

func allowURL(inputURL string) bool {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	host := parsedURL.Hostname()

	for _, parentURL := range parentURLs {
		if strings.Contains(host, parentURL) {
			return true
		}
	}

	return false
}

func saveVisitedURL(inputURL string) error {
	file, err := os.OpenFile("visited_urls.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer file.Close()

	if urlIsVisited(inputURL) {
		log.Info().Msgf("url already saved: %s", inputURL)
		return nil
	}

	_, err = fmt.Fprintln(file, inputURL)
	if err != nil {
		return err
	}

	log.Info().Msgf("url written to file: %v", inputURL)
	return nil
}

func urlIsVisited(inputURL string) bool {
	file, err := os.OpenFile("visited_urls.txt", os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == inputURL {
			return true
		}
	}

	return false
}

func reduceToken(docs []schema.Document) []schema.Document {
	// tokenLimitPerMinute := 60000
	// tokenLimitPerRequest := 16385
	return []schema.Document{}
}

func Stop() {

}
