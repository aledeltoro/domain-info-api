package handler

import (
	hostinfo "domain-info-api/platform/hostinfo"
	"encoding/json"
	"log"
	"net/http"

	validator "github.com/asaskevich/govalidator"
	"github.com/valyala/fasthttp"
)

// DomainPOST returns the route handler for POST /domains
func DomainPOST(host *hostinfo.Connection) func(ctx *fasthttp.RequestCtx) {

	return func(ctx *fasthttp.RequestCtx) {

		hostArg := ctx.URI().QueryArgs().Peek("host")

		if !validator.IsURL(string(hostArg)) {
			ctx.Error("Invalid domain name", http.StatusBadRequest)
			return
		}

		domain, exists, err := host.CheckDomainExists(string(hostArg))
		if err != nil {
			ctx.Error("", http.StatusInternalServerError)
			return
		}

		if !exists {
			domain, err = hostinfo.NewDomain(string(hostArg))
			if err != nil {
				ctx.Error("", http.StatusInternalServerError)
				return
			}
			err = host.InsertDomain(domain)
			if err != nil {
				ctx.Error("", http.StatusInternalServerError)
				return
			}
		}

		ctx.Response.SetStatusCode(http.StatusCreated)
		ctx.Response.Header.SetContentType("application/json")

		err = json.NewEncoder(ctx).Encode(domain.HostInfo)
		if err != nil {
			log.Println("JSON encoding failed: ", err.Error())
			ctx.Error("", http.StatusInternalServerError)
			return
		}

	}

}
