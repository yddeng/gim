# chat

游戏中的聊天服务器

## 对话类型

### 聊天室 

如世界频道

- 群聊人数非常多，群人员上限无限大（5000）
- 在线用户均可参与聊天，无需申请加入群聊，随时加入，随时退出
- 无固定成员列表，成员之间可能是无好友关系的陌生人
- 离线后不需要收到消息推送

### 普通对话

玩家间的单聊、私聊 或者 群聊都属于普通对话，只不过单聊的成员（members）人数为2，群聊成员人数则大于 2。

- 有人数限制，一个群聊中最大支持 500 个成员
- 可查成员列表、成员在线状态，关注成员信息
- 可查历史消息记录
- 允许增加、删除成员
- 成员离线状态下，可以收到消息推送通知，且上线后会进行消息同步


### 临时对话

类似组队聊天，副本内的聊天

- 不需要保存聊天记录，不需要查询历史消息
- 玩家若突然离线，无需发离线消息推送
-「房间」具有一次性，游戏结束「房间」即消失
- 聊天参与的人数较少（最多 8 人组队）


## 离线消息同步实现

客户端主动拉取，即使切换了设备也能拉取到之前的消息。
同时应将消息保存到本地，以减少客户端对云端消息记录的查询次数，并且在设备离线情况下，也能展示出部分数据。

所有消息都维护在会话中，邀请、踢人、加入、离开、聊天消息等。
聊天室记录每条消息并打上唯一ID（自增），保留一段时间或一定数目的消息。
可能用户半月没有登陆，之前的消息已经被丢弃，对于游戏聊天系统来看影响不大。用户身上记录每个聊天已读的最大ID。

1. 邀请加入会话、被踢出会话 消息

每个用户维持一个自己的会话列表。被邀请或踢出，直接增删对应的会话。
用户登陆后 发送获取会话列表协议，比较本地的列表。
       
2. 会话消息同步

用户根据会话本地已保存的最大消息ID，请求服务器拉取未读消息。当服务器与本地最大ID一致时消息同步完成。
tips： 新加入的会话最大消息ID为会话当前ID。从新到旧依此拉取。

用户在线的情况下，服务器主动同步消息。不需要前端拉取。


## 消息记录

由配置表确定保存几天的记录。默认3天。

根据日期的创建每日的消息表 例： 2021.11.14 = message_20211114。
表里保存所有的会话记录，id规则为 conversationID+messageID 例： 会话12的23条消息 = 12_13.
每日刷新后会删除记录外的表，即所有的会话以前的记录都会删除。


## 流程
用户本地应保存一些数据
- 好友申请，登陆后获取到的好友申请列表与本地做对比，确认哪些室新申请提示用户
- 每个聊天室已读消息索引，由用户来选择加载未读消息、或者忽略


1. 登陆 userLogin 
2. 获取好友列表 getFriend(包含好友申请)、获取聊天室 getGroupList（包含当前聊天室最新消息索引）。
3. 获取用户的基本信息 getUserInfo（获取好友列表后立即加载、查看消息时加载聊天室内的发言用户）。