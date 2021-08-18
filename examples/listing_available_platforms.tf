/*
 * This simple example will demonstrate how to list all available platforms.
 */
provider "krok" {
}

// Load in all the available datasets
data "krok_platforms" "platforms" {}

// Output the data after it has been synced.
output "platforms" {
  value = data.krok_platforms.platforms
}
