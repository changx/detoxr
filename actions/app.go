package actions

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/changx/detoxr/dist"
	"github.com/changx/detoxr/ds"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n/v2"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/x/sessions"
	"github.com/gorilla/handlers"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app *buffalo.App
	T   *i18n.Translator
)

var oooops = []byte("Something wrong. Please try again later.")

func ooops(status int, err error, c buffalo.Context) error {
	c.Logger().Error(errors.WithStack(err))
	res := c.Response()
	res.WriteHeader(status)
	res.Write(oooops)
	return nil
}

func JSONLogger(lvl logger.Level) logger.FieldLogger {
	l := logrus.New()
	l.Level = lvl
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	return logger.Logrus{FieldLogger: l}
}

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		ds.InitServerConfig()
		go ds.Serve()

		cors := cors.New(cors.Options{
			Debug: false,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodPut,
				http.MethodHead,
				http.MethodOptions,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler

		// redisHost := envy.Get("REDIS_ADDR", ":6379")

		// worker := gwa.New(gwa.Options{
		// 	Pool: &redis.Pool{
		// 		MaxActive: 5,
		// 		MaxIdle:   2,
		// 		Wait:      true,
		// 		Dial: func() (redis.Conn, error) {
		// 			return redis.Dial("tcp", redisHost)
		// 		},
		// 	},
		// 	Name:           "112mi",
		// 	MaxConcurrency: 500,
		// })

		var lg logger.FieldLogger

		if ENV == "production" {
			lg = JSONLogger(logger.InfoLevel)
		} else {
			lg = JSONLogger(logger.DebugLevel)
		}

		app = buffalo.New(buffalo.Options{
			Env:           ENV,
			SessionStore:  sessions.Null{},
			CompressFiles: true,
			PreWares:      []buffalo.PreWare{cors},
			Logger:        lg,
			Addr:          ds.GetWebAddr(),
			// Worker:        worker,
		})

		app.ErrorHandlers[405] = ooops
		app.ErrorHandlers[500] = ooops

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		app.GET("/api/settings/honeypot", GetHoneypot)
		app.POST("/api/settings/honeypot/try", TryHoneypot)
		app.POST("/api/settings/honeypot", SaveHoneypot)

		app.GET("/api/settings/local_ns", GetLocalNS)
		app.POST("/api/settings/local_ns/try", TryLocalNS)
		app.POST("/api/settings/local_ns", SaveLocalNS)

		app.GET("/api/settings/doh_service", GetDoHServiceURL)
		app.POST("/api/settings/doh_service/try", TryDoHServiceWithURL)
		app.POST("/api/settings/doh_service", SaveDoHServiceURL)

		app.GET("/api/data/safelist", GetSafelist)
		app.GET("/api/data/victims", GetVictims)
		// capture all to /index.html
		app.Muxer().PathPrefix("/").Handler(pwaHandler(http.FS(dist.FS())))
	}

	return app
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func pwaHandler(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			f, _ = fs.Open("/")
			r.URL.Path = "/"
		}

		stat, _ := f.Stat()
		maxAge := envy.Get(buffalo.AssetsAgeVarName, "31536000")
		w.Header().Add("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%s", maxAge))
		fsh.ServeHTTP(w, r)
	})

	return handlers.CompressHandler(baseHandler)
}

func invalidRequest(c buffalo.Context) error {
	err := errors.New("invalid request")
	c.Logger().Error(errors.WithStack(err).Error())
	return errorResponse(c, InvalidRequest, err.Error())
}

func serverErrorWithShortMessage(c buffalo.Context, clientErr string, logErr any) error {
	switch v := logErr.(type) {
	case error:
		c.Logger().Errorf("%s: %s", clientErr, errors.WithStack(v).Error())
	case string:
		c.Logger().Errorf("%s: %s", clientErr, errors.WithStack(errors.New(v)).Error())
	default:
		c.Logger().Errorf("%s: %+v \n %s", clientErr, logErr, errors.WithStack(errors.New(clientErr)).Error())
	}
	return errorResponse(c, ServerError, clientErr)
}

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func errorResponse(c buffalo.Context, code int, msg string) error {
	ret := response{
		Code: code,
		Msg:  msg,
	}

	return c.Render(http.StatusOK, r.JSON(ret))
}

func successResponse(c buffalo.Context, data any) error {
	ret := response{
		Code: SuccessResponse,
		Data: data,
	}
	return c.Render(http.StatusOK, r.JSON(ret))
}
