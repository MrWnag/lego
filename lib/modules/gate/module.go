package gate

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/proto"
)

type Gate struct {
	cbase.ModuleBase
	LocalRouteMgrComp  ILocalRouteMgrComp
	RemoteRouteMgrComp IRemoteRouteMgrComp
	AgentMgrComp       IAgentMgrComp
	//WsServerComp       *WsServerComp
	//TcpServerComp      *TcpServerComp
}

func (this *Gate) UnRegisterRemoteRoute(comId uint16, sType, sId string) {
	this.RemoteRouteMgrComp.UnRegisterRoute(comId, sType, sId)
}
func (this *Gate) UnRegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code int, err string)) {
	this.LocalRouteMgrComp.UnRegisterRoute(comId, f)
}

//需重构处理  内部函数为重构代码
//代理链接
func (this *Gate) Connect(a IAgent) {
	this.AgentMgrComp.Connect(a)
}

//代理关闭
func (this *Gate) DisConnect(a IAgent) {
	this.AgentMgrComp.DisConnect(a)
}

//接收代理消息
func (this *Gate) OnRoute(a IAgent, msg proto.IMessage) (code int, err string) {
	if this.LocalRouteMgrComp.IsHaveRoute(msg.GetComId()) {
		this.LocalRouteMgrComp.OnRoute(a, msg)
	}
	if this.RemoteRouteMgrComp.IsHaveRoute(msg.GetComId()) {
		this.RemoteRouteMgrComp.OnRoute(a, msg)
	}
	return
}

//主动关闭代理
func (this *Gate) CloseAgent(sId string) (result string, err string) {
	return this.AgentMgrComp.Close(sId)
}

//发送代理消息
func (this *Gate) SendMsg(sId string, msg proto.IMessage) (result int, err string) {
	return this.AgentMgrComp.SendMsg(sId, msg)
}

//广播代理消息
func (this *Gate) RadioMsg(sIds []string, msg proto.IMessage) (result int, err string) {
	for _, v := range sIds {
		this.AgentMgrComp.SendMsg(v, msg)
	}
	return
}

//func (this *Gate) OnInstallComp() {
//	this.ModuleBase.OnInstallComp()
//	this.AgentMgrComp = this.RegisterComp(new(AgentMgrComp)).(*AgentMgrComp)
//	this.LocalRouteMgrComp = this.RegisterComp(new(LocalRouteMgrComp)).(*LocalRouteMgrComp)
//	this.RemoteRouteMgrComp = this.RegisterComp(new(LocalRouteMgrComp)).(*RemoteRouteMgrComp)
//	this.TcpServerComp = this.RegisterComp(new(TcpServerComp)).(*TcpServerComp)
//	this.WsServerComp = this.RegisterComp(new(WsServerComp)).(*WsServerComp)
//	this.LocalRouteMgrComp.NewSession = NewLocalSession
//	this.RemoteRouteMgrComp.NewSession = NewRemoteSession
//}
