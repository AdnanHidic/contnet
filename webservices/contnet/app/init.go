package app

import (
	"github.com/AdnanHidic/contnet/core/revel/binders"
	"github.com/AdnanHidic/contnet/core/revel/configreader"
	cors "github.com/AdnanHidic/contnet/core/revel/cors/app/controllers"
	ctrl "github.com/AdnanHidic/contnet/webservices/contnet/app/controllers"
	"github.com/revel/revel"
)

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	cors.SetCORSHeaders(c)

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// on app start hooks
	revel.OnAppStart(InitCors)
	revel.OnAppStart(InitContNet)

	// interceptors
	revel.InterceptMethod((*ctrl.App).SetupContext, revel.BEFORE)

	// Add null.* type binders
	binders.AddTypeBinders()
}

func InitCors() {
	cors.Init(configreader.StringFromConfig("cors.allowed"))
}

func InitContNet() {
	// load config

}
