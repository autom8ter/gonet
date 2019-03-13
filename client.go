package gonet

import (
	"encoding/base64"
	"github.com/autom8ter/util"
	"github.com/autom8ter/util/netutil"
	"github.com/gorilla/sessions"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	req *http.Request
	cli *http.Client
}

func NewClient(u *url.URL, method string) *Client {
	var r = &http.Request{
		URL:    u,
		Method: method,
	}
	return &Client{
		req: r,
		cli: http.DefaultClient,
	}
}
func NewCustomClient(u *url.URL, method string, client *http.Client) *Client {
	var r = &http.Request{
		URL:    u,
		Method: strings.ToUpper(method),
	}
	return &Client{
		req: r,
		cli: client,
	}
}
func (c *Client) SetHeaders(headers map[string]string) {
	c.req = netutil.SetHeaders(headers, c.req)
}

func (c *Client) SetForm(vals map[string]string) {
	c.req = netutil.SetForm(vals, c.req)
}

func (c *Client) SetBasicAuth(userName, password string) {
	c.req = netutil.SetBasicAuth(userName, password, c.req)
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

func (c *Client) Init(headers map[string]string, formvals map[string]string, user, password string) {
	if headers != nil {
		c.SetHeaders(headers)
	}
	if formvals != nil {
		c.SetForm(formvals)
	}
	if user != "" && password != "" {
		c.SetBasicAuth(user, password)
	}
}

func (r *Client) Render(s string, data interface{}) string {
	return util.Render(s, data)
}

func (r *Client) NewSessionFSStore() *sessions.FilesystemStore {
	return netutil.NewSessionFileStore()
}

func (r *Client) NewSessionCookieStore() *sessions.CookieStore {
	return netutil.NewSessionCookieStore()
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
	return netutil.ReadBody(resp)
}
