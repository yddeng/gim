
-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
    "id"           SERIAL,
    "name"         varchar(52)   NOT NULL,
    "age"          int8 NOT NULL DEFAULT '0' ,
    PRIMARY KEY ("id")
) ;
