package im

import (
	"errors"
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/utils/log"
	"github.com/yddeng/utils/lru"
	"github.com/yddeng/utils/task"
	"math/rand"
	"net"
	"runtime"
	"time"
)

func StartTCPGateway(address string) error {
	log.Infof("start tcp gateway on address:%s. ", address)
	return dnet.ServeTCPFunc(address, func(conn net.Conn) {
		log.Debug("new tcp client", conn.RemoteAddr().String())
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
		//dnet.WithTimeout(time.Second*5, 0), // 超时
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
			log.Debug(session.RemoteAddr().String(), reason)
			ctx := session.Context()
			if ctx != nil {
				u := ctx.(*User)
				log.Infof("onClose user(%s) %s. ", u.ID, reason)
				u.sess.SetContext(nil)
				u.sess = nil
			}
		}))
}

var (
	userHandler  = map[uint16]func(session dnet.Session, msg *Message){}
	groupHandler = map[uint16]func(*User, *Message){}
)

func registerGroupHandler(cmd uint16, h func(*User, *Message)) {
	if _, ok := groupHandler[cmd]; ok {
		panic(fmt.Sprintf("group handler cmd %d is alreadly register. ", cmd))
	}
	groupHandler[cmd] = h
}

func registerUserHandler(cmd uint16, h func(session dnet.Session, msg *Message)) {
	if _, ok := userHandler[cmd]; ok {
		panic(fmt.Sprintf("user handler cmd %d is alreadly register. ", cmd))
	}
	userHandler[cmd] = h
}

func dispatchMessage(session dnet.Session, msg *Message) {
	cmd := msg.GetCmd()
	msgType := msg.GetType()

	switch msgType {
	case MESSAGE_UESR:
		if h, ok := userHandler[cmd]; ok {
			h(session, msg)
		}
	case MESSAGE_GROUP:
		if h2, ok := groupHandler[cmd]; ok {
			ctx := session.Context()
			if ctx == nil {
				session.Close(errors.New("user is not login. "))
				return
			}
			h2(ctx.(*User), msg)
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

	//log.SetOutput(logCfg.Path, logCfg.Filename, logCfg.MaxSize*1024*1024)
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

	log.Infof("lru cache userCount=%d groupCount=%d. ", config.UserCacheCount, config.GroupCacheCount)
	userCache = lru.New(config.UserCacheCount)
	groupCache = lru.New(config.GroupCacheCount)

	log.Infof("task count=%d. ", config.MaxTaskCount)
	taskQueue = NewTaskQueue(config.MaxTaskCount * 2)
	go taskQueue.Run()

	taskPool = task.NewTaskPool(runtime.NumCPU(), config.MaxTaskCount)

	go func() {
		if err := StartTCPGateway(config.Address); err != nil {
			panic(err)
		}
	}()

}

func Stop() {
	log.Info("stopping... ")
}
