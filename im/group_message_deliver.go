package im

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/protocol"
	"sort"
	"strings"
	"time"
)

const gcDur = 50

var messageDeliver *MessageDeliver

type MessageDeliver struct {
	maxBackups      int
	maxMessageCount int
	tableName       string
	tables          []string // 从旧到新的表名
	groupMessage    map[int64]*groupMessage
}

type groupMessage struct {
	messages map[int64]*protocol.MessageInfo
}

func (this *groupMessage) gc(maxCount int) {
	if len(this.messages) > maxCount+gcDur {
		ids := make([]int64, 0, len(this.messages))
		for id := range this.messages {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool {
			return ids[i] < ids[j]
		})
		off := len(ids) - maxCount
		for i := 0; i < off; i++ {
			delete(this.messages, ids[i])
		}
	}
}

func NewMessageDeliver(maxBackups, maxMessageCount int) (*MessageDeliver, error) {
	this := &MessageDeliver{
		maxBackups:      maxBackups,
		maxMessageCount: maxMessageCount,
		groupMessage:    map[int64]*groupMessage{},
	}

	tables, err := this.dbLoadAllMessageList()
	if err != nil {
		return nil, err
	}

	tableName := this.makeTableName()
	this.tables = tables
	this.tableName = tableName

	exist := false
	for _, name := range this.tables {
		if name == tableName {
			exist = true
			break
		}
	}
	if !exist {
		if err := this.dbCreateTableMessage(tableName); err != nil {
			return nil, err
		} else if err = this.dbSetMessageList(tableName); err != nil {
			return nil, err
		}
		this.tables = append(this.tables, tableName)
	}

	sort.Strings(this.tables)
	if len(this.tables) > this.maxBackups {
		off := len(this.tables) - this.maxBackups
		delList := this.tables[:off]
		if err := this.dbDelMessageList(delList); err == nil {
			this.tables = this.tables[off:]
			_ = this.dbDropTableMessage(delList)
		}
	}

	return this, nil
}

func (this *MessageDeliver) makeTableName() string {
	return "message_" + time.Now().Format("20060102")
}

func (this *MessageDeliver) pushMessage(groupID int64, msg *protocol.MessageInfo) error {
	tableName := this.makeTableName()
	if this.tableName != tableName {
		if err := this.dbCreateTableMessage(tableName); err != nil {
			return err
		} else if err = this.dbSetMessageList(tableName); err != nil {
			return err
		}
		this.tableName = tableName
		this.tables = append(this.tables, tableName)

		if len(this.tables) > this.maxBackups {
			off := len(this.tables) - this.maxBackups
			delList := this.tables[:off]
			if err := this.dbDelMessageList(delList); err == nil {
				this.tables = this.tables[off:]
				_ = this.dbDropTableMessage(delList)
			}
		}
	}

	if err := this.dbSetNxMessage(groupID, msg, this.tableName); err != nil {
		return err
	}

	gm, ok := this.groupMessage[groupID]
	if !ok {
		gm = &groupMessage{
			messages: make(map[int64]*protocol.MessageInfo, this.maxMessageCount),
		}
		this.groupMessage[groupID] = gm
	}
	gm.messages[msg.GetMsgID()] = msg
	gm.gc(this.maxMessageCount)

	return nil
}

func (this *MessageDeliver) loadMessage(groupID int64, ids []int64) ([]*protocol.MessageInfo, error) {
	gm, ok := this.groupMessage[groupID]
	if !ok {
		gm = &groupMessage{
			messages: make(map[int64]*protocol.MessageInfo, this.maxMessageCount),
		}
		this.groupMessage[groupID] = gm
	}

	infos := make([]*protocol.MessageInfo, 0, len(ids))
	finds := make([]int64, 0, len(ids))
	for _, id := range ids {
		if v, ok := gm.messages[id]; ok {
			infos = append(infos, v)
		} else {
			finds = append(finds, id)
		}
	}

	if len(finds) > 0 {
		for i := len(this.tables) - 1; i >= 0; i-- {
			if ins, err := this.dbLoadMessageBatch(groupID, finds, this.tables[i]); err != nil {
				return nil, err
			} else {
				infos = append(infos, ins...)
				if len(infos) >= len(ids) {
					break
				}
			}
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].GetMsgID() < infos[j].GetMsgID()
	})
	for _, v := range infos {
		gm.messages[v.GetMsgID()] = v
	}
	gm.gc(this.maxMessageCount)

	return infos, nil
}

/***********  db  *************/

func (this *MessageDeliver) dbLoadAllMessageList() ([]string, error) {
	rows, err := sqlDB.Query(`SELECT table_name FROM "message_list";`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := make([]string, 0, 8)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		tables = append(tables, name)
	}
	return tables, nil
}

func (this *MessageDeliver) dbDelMessageList(tables []string) error {
	sqlStr := `
DELETE FROM "message_list" 
WHERE %s;`
	keys := make([]string, 0, len(tables))
	for _, name := range tables {
		keys = append(keys, fmt.Sprintf("table_name = '%s'", name))
	}
	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(keys, " OR "))
	//log.Debug(sqlStatement)
	_, err := sqlDB.Exec(sqlStatement)
	return err
}

func (this *MessageDeliver) dbSetMessageList(tableName string) error {
	sqlStatement := `
INSERT INTO "message_list" (table_name)  
VALUES ($1);`
	_, err := sqlDB.Exec(sqlStatement, tableName)
	return err
}

func (this *MessageDeliver) dbCreateTableMessage(tableName string) error {
	sqlStr := `DROP TABLE IF EXISTS "%s";
CREATE TABLE "%s" (
    "id"           varchar(255) NOT NULL,
    "group_id"      int8 NOT NULL ,
    "message_id"   int8 NOT NULL,
    "message"      bytea NOT NULL,
    PRIMARY KEY ("id")
);`
	sqlStatement := fmt.Sprintf(sqlStr, tableName, tableName)
	_, err := sqlDB.Exec(sqlStatement)
	return err
}

func (this *MessageDeliver) dbDropTableMessage(tables []string) error {
	sqlStr := `DROP TABLE IF EXISTS "%s";`
	sqlStatement := ""
	for _, name := range tables {
		sqlStatement += fmt.Sprintf(sqlStr, name)
	}
	_, err := sqlDB.Exec(sqlStatement)
	return err
}

func (this *MessageDeliver) dbSetNxMessage(groupID int64, msg *protocol.MessageInfo, tableName string) error {
	sqlStr := `
INSERT INTO "%s" (id,group_id,message_id,message)  
VALUES ($1,$2,$3,$4)
ON conflict(id) DO 
UPDATE SET message = $4;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	id := fmt.Sprintf("%d_%d", groupID, msg.GetMsgID())
	data, _ := proto.Marshal(msg)
	_, err := sqlDB.Exec(sqlStatement, id, groupID, msg.GetMsgID(), data)
	return err
}

func (this *MessageDeliver) dbLoadMessageBatch(groupID int64, ids []int64, tableName string) ([]*protocol.MessageInfo, error) {
	sqlStr := `
SELECT message FROM "%s" 
WHERE %s;`

	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf("id = '%d_%d'", groupID, id))
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(keys, " OR "))
	//log.Debug(sqlStatement)

	rows, err := sqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	infos := make([]*protocol.MessageInfo, 0, len(ids))
	defer rows.Close()
	for rows.Next() {
		var info protocol.MessageInfo
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			return nil, err
		}

		_ = proto.Unmarshal(data, &info)
		infos = append(infos, &info)
	}

	return infos, nil
}
