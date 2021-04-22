resource "idefix_ci" "example" {
  name        = "myci"
  company_id  = 1234
  project_ids = [1, 2]
}
