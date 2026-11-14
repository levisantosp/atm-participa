env "local" {
  src = "ent://ent/schema"
  dev = "docker://postgres/18/dev?search_path=public"
  url = getenv("DATABASE_URL")
  
  migration {
    dir = "file://migrations"
  }
}
