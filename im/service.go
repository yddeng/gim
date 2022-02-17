package im

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/gim/im/protocol"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
	"github.com/yddeng/utils/task"
	"math/rand"
	"net"
	"runtime"
	"sync/atomic"
	"time"
)

var crtConnCount int32

func StartTCPGateway(address string) error {
	log.Infof("start tcp gateway on address:%s. ", address)
	return dnet.ServeTCPFunc(address, func(conn net.Conn) {
		log.Debug("new tcp client", conn.RemoteAddr().String())
		atomic.AddInt32(&crtConnCount, 1)
		_ = createSession(conn)
	})
}

func StartWSGateway(address string) error {
	log.Infof("start ws gateway on address:%s. ", address)
	return dnet.ServeWSFunc(address, func(conn net.Conn) {
		log.Debug("new ws client", conn.RemoteAddr().String())
		_ = createSession(conn)
	})
}

func createSession(conn net.Conn) dnet.Session {
	return dnet.NewTCPSession(conn,
		dnet.WithTimeout(time.Second*time.Duration(config.HeartbeatTimeout), 0),
		dnet.WithCodec(Codec{}),
		//dnet.WithErrorCallback(func(session dnet.Session, err error) {
		//	fmt.Println("onError", err)
		//}),
		dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
			if err := taskQueue.Push(func() {
				dispatchMessage(session, data.(*Message))
			}); err != nil {
				log.Error(err)
			}
		}),
		dnet.WithCloseCallback(func(session dnet.Session, reason error) {
			atomic.AddInt32(&crtConnCount, -1)
			log.Debug(session.RemoteAddr().String(), reason)
			_ = taskQueue.Push(func() {
				ctx := session.Context()
				if ctx != nil {
					u := ctx.(*User)
					log.Infof("onClose user(%s) %s. ", u.ID, reason)
					if u.sess != session {
						return
					}
					u.sess.SetContext(nil)
					u.sess = nil
					u.notifyOnline()

				}
			})
		}))
}

var (
	msgHandler = map[uint16]func(*User, *Message){}
)

func registerHandler(cmd uint16, h func(*User, *Message)) {
	if _, ok := msgHandler[cmd]; ok {
		panic(fmt.Sprintf("handler cmd %d is alreadly register. ", cmd))
	}
	msgHandler[cmd] = h
}

func dispatchMessage(session dnet.Session, msg *Message) {
	cmd := msg.GetCmd()
	switch cmd {
	case uint16(protocol.CmdType_CmdHeartbeat):
		_ = session.Send(NewMessage(msg.GetSeq(), &protocol.Heartbeat{Timestamp: proto.Int64(time.Now().Unix())}))
	case uint16(protocol.CmdType_CmdUserLoginReq):
		onUserLogin(session, msg)
	default:
		if h, ok := msgHandler[cmd]; ok {
			ctx := session.Context()
			if ctx == nil {
				session.Close(errors.New("user is not login. "))
				return
			}
			h(ctx.(*User), msg)
		}
	}
}

func initLog(conf *Config) {
	logCfg := conf.LogConfig
	log.Infof("init logger debug=%v enableStdout%v. ", logCfg.Debug, logCfg.EnableStdout)
	if !logCfg.Debug {
		log.CloseDebug()
	}
	if !logCfg.EnableStdout {
		log.CloseStdOut()
	}

	if logCfg.MaxSize != 0 {
		log.SetOutput(logCfg.Path, logCfg.Filename, logCfg.MaxSize*1024*1024)
	}
}

func Service(cfgPath string) {
	rand.Seed(time.Now().UnixNano())

	log.Info("startup...")
	log.Info("load config. ")
	config = loadCfg(cfgPath)
	initLog(config)

	dbConfig := config.DBConfig
	log.Infof("init db sqlType=%s host=%s port=%d database=%s user=%s password=%s. ",
		dbConfig.SqlType, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password)
	var err error
	if err = dbInit(dbConfig.SqlType, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password); err != nil {
		panic(err)
	}

	log.Infof("init message deliver maxBackups=%d maxMessageCount=%d. ", config.MaxBackups, config.MaxMessageCount)
	messageDeliver, err = NewMessageDeliver(config.MaxBackups, config.MaxMessageCount)
	if err != nil {
		panic(err)
	}

	log.Infof("lru cache userCount=%d groupCount=%d. ", config.MaxConnCount*2, config.MaxConnCount)
	userCache = lru.New(config.MaxConnCount * 2)
	userCache.OnEvicted = onUserEvicted
	groupCache = lru.New(config.MaxConnCount)
	friendCache = lru.New(config.MaxConnCount * 2)

	log.Infof("task queue channel size is %d. ", config.MaxTaskCount)
	taskQueue = NewTaskQueue(config.MaxTaskCount * 2)
	go taskQueue.Run()

	taskPool = task.NewTaskPool(runtime.NumCPU(), config.MaxTaskCount)

	log.Infof("connection max count is %d. ", config.MaxConnCount)
	go func() {
		if err := StartTCPGateway(config.Address); err != nil {
			panic(err)
		}
	}()

}

func Stop() {
	log.Info("stopping... ")
}
