package controllers

import (
	"github.com/AdnanHidic/contnet/contnet"
	base "github.com/AdnanHidic/contnet/core/revel/base/app/controllers"
	"github.com/guregu/null"
	"github.com/revel/revel"
	"net/http"
)

var Net *contnet.Net

type App struct {
	base.Base
}

func (c *App) Index() revel.Result {
	return c.RenderText("Hello from contnet server.")
}

func (c *App) NotFound() revel.Result {
	return c.Error(base.ERROR_NO_ACTION, http.StatusNotFound)
}

func (c *App) GetNthFrontpage(profileID, pageID null.Int) revel.Result {
	if !profileID.Valid || !pageID.Valid {
		return c.ErrorBadRequest()
	}

	output := Net.Select(profileID.Int64, uint8(pageID.Int64))
	return c.RenderJson(output)
}

func (c *App) GetDescription() revel.Result {
	output := Net.Describe()
	return c.RenderJson(output)
}

func (c *App) PostContent() revel.Result {
	// get input
	content := &contnet.Content{}
	if err := c.FromJson(content, nil); err != nil {
		return err
	}

	Net.SaveContent(content)
	return c.RenderText("OK.")
}

func (c *App) PostContentAction(contentID null.Int) revel.Result {
	// get input
	action := &contnet.Action{}
	if err := c.FromJson(action, nil); err != nil {
		return err
	}

	if err := Net.SaveAction(action); err != nil {
		return c.ErrorNotFound()
	}

	return c.RenderText("OK.")
}
