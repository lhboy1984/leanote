package tests

import (
	"github.com/lhboy1984/leanote/app/db"
	"testing"
	//	. "github.com/lhboy1984/leanote/app/lea"
	//	"github.com/lhboy1984/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func TestDBConnect(t *testing.T) {
	db.Init("mongodb://localhost:27017/leanote", "leanote")
}
