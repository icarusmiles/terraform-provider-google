// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iam3_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIAM3FoldersPolicyBinding_iamFoldersPolicyBindingExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckIAM3FoldersPolicyBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAM3FoldersPolicyBinding_iamFoldersPolicyBindingExample_full(context),
			},
			{
				ResourceName:            "google_iam_folders_policy_binding.my-folder-binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "folder", "location", "policy_binding_id"},
			},
			{
				Config: testAccIAM3FoldersPolicyBinding_iamFoldersPolicyBindingExample_update(context),
			},
			{
				ResourceName:            "google_iam_folders_policy_binding.my-folder-binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "folder", "location", "policy_binding_id"},
			},
		},
	})
}

func testAccIAM3FoldersPolicyBinding_iamFoldersPolicyBindingExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_principal_access_boundary_policy" "pab_policy" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test folder binding%{random_suffix}"
  principal_access_boundary_policy_id = "tf-test-my-pab-policy%{random_suffix}"
}

resource "google_folder" "folder" {
  display_name        = "test folder%{random_suffix}"
  parent              = "organizations/%{org_id}"
  deletion_protection = false
}

resource "time_sleep" "wait_120s" {
  depends_on      = [google_folder.folder]
  create_duration = "120s"
}

resource "google_iam_folders_policy_binding" "my-folder-binding" {
  folder         = google_folder.folder.folder_id
  location       = "global"
  display_name   = "test folder binding%{random_suffix}"
  policy_kind    = "PRINCIPAL_ACCESS_BOUNDARY"
  policy_binding_id = "tf-test-folder-binding%{random_suffix}"
  policy         = "organizations/%{org_id}/locations/global/principalAccessBoundaryPolicies/${google_iam_principal_access_boundary_policy.pab_policy.principal_access_boundary_policy_id}"
  target {
    principal_set = "//cloudresourcemanager.googleapis.com/folders/${google_folder.folder.folder_id}"
  }
  depends_on = [time_sleep.wait_120s]
}
`, context)
}

func testAccIAM3FoldersPolicyBinding_iamFoldersPolicyBindingExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_principal_access_boundary_policy" "pab_policy" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test folder binding%{random_suffix}"
  principal_access_boundary_policy_id = "tf-test-my-pab-policy%{random_suffix}"
}

resource "google_folder" "folder" {
  display_name        = "test folder%{random_suffix}"
  parent              = "organizations/%{org_id}"
  deletion_protection = false
}

resource "time_sleep" "wait_120s" {
  depends_on      = [google_folder.folder]
  create_duration = "120s"
}

resource "google_iam_folders_policy_binding" "my-folder-binding" {
  folder         = google_folder.folder.folder_id
  location       = "global"
  display_name   = "test folder binding%{random_suffix}"
  policy_kind    = "PRINCIPAL_ACCESS_BOUNDARY"
  policy_binding_id = "tf-test-folder-binding%{random_suffix}"
  policy         = "organizations/%{org_id}/locations/global/principalAccessBoundaryPolicies/${google_iam_principal_access_boundary_policy.pab_policy.principal_access_boundary_policy_id}"
  annotations    = {"foo": "bar"}
  target {
    principal_set = "//cloudresourcemanager.googleapis.com/folders/${google_folder.folder.folder_id}"
  }
  condition {
    description   = "test condition"
    expression    = "principal.subject == 'al@a.com'"
    location      = "test location"
    title         = "test title"
  }
  depends_on = [time_sleep.wait_120s]
}
`, context)
}