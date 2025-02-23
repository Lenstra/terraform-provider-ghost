# Check: ghost_theme.theme
resource "ghost_theme" "theme" {
  name     = "test-theme"
  activate = true
  source   = "cwd/../../../tests/casper.zip"
  hash     = filesha256("cwd/../../../tests/casper.zip")
}
