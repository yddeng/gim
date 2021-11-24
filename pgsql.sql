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
    "id"           SERIAL,
    "name"         varchar(52)   NOT NULL,
    "age"          int8 NOT NULL DEFAULT '0' ,
    PRIMARY KEY ("id")
) ;