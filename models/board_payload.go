package models

import (
	"encoding/json"
	"errors"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

type BoardPayload struct {
	Id      int64  `json:"id" gorm:"primaryKey"`
	Payload string `json:"payload"`
}

func (p *BoardPayload) TableName() string {
	return "board_payload"
}

func (p *BoardPayload) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	return DB(ctx).Model(p).Select(selectField, selectFields...).Updates(p).Error
}

func BoardPayloadGets(ctx *ctx.Context, ids []int64) ([]*BoardPayload, error) {
	if len(ids) == 0 {
		return nil, errors.New("empty ids")
	}

	var arr []*BoardPayload
	err := DB(ctx).Where("id in ?", ids).Find(&arr).Error
	return arr, err
}

func BoardPayloadGet(ctx *ctx.Context, id int64) (string, error) {
	payloads, err := BoardPayloadGets(ctx, []int64{id})
	if err != nil {
		return "", err
	}

	if len(payloads) == 0 {
		return "", nil
	}

	return payloads[0].Payload, nil
}

func BoardPayloadSave(ctx *ctx.Context, id int64, payload string) error {
	var bp BoardPayload
	err := DB(ctx).Where("id = ?", id).Find(&bp).Error
	if err != nil {
		return err
	}

	if bp.Id > 0 {
		// already exists
		bp.Payload = payload
		return bp.Update(ctx, "payload")
	}

	return Insert(ctx, &BoardPayload{
		Id:      id,
		Payload: payload,
	})
}

func GetPanelContent(ctx *ctx.Context, id uint, panelID string) (*Panel, error) {
	var result struct {
		PanelContent string `gorm:"column:panel_content"`
	}

	err := DB(ctx).Table("board_payload").
		Select("JSON_EXTRACT(payload, REPLACE(JSON_UNQUOTE(JSON_SEARCH(payload, 'one', ?, NULL, '$.panels[*].id')), '.id', '')) as panel_content",
			panelID).Where("id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	var ret = &Panel{}
	if err = json.Unmarshal([]byte(result.PanelContent), ret); err != nil {
		return nil, err
	}

	return ret, nil
}

type Panel struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Type            string           `json:"type"`
	Links           []Link           `json:"links"`
	Custom          Custom           `json:"custom"`
	Layout          Layout           `json:"layout"`
	Options         Options          `json:"options"`
	Targets         []Tgt            `json:"targets"`
	Version         string           `json:"version"`
	MaxPerRow       int              `json:"maxPerRow"`
	Description     string           `json:"description"`
	DatasourceCate  string           `json:"datasourceCate"`
	DatasourceValue string           `json:"datasourceValue"`
	Transformations []Transformation `json:"transformations"`
}

type Link struct {
	// Add fields as needed
}

type Custom struct {
	Calc       string `json:"calc"`
	TextMode   string `json:"textMode"`
	ValueField string `json:"valueField"`
}

type Layout struct {
	H           int    `json:"h"`
	I           string `json:"i"`
	W           int    `json:"w"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	IsResizable bool   `json:"isResizable"`
}

type Options struct {
	Thresholds      Thresholds      `json:"thresholds"`
	StandardOptions StandardOptions `json:"standardOptions"`
}

type Thresholds struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Type  string   `json:"type"`
	Color string   `json:"color"`
	Value *float64 `json:"value"` // Using pointer to handle null values
}

type StandardOptions struct {
	Max      float64 `json:"max"`
	Min      float64 `json:"min"`
	Util     string  `json:"util"`
	Decimals int     `json:"decimals"`
}

type Tgt struct {
	Expr          string `json:"expr"`
	RefID         string `json:"refId"`
	MaxDataPoints int    `json:"maxDataPoints"`
}

type Transformation struct {
	ID      string                 `json:"id"`
	Options map[string]interface{} `json:"options"`
}
