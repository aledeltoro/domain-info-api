package handler

import (
	hostinfo "domain-info-api/platform/hostinfo"
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// DomainPOST returns the route handler for POST /domains
func DomainPOST(host *hostinfo.Connection) func(ctx *fasthttp.RequestCtx) {

	return func(ctx *fasthttp.RequestCtx) {

		hostArg := ctx.URI().QueryArgs().Peek("host")

		domain := hostinfo.NewDomain(string(hostArg))
		host.InsertDomain(domain)

		ctx.Response.SetStatusCode(http.StatusCreated)
		ctx.Response.Header.SetContentType("application/json")

		json.NewEncoder(ctx).Encode(domain.HostInfo)

	}

}
