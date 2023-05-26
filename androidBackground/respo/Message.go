package respo

type Message struct {
	Id      int    `json:"-"`
	Role    string `json:"role"`
	Content string `gorm:"column:message" json:"content"` //与message映射
}

func (m Message) TableName() string {
	return "admin"
}

func (m Message) AddMessage() bool {
	//向数据库种插入一条数据
	tx := GolbalDB.Create(&m)
	if tx.Error != nil {
		return false
	}
	return true
}

func (m Message) GetMessage() []Message {
	//获取最新的(id最大)四条数据
	var messages []Message
	tx := GolbalDB.Order("id desc").Limit(4).Find(&messages)
	if tx.Error != nil {
		return nil
	}
	return messages
}
