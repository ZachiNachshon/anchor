package installer

import (
	"fmt"
	"github.com/anchor/pkg/logger"
	"github.com/anchor/pkg/utils/shell"

	"github.com/anchor/pkg/utils/input"
	"github.com/pkg/errors"
)

type baseInstaller struct {
	shellExec shell.Shell
}

// region HomeBrew Installer
type brewInstaller struct {
	baseInstaller
}

func (b *brewInstaller) verify() error {
	if _, err := b.shellExec.ExecuteWithOutput("which brew"); err != nil {
		return err
	}
	return nil
}

func (b *brewInstaller) install() error {
	logger.Info("==> Installing Homebrew...")
	installCmd := "/usr/bin/ruby -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\""
	if err := b.shellExec.Execute(installCmd); err != nil {
		return err
	}
	return nil
}

func (b *brewInstaller) installCask(cask string) error {
	if err := b.Check(); err != nil {
		return err
	} else {
		logger.Infof("  ==> Installing Homebrew cask %v...", cask)
		caskInstallFormat := "brew update && brew tap caskroom/cask && brew search %v && brew cask info %v && brew cask install %v && brew cleanup"
		installCmd := fmt.Sprintf(caskInstallFormat, cask, cask, cask)
		if err := b.shellExec.Execute(installCmd); err != nil {
			return err
		}
	}
	return nil
}

func (b *brewInstaller) installPackage(pkg string) error {
	if err := b.Check(); err != nil {
		return err
	} else {
		logger.Infof("  ==> Installing Homebrew package %v...", pkg)
		pkgInstallFormat := "brew update && brew search %v && brew install %v && brew cleanup"
		installCmd := fmt.Sprintf(pkgInstallFormat, pkg, pkg)
		if err := b.shellExec.Execute(installCmd); err != nil {
			return err
		}
	}
	return nil
}

func (b *brewInstaller) linkPackageFiles(pkg string) error {
	if err := b.Check(); err != nil {
		return err
	} else {
		linkPkgFilesFormat := "brew link --force %v"
		installCmd := fmt.Sprintf(linkPkgFilesFormat, pkg)
		if err := b.shellExec.Execute(installCmd); err != nil {
			return err
		}
	}
	return nil
}

func (b *brewInstaller) Check() error {
	if err := b.verify(); err != nil {
		in := input.NewYesNoInput()
		if result, err := in.WaitForInput("Homebrew is not installed, install?"); err != nil || !result {
			return errors.Errorf("brew is missing, must be installed, cannot proceed.")
		} else {
			return b.install()
		}
	}
	return nil
}

func NewBrewInstaller(shellExec shell.Shell) *brewInstaller {
	return &brewInstaller{
		baseInstaller: baseInstaller{
			shellExec: shellExec,
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
	logger.Info("==> Installing Docker...")
	brew := NewBrewInstaller(d.shellExec)
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

func NewDockerInstaller(shellExec shell.Shell) Installer {
	return &dockerInstaller{
		baseInstaller: baseInstaller{
			shellExec: shellExec,
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
	logger.Info("==> Installing Kind...")
	if err := d.shellExec.Execute("export GO111MODULE=\"on\" && go get sigs.k8s.io/kind@v0.4.0"); err != nil {
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

func NewKindInstaller(shellExec shell.Shell) Installer {
	return &kindInstaller{
		baseInstaller: baseInstaller{
			shellExec: shellExec,
		},
	}
}

// endregion

// region Kubectl Installer
type kubectlInstaller struct {
	baseInstaller
}

func (k *kubectlInstaller) verify() error {
	if _, err := k.shellExec.ExecuteWithOutput("which kubectl"); err != nil {
		return err
	}
	return nil
}

func (k *kubectlInstaller) install() error {
	brew := NewBrewInstaller(k.shellExec)
	logger.Info("==> Installing kubectl...")
	if err := brew.installCask("kubernetes-cli"); err != nil {
		return err
	}
	return nil
}

func (k *kubectlInstaller) Check() error {
	if err := k.verify(); err != nil {
		in := input.NewYesNoInput()
		if result, err := in.WaitForInput("Kubectl is not installed, install?"); err != nil || !result {
			return errors.Errorf("kubectl is missing, must be installed, cannot proceed.")
		} else {
			return k.install()
		}
	}
	return nil
}

func NewKubectlInstaller(shellExec shell.Shell) Installer {
	return &kubectlInstaller{
		baseInstaller: baseInstaller{
			shellExec: shellExec,
		},
	}
}

// endregion

// region Helm Installer
type helmInstaller struct {
	baseInstaller
}

func (h *helmInstaller) verify() error {
	if _, err := h.shellExec.ExecuteWithOutput("which helm"); err != nil {
		return err
	}
	return nil
}

func (h *helmInstaller) install() error {
	brew := NewBrewInstaller(h.shellExec)
	logger.Info("==> Installing Helm...")
	if err := brew.installCask("kubernetes-helm"); err != nil {
		return err
	}
	return nil
}

func (h *helmInstaller) Check() error {
	if err := h.verify(); err != nil {
		in := input.NewYesNoInput()
		if result, err := in.WaitForInput("Helm is not installed, install?"); err != nil || !result {
			return errors.Errorf("helm is missing, must be installed, cannot proceed.")
		} else {
			return h.install()
		}
	}
	return nil
}

func NewHelmlInstaller(shellExec shell.Shell) Installer {
	return &helmInstaller{
		baseInstaller: baseInstaller{
			shellExec: shellExec,
		},
	}
}

// endregion

// region Env Substituter Installer
type envsubstInstaller struct {
	baseInstaller
}

func (e *envsubstInstaller) verify() error {
	if _, err := e.shellExec.ExecuteWithOutput("which envsubst"); err != nil {
		return err
	}
	return nil
}

func (e *envsubstInstaller) install() error {
	brew := NewBrewInstaller(e.shellExec)
	logger.Info("==> Installing envsubst...")
	if err := brew.installPackage("gettext"); err != nil {
		return err
	}
	if err := brew.linkPackageFiles("gettext"); err != nil {
		return err
	}
	return nil
}

func (e *envsubstInstaller) Check() error {
	if err := e.verify(); err != nil {
		logger.Info("envsubst is mandatory for ENV vars substitution on Kubernetes manifests, installing...")
		return e.install()
	}
	return nil
}

func NewEnvsubstInstaller(shellExec shell.Shell) Installer {
	return &envsubstInstaller{
		baseInstaller: baseInstaller{
			shellExec: shellExec,
		},
	}
}

// endregion
