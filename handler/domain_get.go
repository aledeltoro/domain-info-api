package handler

import (
	hostinfo "domain-info-api/platform/hostinfo"
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

// DomainGET returns the route handler for GET /domains
func DomainGET(host *hostinfo.Connection) func(ctx *fasthttp.RequestCtx) {

	return func(ctx *fasthttp.RequestCtx) {

		domains, err := host.GetAllDomains()
		if err != nil {
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("Origin"))
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.SetStatusCode(fasthttp.StatusOK)

		err = json.NewEncoder(ctx).Encode(domains)
		if err != nil {
			log.Println("JSON encoding failed: ", err.Error())
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}

	}

}
