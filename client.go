package gonet

import (
	"context"
	"encoding/base64"
	"github.com/autom8ter/util"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type Client struct {
	req *http.Request
	cli *http.Client
}

func NewClient(ctx context.Context, urL, method, user, password, body string, headers map[string]string, form map[string]string) *Client {
	var r, err = util.NewRequestCtx(ctx, method, urL, user, password, headers, form, strings.NewReader(body))
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Client{
		req: r,
		cli: http.DefaultClient,
	}
}

func (c *Client) WithClient(client *http.Client) {
	c.cli = client
}

func (c *Client) Client() *http.Client {
	return c.cli
}

func (c *Client) Request() *http.Request {
	return c.req
}

func (c *Client) Stringify(obj interface{}) string {
	return util.ToPrettyJsonString(obj)
}

func (c *Client) JSONify(obj interface{}) []byte {
	return util.ToPrettyJson(obj)
}

func (c *Client) AsCsv(s string) ([]string, error) {
	return util.ReadAsCSV(s)
}

func (c *Client) ToWriter(w io.Writer) error {
	return c.req.Write(w)
}

func (c *Client) Prompt(q string) string {
	return util.Prompt(q)
}

func (c *Client) SetRequest(req *http.Request) {
	c.req = req
}

func (c *Client) Do() (*http.Response, error) {
	return c.cli.Do(c.req)
}

func (c *Client) GenerateJWT(signingKey string, claims map[string]interface{}) (string, error) {
	return util.GenerateJWT(signingKey, claims)
}

func (r *Client) Render(s string, data interface{}) string {
	return util.Render(s, data)
}

func (r *Client) RandomTokenString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (r *Client) RandomToken(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func (r *Client) DerivePassword(counter uint32, password_type, password, user, site string) string {
	return util.DerivePassword(counter, password, password, user, site)
}

func (r *Client) GeneratePrivateKey(typ string) string {
	return util.GeneratePrivateKey(typ)
}

func (r *Client) ReadBody(resp *http.Response) ([]byte, error) {
	return util.ReadBody(resp)
}
