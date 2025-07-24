CREATE TABLE IF NOT EXISTS "url_token" (
  id              SERIAL PRIMARY KEY,
  token           VARCHAR UNIQUE NOT NULL, 
  purpose         VARCHAR NOT NULL,        
  user_id         int,                   
  invited_email   VARCHAR,    -- for sending mail to none website users            
  team_id         int,                      
  sender_id       int,                     
  expires_at      TIMESTAMP WITH TIME ZONE,
  created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL,
  FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE,
  FOREIGN KEY ("sender_id") REFERENCES "users"("id") ON DELETE SET NULL
);
