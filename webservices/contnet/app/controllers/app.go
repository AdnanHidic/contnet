package contnet

import (
	base "github.com/AdnanHidic/contnet/core/revel/base/app/controllers"
	"github.com/guregu/null"
	"github.com/revel/revel"
	"net/http"
)

type App struct {
	base.Base
}

func (c *App) Index() revel.Result {
	return c.RenderText("Hello from contnet-backend.")
}

func (c *App) NotFound() revel.Result {
	return c.Error(base.ERROR_NO_ACTION, http.StatusNotFound)
}

func (c *App) SetupContext() revel.Result {
	return nil
}

func (c *App) GetContent(forID, limit, offset null.Int) revel.Result {
	return c.ErrorNotImplemented()
}

func (c *App) GetDescription() revel.Result {
	return c.ErrorNotImplemented()
}

func (c *App) PostContent() revel.Result {
	return c.ErrorNotImplemented()
}

func (c *App) PostContentRead(contentID null.Int) revel.Result {
	return c.ErrorNotImplemented()
}
