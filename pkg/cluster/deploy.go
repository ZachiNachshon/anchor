package cluster

import (
	"github.com/ZachiNachshon/anchor/pkg/docker"
)

func Deploy(identifier string, namespace string) error {

	// TODO: Should verify delete before starting again ?

	if err := docker.Build(identifier); err != nil {
		return err
	}

	if err := docker.Tag(identifier); err != nil {
		return err
	}

	// Verify Registry is up and running, else start before try to tag & push
	if err := Registry(false); err != nil {
		return err
	}

	if err := docker.Push(identifier); err != nil {
		return err
	}

	if _, err := Apply(identifier, namespace); err != nil {
		return err
	}

	if err := Expose(identifier); err != nil {
		return err
	}

	return nil
}
