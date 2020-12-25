CREATE TABLE "users" (
    "id"            INTEGER,
    "username"      VARCHAR(50) NOT NULL UNIQUE,
    "email"         VARCHAR(255) NOT NULL UNIQUE,
    "password"      VARCHAR(255) NOT NULL,
    "last_login"    DATETIME,
    "created_at"    DATETIME DEFAULT CURRENT_TIMESTAMP,
    "modified_at"   DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE "pages" (
    "id"            INTEGER,
    "parent_id"     INTEGER,
    "title"         VARCHAR(255) NOT NULL,
    "content"       TEXT,
    "is_blog"       BOOLEAN DEFAULT 'FALSE',
    "created_at"    DATETIME DEFAULT CURRENT_TIMESTAMP,
    "modified_at"   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY("parent_id") REFERENCES "pages"("id"),
    PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE "posts" (
    "id"            INTEGER,
    "blog_id"       INTEGER NOT NULL,
    "title"         VARCHAR(255) NOT NULL,
    "content"       TEXT NOT NULL,
    "created_at"    DATETIME DEFAULT CURRENT_TIMESTAMP,
    "modified_at"   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY("blog_id") REFERENCES "pages"("id"),
    PRIMARY KEY("id" AUTOINCREMENT)
);
