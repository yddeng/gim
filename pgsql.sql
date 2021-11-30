DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
    "id"         varchar(255) NOT NULL,
    "create_at"  int8 NOT NULL,
    "update_at"  int8 NOT NULL,
    "extra"      bytea NOT NULL, /* 附加属性 */
PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "groups";
CREATE TABLE "groups" (
    "id"               SERIAL8,
    "type"             int4 NOT NULL ,
    "creator"          varchar(255) NOT NULL ,
    "create_at"        int8 NOT NULL,
    "extra"            bytea NOT NULL, /* 附加属性 */
    "last_message_id"  int8 NOT NULL DEFAULT 0,
    "last_message_at"  int8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "group_member";
CREATE TABLE "group_member" (
    "id"               varchar(255) NOT NULL,     /* 组合建 group_id_user_id */
    "group_id"          int8 NOT NULL,             /* 会话ID */
    "user_id"          varchar(255) NOT NULL ,    /* 用户ID */
    "nickname"         varchar(255) ,             /* 群内昵称 */
    "create_at"        int8 NOT NULL,             /* 入群时间,单位：秒 */
    "update_at"        int8 NOT NULL,             /* 上一次活跃时间,单位：秒 */
    "mute"             int4 NOT NULL DEFAULT 0 ,  /* 群内禁言，1:禁言 */
    "role"             int4 NOT NULL DEFAULT 0,   /* 成员角色 0:普通成员 1:管理员 */
    PRIMARY KEY ("id") /* 组合建 */
);

DROP TABLE IF EXISTS "message_list";
CREATE TABLE "message_list" (
    "table_name"    varchar(255) NOT NULL,
    PRIMARY KEY ("table_name")
);

DROP TABLE IF EXISTS "friend";
CREATE TABLE "friend" (
    "id"               varchar(255) NOT NULL,     /* 组合建，user1_user2 */
    "user1_id"         varchar(255) NOT NULL ,    /* 用户1ID, code较小 */
    "user2_id"         varchar(255) ,             /* 用户2ID, code较大*/
    "create_at"        int8 NOT NULL,             /* 创建时间,单位：秒 */
    "status"           int4 NOT NULL,             /* 1:u1->u2 2:u2->u1 3:both 4:friend */
    PRIMARY KEY ("id") /* 组合建 */
);

/*
DROP TABLE IF EXISTS "message_2021125";
CREATE TABLE "message_2021125" (
    "id"           varchar(255) NOT NULL, /* 组合建 group_id_message_id */
    "group_id"      int8 NOT NULL ,
    "message_id"   int8 NOT NULL,
    "message"      bytea NOT NULL,
    PRIMARY KEY ("id")
);
 */

