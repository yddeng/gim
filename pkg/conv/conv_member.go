package conv

import (
	"github.com/yddeng/gim/internal/codec"
	"github.com/yddeng/gim/internal/protocol/pb"
	"github.com/yddeng/gim/pkg/gate"
	"github.com/yddeng/gim/pkg/user"
	"github.com/yddeng/utils/log"
)

func onAddMember(u *user.User, message *codec.Message) {
	req := message.GetData().(*pb.AddMemberReq)
	log.Debugf("onAddMember %v", req)

}

func init() {
	gate.RegisterHandler(uint16(pb.CmdType_CmdAddMemberReq), onAddMember)
}
