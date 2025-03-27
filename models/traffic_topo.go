package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
	"strings"
)

type TrafficTopoType string

func (t TrafficTopoType) String() string {
	return string(t)
}

const (
	TrafficTopoType_Spine TrafficTopoType = "spine"
	TrafficTopoType_Leaf  TrafficTopoType = "leaf"
	TrafficTopoType_Node  TrafficTopoType = "node"
)

type TrafficTopo struct {
	ID        uint                 `json:"id" gorm:"primarykey"`
	IP        string               `json:"ip" gorm:"uniqueIndex:topo_ip;size:255"`
	Type      TrafficTopoType      `json:"type" gorm:"size:32"`
	Connects  ConnectSlice         `json:"connects" gorm:"type:json"`
	InPromql  string               `json:"-" gorm:"type:text"`
	OutPromql string               `json:"-" gorm:"type:text"`
	Labels    LabelExpressionSlice `json:"-" gorm:"type:json"`
	CreatedAt DateTime             `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt DateTime             `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP"`
	gorm.DeletedAt
}

func TrafficTopoListAll(ctx *ctx.Context) ([]*TrafficTopo, error) {
	ret := make([]*TrafficTopo, 0)

	db := DB(ctx).Where("id > ?", 0)

	err := db.Find(&ret).Error

	return ret, err
}

type LabelExpression struct {
	Label    string `json:"label"`    // 标签名
	Operator string `json:"operator"` // 运算符:(=/=~)
	Value    string `json:"value"`    //  标签值
}

func (l LabelExpression) ToString() (string, error) {
	if l.Label == "" {
		return "", fmt.Errorf("label can't be empty")
	}

	if l.Operator == "" {
		return "", fmt.Errorf("label [%s]'s operator  is empty", l.Label)
	}

	if l.Value == "" {
		return "", fmt.Errorf("label [%s]'s value is empty", l.Label)
	}
	return fmt.Sprintf(`%s%s"%s"`, l.Label, l.Operator, l.Value), nil
}

type LabelExpressionSlice []LabelExpression

func (l LabelExpressionSlice) KeyToLabelExpression() map[string]*LabelExpression {
	if len(l) == 0 {
		return nil
	}

	ret := make(map[string]*LabelExpression)
	for _, v := range l {
		ret[v.Label] = &v
	}

	return ret
}

func ConvertKeyToLabelExpressionToString(mm map[string]*LabelExpression) (string, error) {
	if len(mm) == 0 {
		return "", nil
	}

	var ret string
	for _, v := range mm {
		str, err := v.ToString()
		if err != nil {
			return "", err
		}
		ret += str + ","
	}

	return strings.TrimRight(ret, ","), nil
}

func (l *LabelExpressionSlice) Scan(value interface{}) error {
	if l == nil {
		*l = make(LabelExpressionSlice, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, l)
}

// Value 实现 driver.Valuer 接口，存储数据时序列化成 JSON
func (l LabelExpressionSlice) Value() (driver.Value, error) {
	if l == nil {
		return nil, nil
	}

	return json.Marshal(l)
}

// GormDataType 声明数据库类型（可选，增强可读性）
func (LabelExpressionSlice) GormDataType() string {
	return "json"
}

type Connect struct {
	ParentIP       string `json:"parentIP"`
	SelfPortDevice string `json:"selfPortDevice"`
}
type ConnectSlice []Connect

func (c *ConnectSlice) Scan(value interface{}) error {
	if c == nil {
		*c = make(ConnectSlice, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

// Value 实现 driver.Valuer 接口，存储数据时序列化成 JSON
func (c ConnectSlice) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

// GormDataType 声明数据库类型（可选，增强可读性）
func (ConnectSlice) GormDataType() string {
	return "json"
}
