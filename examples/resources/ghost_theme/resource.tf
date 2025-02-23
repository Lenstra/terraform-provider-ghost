resource "ghost_theme" "theme" {
  name     = "my-theme"
  activate = true
  source   = "my-theme.zip"
  hash     = filesha256("my-theme.zip")
}
