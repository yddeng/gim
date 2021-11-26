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

var (
	createMessageTableStr = `DROP TABLE IF EXISTS "%s";
CREATE TABLE "%s" (
    "id"           varchar(255) NOT NULL,
    "conv_id"      int8 NOT NULL ,
    "message_id"   int8 NOT NULL,
    "message"      bytea NOT NULL,
    PRIMARY KEY ("id")
);`

	tableName string
)

func InitMessageTable() {
	tableName = makeMessageTableName()
	if exist := existMessageTable(tableName); !exist {
		createMessageTable(tableName)
	}
}

func makeMessageTableName() string {
	date := time.Now().Format("20060102")
	return "message_" + date
}

func existMessageTable(tableName string) bool {
	sqlStr := `
select count(*) from "%s";`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	smt, err := db.SqlDB.Prepare(sqlStatement)
	if err != nil {
		return false
	}
	row := smt.QueryRow()
	var count int
	err = row.Scan(&count)
	if err != nil {
		return false
	}
	return true
}

func createMessageTable(tableName string) error {
	sqlStatement := fmt.Sprintf(createMessageTableStr, tableName, tableName)
	_, err := db.SqlDB.Exec(sqlStatement)
	return err
}

func setNxMessage(convID int64, msg *pb.MessageInfo, tableName string) error {
	sqlStr := `
INSERT INTO "%s" (id,conv_id,message_id,message)  
VALUES ($1,$2,$3,$4)
ON conflict(id) DO 
UPDATE SET message = $4;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	id := fmt.Sprintf("%d_%d", convID, msg.GetMsgID())
	data, _ := proto.Marshal(msg)
	_, err := db.SqlDB.Exec(sqlStatement, id, convID, msg.GetMsgID(), data)
	return err
}

func loadMessageBatch(convID int64, start, limit int, tableName string) ([]*pb.MessageInfo, error) {
	sqlStr := `
SELECT message FROM "%s" 
WHERE %s;`

	keys := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		seq := start + i
		keys = append(keys, fmt.Sprintf("id = '%d_%d'", convID, seq))
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
