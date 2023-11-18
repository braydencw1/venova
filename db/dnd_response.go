package db

import (
	"fmt"
)

type DndResponses struct {
	Id       int `gorm:"primaryKey"`
	Response string
}

func DndMsgResponse() string {
	var res DndResponses
	query := db.Order("RANDOM()").Find(&res)
	if query.Error != nil {
		erMsg := fmt.Sprintf("Error, %s", query.Error)
		return erMsg
	}
	return res.Response

}
