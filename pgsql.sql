DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
"id" varchar(255) NOT NULL,
"create_at" int8 NOT NULL,
"update_at" int8 NOT NULL,
"attr" bytea NOT NULL,
"convs" bytea NOT NULL,
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
    "members"          bytea NOT NULL,
    PRIMARY KEY ("id")
);