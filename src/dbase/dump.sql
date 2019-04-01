DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "thread";
DROP TABLE IF EXISTS "vote";
DROP TABLE IF EXISTS "post";

CREATE TABLE "user" (
    "uid" SERIAL PRIMARY KEY,
    "nickname" CITEXT UNIQUE,
    "fullname" VARCHAR(100) NOT NULL,
    "about" TEXT,
    "email" CITEXT UNIQUE
);

CREATE TABLE "forum" (
    "fid" SERIAL PRIMARY KEY,
    "title" VARCHAR(100) NOT NULL,
    "author" CITEXT REFERENCES "user"("nickname") ON DELETE CASCADE,
    "slug" CITEXT NOT NULL UNIQUE,
    "posts" INT DEFAULT '0',
    "threads" INT DEFAULT '0'
);

CREATE TABLE "thread" (
    "tid" BIGSERIAL PRIMARY KEY,
    "title" VARCHAR(100) NOT NULL,
    "author" CITEXT REFERENCES "user"("nickname") ON DELETE CASCADE,
    "message" TEXT NOT NULL,
    "forumSlug" CITEXT REFERENCES "forum"("slug") ON DELETE CASCADE,
    "votes" INT DEFAULT '0',
    "slug" CITEXT UNIQUE DEFAULT NULL,
    "created" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE "vote" (
    "vid" SERIAL PRIMARY KEY,
    "author" CITEXT REFERENCES "user"("nickname") ON DELETE CASCADE,
    "tid" INT8 REFERENCES "thread"("tid") ON DELETE CASCADE,
    "voice" INT2 NOT NULL
);

CREATE TABLE "post" (
    "pid" BIGSERIAL PRIMARY KEY,
    "parent" INT8 DEFAULT '0',
    "path" BIGINT[],
    "author" CITEXT REFERENCES "user"("nickname") ON DELETE CASCADE,
    "message" TEXT NOT NULL,
    "isEdited" BOOLEAN DEFAULT 'false',
    "forumSlug" CITEXT REFERENCES "forum"("slug") ON DELETE CASCADE,
    "tid" BIGINT REFERENCES "thread"("tid") ON DELETE CASCADE,
    "created" TIMESTAMP WITH TIME ZONE DEFAULT '1970-01-01T00:00:00.000Z'
);

