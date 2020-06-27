package handler

import (
	hostinfo "domain-info-api/platform/hostinfo"
	"encoding/json"
	"fmt"
	"log"

	validator "github.com/asaskevich/govalidator"
	"github.com/valyala/fasthttp"
)

// DomainPOST returns the route handler for POST /domains
func DomainPOST(host *hostinfo.Connection) func(ctx *fasthttp.RequestCtx) {

	return func(ctx *fasthttp.RequestCtx) {

		hostArg := ctx.URI().QueryArgs().Peek("host")

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("Origin"))

		if !validator.IsURL(string(hostArg)) {
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			fmt.Fprintln(ctx, "Invalid domain name")
			return
		}

		domain, exists, err := host.CheckDomainExists(string(hostArg))
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		if !exists {
			domain, err = hostinfo.NewDomain(string(hostArg))
			if err != nil {
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			err = host.InsertDomain(domain)
			if err != nil {
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}

		ctx.Response.SetStatusCode(fasthttp.StatusCreated)
		ctx.Response.Header.SetContentType("application/json")

		err = json.NewEncoder(ctx).Encode(domain.HostInfo)
		if err != nil {
			log.Println("JSON encoding failed: ", err.Error())
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

	}

}
