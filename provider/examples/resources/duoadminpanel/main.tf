resource "sccfm_duo_admin_panel" "panel" {
  name            = var.duo_admin_panel_name
  host            = var.duo_admin_panel_host
  integration_key = var.duo_admin_panel_integration_key
  secret_key      = var.duo_admin_panel_secret_key
  labels          = var.duo_admin_panel_labels
}