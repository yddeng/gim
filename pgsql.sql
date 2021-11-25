DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
"id" varchar(255) NOT NULL,
"create_at" int8 NOT NULL,
"update_at" int8 NOT NULL,
"extra"     bytea NOT NULL, /* 附加属性 */
PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "conversation_list";
CREATE TABLE "conversation_list" (
    "id"               SERIAL8,
    "type"             int4 NOT NULL ,
    "name"             varchar(255),
    "creator"          varchar(255) NOT NULL ,
    "create_at"        int8 NOT NULL,
    "last_message_id"  int8 NOT NULL DEFAULT 0,
    "last_message_at"  int8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "conv_user";
CREATE TABLE "conv_user" (
    "id"               varchar(255) NOT NULL, /* 组合建 conversation_id_user_id */
    "conv_id"          int8 NOT NULL,
    "user_id"          varchar(255) NOT NULL ,
    "role"             int4 NOT NULL DEFAULT 0, /* 成员角色 0:普通成员 1:管理员*/
    PRIMARY KEY ("id")
);


/*
DROP TABLE IF EXISTS "message_2021125";
CREATE TABLE "message_2021125" (
    "id"           varchar(255) NOT NULL, /* 组合建 conv_id_message_id */
    "conv_id"      int8 NOT NULL ,
    "message_id"   int8 NOT NULL,
    "message"      bytea NOT NULL,
    PRIMARY KEY ("id")
);
 */

