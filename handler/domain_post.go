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

type APP struct {
	*hostinfo.Connection
}

// DomainPOST returns the route handler for POST /domains
func (app *APP) DomainPOST(ctx *fasthttp.RequestCtx) {

	hostArg := ctx.URI().QueryArgs().Peek("host")

	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("Origin"))

	// facebook.com
	// facebook

	if !validator.IsURL(string(hostArg)) {
		customErr := wrappedErr.New(fasthttp.StatusBadRequest, "DomainPOST", "Invalid domain name")
		log.Println(customErr)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintln(ctx, customErr.Message.Error())
		return
	}

	domain, exists, customErr := app.CheckDomainExists(string(hostArg))
	if customErr != nil {

		switch customErr.Status {
		case fasthttp.StatusRequestTimeout:
			ctx.Response.SetStatusCode(fasthttp.StatusRequestTimeout)
			fmt.Fprintln(ctx, customErr.Message)
		case fasthttp.StatusInternalServerError:
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		case fasthttp.StatusNotImplemented:
			ctx.Response.SetStatusCode(fasthttp.StatusNotImplemented)
			fmt.Fprintln(ctx, customErr.Message)
		}

		return
	}

	if !exists {
		domain, customErr = hostinfo.NewDomain(string(hostArg))
		if customErr != nil {

			switch customErr.Status {
			case fasthttp.StatusRequestTimeout:
				ctx.Response.SetStatusCode(fasthttp.StatusRequestTimeout)
				fmt.Fprintln(ctx, customErr.Message)
			case fasthttp.StatusInternalServerError:
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			case fasthttp.StatusNotImplemented:
				ctx.Response.SetStatusCode(fasthttp.StatusNotImplemented)
				fmt.Fprintln(ctx, customErr.Message)
			}

			return
		}

		customErr = app.InsertDomain(domain)
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
		customErr := wrappedErr.New(fasthttp.StatusInternalServerError, "DomainPOST", errMessage)
		log.Println(customErr)
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

}
