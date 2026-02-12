output "admin_session_token" {
  value     = rancher2_login.initial_admin.session_token
  sensitive = true
}
output "admin_user_token" {
  value     = rancher2_login.initial_admin.user_token
  sensitive = true
}
output "admin_user_token_start_date" {
  value = rancher2_login.initial_admin.user_token_start_date
}
output "admin_user_token_end_date" {
  value = rancher2_login.initial_admin.user_token_end_date
}
output "admin_refresh_date" {
  value = rancher2_login.initial_admin.user_token_refresh_date
}

# output "kate_session_token" {
#   value = rancher2_login.kate.session_token
# }
# output "kate_user_token" {
#   value = rancher2_login.kate.user_token
# }
# output "kate_user_token_start_date" {
#   value = rancher2_login.kate.user_token_start_date
# }
# output "kate_user_token_end_date" {
#   value = rancher2_login.kate.user_token_end_date
# }
# output "kate_refresh_date" {
#   value = rancher2_login.kate.user_token_refresh_date
# }
