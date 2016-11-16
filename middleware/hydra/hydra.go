package hydra

import (
    "fmt"

	hydraSDK "github.com/ory-am/hydra/sdk"
    hydraClient "github.com/ory-am/hydra/client"
    "github.com/ory-am/hydra/firewall"
	"github.com/webx-top/echo"
)

var DefaultClient *hydraSDK.Client

func Options struct{
    Skipper echo.Skipper
    ClientID string
    ClientSecret string
    ClientURL string
}

func Connect(val Options) (hc *hydraSDK.Client,err error){
    hc, err = hydraSDK.Connect(
		hydraSDK.ClientID(val.ClientID),
		hydraSDK.ClientSecret(val.ClientSecret),
		hydraSDK.ClusterURL(val.ClusterURL),
	)
    return
}

func NewClient(hc *hydraSDK.Client,clientConfig *hydraClient.Client)(*hydraClient.Client, error){
    return hc.Clients.CreateClient(clientConfig)
}

func GetClient(hc *hydraSDK.Client,id string)(*hydraClient.Client, error){
    return hc.Clients.GetClient(id)
}

func GetContext(c echo.Context) *firewall.Context {
    ctx, _ := c.Get("hydra").(*firewall.Context)
    return ctx
}

func ScopesRequired(opt interface{},scopes ...string) echo.MiddlewareFunc {
    var hc *hydraSDK.Client
    var err error
    var skipper echo.Skipper
    if client,ok := opt.(*hydraSDK.Client); ok {
        hc = client
    } else if val,ok := opt.(*Options); ok {
        skipper = val.Skipper
	    hc, err = Connect(*val)
    } else if val,ok := opt.(Options); ok {
        skipper = val.Skipper
	    hc, err = Connect(val)
    } else if DefaultClient != nil {
        hc = DefaultClient
    } else {
        err = fmt.Errorf("invalid parameter: %T",opt)
    }
    
    if err != nil {
		panic(err.Error())
	}
    if DefaultClient == nil {
        DefaultClient = hc
    }
    if skipper == nil {
        skipper = echo.DefaultSkipper
    }

	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
            if skipper(c) {
                return h.Handle(c)
            }
			ctx, err := hc.Warden.TokenValid(c, hc.Warden.TokenFromRequest(c.Request().StdRequest()), scopes...)
			if err != nil {
				return err
			}
			// All required scopes are found
			c.Set("hydra", ctx)
			return h.Handle(c)
		})
	}
}
