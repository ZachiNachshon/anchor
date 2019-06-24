package shell

import (
	"fmt"

	"github.com/kit/pkg/utils/input"
	"github.com/pkg/errors"
)

type baseInstaller struct {
	shellExec Shell
}

// region HomeBrew Installer
type brewInstaller struct {
	baseInstaller
}

func (d *brewInstaller) verify() error {
	if _, err := d.shellExec.ExecuteWithOutput("which brew"); err != nil {
		return err
	}
	return nil
}

func (d *brewInstaller) install() error {
	installCmd := "/usr/bin/ruby -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\""
	if err := d.shellExec.Execute(installCmd); err != nil {
		return err
	}
	return nil
}

func (d *brewInstaller) installCask(cask string) error {
	caskInstallFormat := "brew update && brew tap caskroom/cask && brew search %v && brew cask info %v && brew cask install %v && brew cleanup"
	installCmd := fmt.Sprintf(caskInstallFormat, cask, cask, cask)
	if err := d.shellExec.Execute(installCmd); err != nil {
		return err
	}
	return nil
}

func (d *brewInstaller) installPackage(pkg string) error {
	pkgInstallFormat := "brew update && brew search %v && brew install %v && brew cleanup"
	installCmd := fmt.Sprintf(pkgInstallFormat, pkg, pkg)
	if err := d.shellExec.Execute(installCmd); err != nil {
		return err
	}
	return nil
}

func (d *brewInstaller) Check() error {
	if err := d.verify(); err != nil {
		in := input.NewYesNoInput()
		if result, err := in.WaitForInput("Homebrew is not installed, install?"); err != nil || !result {
			return errors.Errorf("brew is missing, must be installed, cannot proceed.")
		} else {
			return d.install()
		}
	}
	return nil
}

func NewBrewInstaller() *brewInstaller {
	return &brewInstaller{
		baseInstaller: baseInstaller{
			shellExec: NewShellExecutor(BASH),
		},
	}
}

// endregion

// region Docker Installer
type dockerInstaller struct {
	baseInstaller
}

func (d *dockerInstaller) verify() error {
	if _, err := d.shellExec.ExecuteWithOutput("which docker"); err != nil {
		return err
	}
	return nil
}

func (d *dockerInstaller) install() error {
	brew := NewBrewInstaller()
	if err := brew.installCask("docker"); err != nil {
		return err
	}
	return nil
}

func (d *dockerInstaller) Check() error {
	if err := d.verify(); err != nil {
		in := input.NewYesNoInput()
		if result, err := in.WaitForInput("Docker is not installed, install?"); err != nil || !result {
			return errors.Errorf("docker is missing, must be installed, cannot proceed.")
		} else {
			return d.install()
		}
	}
	return nil
}

func NewDockerInstaller() Installer {
	return &dockerInstaller{
		baseInstaller: baseInstaller{
			shellExec: NewShellExecutor(BASH),
		},
	}
}

// endregion

// region Kind Installer
type kindInstaller struct {
	baseInstaller
}

func (d *kindInstaller) verify() error {
	if _, err := d.shellExec.ExecuteWithOutput("which kind"); err != nil {
		return err
	}
	return nil
}

func (d *kindInstaller) install() error {
	if err := d.shellExec.Execute("install kind"); err != nil {
		return err
	}
	return nil
}

func (d *kindInstaller) Check() error {
	if err := d.verify(); err != nil {
		in := input.NewYesNoInput()
		if result, err := in.WaitForInput("Kind is not installed, install?"); err != nil || !result {
			return errors.Errorf("kind is missing, must be installed, cannot proceed.")
		} else {
			return d.install()
		}
	}
	return nil
}

func NewKindInstaller() Installer {
	return &kindInstaller{
		baseInstaller: baseInstaller{
			shellExec: NewShellExecutor(BASH),
		},
	}
}

// endregion
