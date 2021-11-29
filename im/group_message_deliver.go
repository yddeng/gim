package im

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/internal/db"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/utils/log"
	"strings"
	"time"
)

var messageDeliver *MessageDeliver

type MessageDeliver struct {
	tableName  string
	maxBackups int      // 日志保留时长，单位：天
	tables     []string // 从旧到新的表名
}

func NewMessageDeliver(maxBackups int) (*MessageDeliver, error) {
	this := &MessageDeliver{maxBackups: maxBackups}

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
		if err := this.dbCreateMessageTable(tableName); err != nil {
			return nil, err
		} else if err = this.dbSetMessageList(tableName); err != nil {
			return nil, err
		}
		this.tables = append(this.tables, tableName)
	}

	if len(this.tables) > this.maxBackups {
		off := len(this.tables) - this.maxBackups
		if err := this.dbDelMessageList(this.tables[0 : off-1]); err == nil {
			this.tables = this.tables[off:]
		}
	}

	return this, nil
}

func (this *MessageDeliver) makeTableName() string {
	return "message_" + time.Now().Format("20060102")
}

func (this *MessageDeliver) pushMessage(groupID int64, msg *pb.MessageInfo) error {
	tableName := this.makeTableName()
	if this.tableName != tableName {
		if err := this.dbCreateMessageTable(tableName); err != nil {
			return err
		} else if err = this.dbSetMessageList(tableName); err != nil {
			return err
		}
		this.tableName = tableName
		this.tables = append(this.tables, tableName)

		if len(this.tables) > this.maxBackups {
			off := len(this.tables) - this.maxBackups
			if err := this.dbDelMessageList(this.tables[0 : off-1]); err == nil {
				this.tables = this.tables[off:]
			}
		}
	}

	if err := this.dbSetNxMessage(groupID, msg, this.tableName); err != nil {
		return err
	}
	return nil
}

func (this *MessageDeliver) loadMessage() {

}

/***********  db  *************/

func (this *MessageDeliver) dbLoadAllMessageList() ([]string, error) {
	rows, err := db.SqlDB.Query(`SELECT table_name FROM "message_list";`)
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
	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(tables, " OR "))
	//log.Debug(sqlStatement)
	_, err := db.SqlDB.Exec(sqlStatement)
	return err
}

func (this *MessageDeliver) dbSetMessageList(tableName string) error {
	sqlStatement := `
INSERT INTO "message_list" (table_name)  
VALUES ($1);`
	_, err := db.SqlDB.Exec(sqlStatement, tableName)
	return err
}

func (this *MessageDeliver) dbCreateMessageTable(tableName string) error {
	sqlStr := `DROP TABLE IF EXISTS "%s";
CREATE TABLE "%s" (
    "id"           varchar(255) NOT NULL,
    "group_id"      int8 NOT NULL ,
    "message_id"   int8 NOT NULL,
    "message"      bytea NOT NULL,
    PRIMARY KEY ("id")
);`
	sqlStatement := fmt.Sprintf(sqlStr, tableName, tableName)
	_, err := db.SqlDB.Exec(sqlStatement)
	return err
}

func (this *MessageDeliver) dbSetNxMessage(groupID int64, msg *pb.MessageInfo, tableName string) error {
	sqlStr := `
INSERT INTO "%s" (id,group_id,message_id,message)  
VALUES ($1,$2,$3,$4)
ON conflict(id) DO 
UPDATE SET message = $4;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	id := fmt.Sprintf("%d_%d", groupID, msg.GetMsgID())
	data, _ := proto.Marshal(msg)
	_, err := db.SqlDB.Exec(sqlStatement, id, groupID, msg.GetMsgID(), data)
	return err
}

func (this *MessageDeliver) dbLoadMessageBatch(groupID int64, start, limit int, tableName string) ([]*pb.MessageInfo, error) {
	sqlStr := `
SELECT message FROM "%s" 
WHERE %s;`

	keys := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		seq := start + i
		keys = append(keys, fmt.Sprintf("id = '%d_%d'", groupID, seq))
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(keys, " OR "))
	log.Debug(sqlStatement)

	rows, err := db.SqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	infos := make([]*pb.MessageInfo, 0, limit)
	defer rows.Close()
	for rows.Next() {
		var info pb.MessageInfo
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
