package handler

import (
	"encoding/json"
	"fmt"
	"log"

	wrappedErr "domain-info-api/platform/errorhandling"
	hostinfo "domain-info-api/platform/hostinfo"

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

		domain, exists, customErr := host.CheckDomainExists(string(hostArg))
		if customErr != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		if !exists {
			domain, customErr = hostinfo.NewDomain(string(hostArg))
			if customErr != nil {
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			customErr = host.InsertDomain(domain)
			if customErr != nil {
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}

		ctx.Response.SetStatusCode(fasthttp.StatusCreated)
		ctx.Response.Header.SetContentType("application/json")

		err := json.NewEncoder(ctx).Encode(domain.HostInfo)
		if err != nil {
			errMessage := fmt.Sprintf("JSON encoding failed: %s", err.Error())
			customErr := wrappedErr.New(500, "DomainPOST", errMessage)
			log.Println(customErr)
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

	}

}
