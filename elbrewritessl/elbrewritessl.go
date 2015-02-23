package elbrewritessl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mailgun/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/mailgun/vulcand/Godeps/_workspace/src/github.com/mailgun/oxy/utils"
	"github.com/mailgun/vulcand/plugin"
)

const Type = "elb-rewrite-ssl"

type Rewrite struct {
	Redirect bool
}

func NewRewrite(redirect bool) (*Rewrite, error) {
	return &Rewrite{redirect}, nil
}

func (rw *Rewrite) NewHandler(next http.Handler) (http.Handler, error) {
	return newRewriteHandler(next, rw)
}

func (rw *Rewrite) String() string {
	return fmt.Sprintf("redirect=%v", rw.Redirect)
}

type rewriteHandler struct {
	next       http.Handler
	errHandler utils.ErrorHandler
	redirect   bool
}

func newRewriteHandler(next http.Handler, spec *Rewrite) (*rewriteHandler, error) {
	return &rewriteHandler{
		redirect:   spec.Redirect,
		next:       next,
		errHandler: utils.DefaultHandler,
	}, nil
}

func (rw *rewriteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	protocol := r.Header.Get("X-Forwarded-Proto")

	if rw.redirect && protocol != "https" {
		httpsUrl := strings.Join([]string{"https://", r.Host, r.RequestURI}, "")
		parsedURL, err := url.Parse(httpsUrl)
		if err != nil {
			rw.errHandler.ServeHTTP(w, r, err)
			return
		}

		(&redirectHandler{u: parsedURL}).ServeHTTP(w, r)
		return
	}
	rw.next.ServeHTTP(w, r)
}

func FromOther(rw Rewrite) (plugin.Middleware, error) {
	return NewRewrite(rw.Redirect)
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	return NewRewrite(c.Bool("redirect"))
}

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

func CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "elb-rewrite-ssl",
			Usage: "if provided, request is redirected to the HTTPS URL",
		},
	}
}

type redirectHandler struct {
	u *url.URL
}

func (f *redirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Location", f.u.String())
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(http.StatusText(http.StatusFound)))
}
