package db

import (
	"fmt"
)

type DndResponses struct {
	Id       int `gorm:"primaryKey"`
	Response string
}

func DndMsgResponse() string {
	var query DndResponses
	res := db.Order("RANDOM()").Find(&query)
	if res.Error != nil {
		err := fmt.Sprintf("Error, %s", res.Error)
		return err
	}
	return query.Response

}
