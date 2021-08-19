/*
 * This example will create a command and a repository
 * and create a relationship between that command and
 * the repo.
 */
provider "krok" {
}

/*
 * Create github platform token.
 */
resource "krok_platform" "github" {
  platform_id = 1
  token = "token"
}

/*
 * Create a setting for this command.
 */
resource "krok_command_setting" "channel" {
  key = "channel"
  value = "krok"
  in_vault = false
  command_id = krok_command.slack_notification.id
}

/*
 * Create a Slack notification command.
 */
resource "krok_command" "slack_notification" {
  name = "slack-notification"
  url = "https://github.com/krok-o/plugins/releases/download/v0.1.0/slack-notification.tar.gz"
  platforms = [1]
}

/*
 * Create a repository.
 */
resource "krok_repository" "skarlso_test" {
  name = "skarlso-test"
  url = "https://github.com/Skarlso/test"
  vcs = 1
  auth {
    secret = "secret"
  }
  commands = [krok_command.slack_notification.id]
  events = ["push"]
}
