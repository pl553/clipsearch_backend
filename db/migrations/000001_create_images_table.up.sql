CREATE TABLE IF NOT EXISTS Images(
   ImageID serial PRIMARY KEY,
   SourceUrl TEXT,
   ThumbnailUrl TEXT,
   Sha256 CHAR(64)
);
