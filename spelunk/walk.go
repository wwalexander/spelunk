package spelunk

import (
	"path"

	vault "github.com/hashicorp/vault/api"
)

type WalkFunc = func(name string, data map[string]interface{}, err error) error

func Walk(client *vault.Logical, root string, fn WalkFunc) error {
	secret, err := client.List(root)
	if err != nil {
		return fn(root, nil, err)
	}
	if secret == nil {
		secret, err := client.Read(root)
		if err != nil {
			return fn(root, nil, err)
		}
		return fn(root, secret.Data, nil)
	}
	if err := fn(root, nil, nil); err != nil {
		return err
	}
	keys := secret.Data["keys"]
	for _, key := range keys.([]interface{}) {
		name := path.Join(root, key.(string))
		if err := Walk(client, name, fn); err != nil {
			return err
		}
	}
	return nil
}
