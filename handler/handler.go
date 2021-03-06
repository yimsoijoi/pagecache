package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	myredis "github.com/yimsoijoi/pagecache/redis"
)

type Handler interface {
}
type handler struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *handler {
	return &handler{
		rdb: rdb,
	}
}

type Req struct {
	Websites []string `json:"websites"`
}

type Res map[string]string

func (h *handler) Handle(c *fiber.Ctx) error {
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(map[string]interface{}{
			"error message": "bad request",
			"error":         err.Error(),
		})
	}
	res := make(Res)
	var mut sync.RWMutex
	var wg sync.WaitGroup
	for _, site := range req.Websites {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			body, err := worker(u)
			if err != nil {
				c.Status(400).JSON(map[string]interface{}{
					"error": err.Error(),
				})
			}
			mut.Lock()
			res[u] = string(body)
			mut.Unlock()
		}(site)
	}
	return c.Status(200).JSON(res)
}

func worker(url string) ([]byte, error) {
	body, err := myredis.ReadFromRedis(url)
	if errors.Is(err, redis.Nil) {
		body, err = GetBody(url)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get HTML body")
		}
	}
	if err := myredis.WriteToRedis(url, string(body)); err != nil {
		return nil, errors.Wrap(err, "failed to write body to Redis")
	}
	return body, nil
}

func GetBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("http.Get failed for URL %s", url))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read response body from URL %s", url))
	}
	return body, nil
}
