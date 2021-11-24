
DROP TABLE IF EXISTS "conversation_list";
CREATE TABLE "conversation_list" (
    "id"           SERIAL,
    "name"         varchar(52)   NOT NULL,
    "age"          int8 NOT NULL DEFAULT '0' ,
    PRIMARY KEY ("id")
) ;