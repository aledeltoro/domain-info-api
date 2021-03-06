package handler

import (
	"encoding/json"
	"fmt"
	"log"

	wrappedErr "domain-info-api/platform/errorhandling"

	"github.com/valyala/fasthttp"
)

// DomainGET returns the route handler for GET /domains
func (app *APP) DomainGET(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("Origin"))

	domains, customErr := app.GetAllDomains()
	if customErr != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	err := json.NewEncoder(ctx).Encode(domains)
	if err != nil {
		errMessage := fmt.Sprintf("JSON encoding failed: %s", err.Error())
		customErr := wrappedErr.New(fasthttp.StatusInternalServerError, "DomainGET", errMessage)
		log.Println(customErr)
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

}
