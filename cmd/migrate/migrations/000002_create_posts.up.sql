CREATE TABLE IF NOT EXISTS posts(  
    "id" bigserial PRIMARY KEY,
    "title" TEXT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "content" TEXT NOT NULL,
    "tags" VARCHAR(100) [],
    "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_user FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE
);