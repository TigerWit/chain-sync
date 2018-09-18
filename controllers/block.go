package controllers

import (
	"chain-sync/models"
	"fmt"
)

type RLBData struct {
	Status *Status         `json:"status"`
	Blocks []*models.Block `json:"blocks"`
}

func (m *MainController) GetLastBlocks() {
	rData := &RLBData{}
	status := &Status{}
	num, _ := m.GetInt("num")
	blocks, err := models.GetLastBlocks(num)
	if err != nil {
		status.Code = 500
		status.Msg = fmt.Sprintf("Query blocks fail: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}
	status.Success()
	rData.Status = status
	rData.Blocks = blocks
	m.Data["json"] = rData
	m.ServeJSON()
}
