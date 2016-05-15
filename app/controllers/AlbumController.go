package controllers

import (
	"github.com/lhboy1984/leanote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// Album controller
type Album struct {
	BaseController
}

// 图片管理, iframe
func (c Album) Index() revel.Result {
	c.SetLocale()
	return c.RenderTemplate("album/index.html")
}

// all albums by userId
func (c Album) GetAlbums() revel.Result {
	re := albumService.GetAlbums(c.GetUserId())
	return c.RenderJson(re)
}
func (c Album) DeleteAlbum(albumId string) revel.Result {
	re, msg := albumService.DeleteAlbum(c.GetUserId(), albumId)
	return c.RenderJson(info.Re{Ok: re, Msg: msg})
}

// add album
func (c Album) AddAlbum(name string) revel.Result {
	album := info.Album{
		AlbumId: bson.NewObjectId(),
		Name:    name,
		Seq:     -1,
		UserId:  c.GetObjectUserId()}
	re := albumService.AddAlbum(album)

	if re {
		return c.RenderJson(album)
	} else {
		return c.RenderJson(false)
	}
}

// update alnum name
func (c Album) UpdateAlbum(albumId, name string) revel.Result {
	return c.RenderJson(albumService.UpdateAlbum(albumId, c.GetUserId(), name))
}
