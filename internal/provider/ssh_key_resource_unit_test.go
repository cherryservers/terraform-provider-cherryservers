package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// newSSHKeyModel runs f over a model with zero valued fields set to null,
// as required for a valid model, and returns the model. Nil f is ignored.
func newSSHKeyModel(f func(m *sshKeyResourceModel)) sshKeyResourceModel {
	m := sshKeyResourceModel{
		Name:        types.StringNull(),
		PublicKey:   types.StringNull(),
		Fingerprint: types.StringNull(),
		Created:     types.StringNull(),
		Updated:     types.StringNull(),
		ID:          types.StringNull(),
	}

	if f != nil {
		f(&m)
	}
	return m
}
