package installer

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func logInstallHeader(name string) {
	//logger.PrintHeadline(logger.InstallerHeadline, name)
}

func logInstallPackage(cask string) {
	//msg := fmt.Sprintf("Installing Homebrew package %v (Please be patient, might take some time)", cask)
	//logger.PrintCommandHeader(msg)
}

func logInstallCask(cask string) {
	//msg := fmt.Sprintf("Installing Homebrew cask %v (Please be patient, might take some time)", cask)
	//logger.PrintCommandHeader(msg)
}

type baseInstaller struct {
	shellExec shell.Shell
}

//// region HomeBrew Installer
//type brewInstaller struct {
//	baseInstaller
//}
//
//func (b *brewInstaller) verify() error {
//	if _, err := b.shellExec.ExecuteWithOutput("which brew"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (b *brewInstaller) install() error {
//	installCmd := "/usr/bin/ruby -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\""
//	if err := b.shellExec.Execute(installCmd); err != nil {
//		return err
//	}
//	// Homebrew Permissions Denied Issues Solution
//	// https://gist.github.com/irazasyed/7732946
//	if err := b.shellExec.Execute("sudo chown -R $(whoami) $(brew --prefix)/*"); err != nil {
//		return err
//	}
//	if err := b.shellExec.Execute("sudo chown -R $(whoami) ${HOME}/Library/Caches/Homebrew/*"); err != nil {
//		return err
//	}
//	if err := b.shellExec.Execute("sudo chown -R $(whoami) ${HOME}/Library/Caches/Homebrew/.*"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (b *brewInstaller) installCask(cask string) error {
//	if err := b.Check(); err != nil {
//		return err
//	} else {
//		logInstallCask(cask)
//		caskInstallFormat := "brew tap caskroom/cask && brew cask install %v"
//		installCmd := fmt.Sprintf(caskInstallFormat, cask)
//		if err := b.shellExec.Execute(installCmd); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (b *brewInstaller) installPackage(pkg string) error {
//	if err := b.Check(); err != nil {
//		return err
//	} else {
//		logInstallPackage(pkg)
//		pkgInstallFormat := "brew install %v"
//		installCmd := fmt.Sprintf(pkgInstallFormat, pkg)
//		if err := b.shellExec.Execute(installCmd); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (b *brewInstaller) linkPackageFiles(pkg string) error {
//	if err := b.Check(); err != nil {
//		return err
//	} else {
//		linkPkgFilesFormat := "brew link --force %v"
//		installCmd := fmt.Sprintf(linkPkgFilesFormat, pkg)
//		if err := b.shellExec.Execute(installCmd); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (b *brewInstaller) Check() error {
//	if err := b.verify(); err != nil {
//		logInstallHeader("Homebrew")
//		in := input.New()
//		if result, err := in.WaitForInput("Homebrew is not installed, install?"); err != nil || !result {
//			return errors.Errorf("brew is missing, must be installed, cannot proceed.")
//		} else {
//			return b.install()
//		}
//	}
//	return nil
//}
//
//func NewBrewInstaller(shellExec shell.Shell) *brewInstaller {
//	return &brewInstaller{
//		baseInstaller: baseInstaller{
//			shellExec: shellExec,
//		},
//	}
//}
//
//// endregion
//
//// region Docker Installer
//type dockerInstaller struct {
//	baseInstaller
//}
//
//func (d *dockerInstaller) verify() error {
//	if _, err := d.shellExec.ExecuteWithOutput("which docker"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (d *dockerInstaller) install() error {
//	brew := NewBrewInstaller(d.shellExec)
//	if err := brew.installCask("docker"); err != nil {
//		return err
//	}
//	// Force open Docker
//	if err := d.shellExec.Execute("open -a Docker"); err != nil {
//		return err
//	} else {
//		return errors.Errorf("Please configure Docker and wait for it to become available before using anchor again")
//	}
//}
//
//func (d *dockerInstaller) Check() error {
//	if err := d.verify(); err != nil {
//		logInstallHeader("Docker")
//		in := input.New()
//		if result, err := in.WaitForInput("Docker is not installed, install?"); err != nil || !result {
//			return errors.Errorf("docker is missing, must be installed, cannot proceed.")
//		} else {
//			return d.install()
//		}
//	}
//	return nil
//}
//
//func NewDockerInstaller(shellExec shell.Shell) Installer {
//	return &dockerInstaller{
//		baseInstaller: baseInstaller{
//			shellExec: shellExec,
//		},
//	}
//}
//
//// endregion
//
//// region Kind Installer
//type kindInstaller struct {
//	baseInstaller
//}
//
//func (d *kindInstaller) verify() error {
//	if _, err := d.shellExec.ExecuteWithOutput("which kind"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (d *kindInstaller) install() error {
//	brew := NewBrewInstaller(d.shellExec)
//	if err := brew.installPackage("kind"); err != nil {
//		return err
//	}
//	if err := brew.linkPackageFiles("kind"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (d *kindInstaller) Check() error {
//	if err := d.verify(); err != nil {
//		logInstallHeader("Kind")
//		in := input.New()
//		if result, err := in.WaitForInput("Kind is not installed, install?"); err != nil || !result {
//			return errors.Errorf("kind is missing, must be installed, cannot proceed.")
//		} else {
//			return d.install()
//		}
//	}
//	return nil
//}
//
//func NewKindInstaller(shellExec shell.Shell) Installer {
//	return &kindInstaller{
//		baseInstaller: baseInstaller{
//			shellExec: shellExec,
//		},
//	}
//}
//
//// endregion
//
//// region Kubectl Installer
//type kubectlInstaller struct {
//	baseInstaller
//}
//
//func (k *kubectlInstaller) verify() error {
//	if _, err := k.shellExec.ExecuteWithOutput("which kubectl"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (k *kubectlInstaller) install() error {
//	brew := NewBrewInstaller(k.shellExec)
//	if err := brew.installPackage("kubernetes-cli"); err != nil {
//		return err
//	}
//	if err := brew.linkPackageFiles("kubernetes-cli"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (k *kubectlInstaller) Check() error {
//	if err := k.verify(); err != nil {
//		logInstallHeader("kubectl")
//		in := input.New()
//		if result, err := in.WaitForInput("Kubectl is not installed, install?"); err != nil || !result {
//			return errors.Errorf("kubectl is missing, must be installed, cannot proceed.")
//		} else {
//			return k.install()
//		}
//	}
//	return nil
//}
//
//func NewKubectlInstaller(shellExec shell.Shell) Installer {
//	return &kubectlInstaller{
//		baseInstaller: baseInstaller{
//			shellExec: shellExec,
//		},
//	}
//}
//
//// endregion
//
//// region Helm Installer
//type helmInstaller struct {
//	baseInstaller
//}
//
//func (h *helmInstaller) verify() error {
//	if _, err := h.shellExec.ExecuteWithOutput("which helm"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (h *helmInstaller) install() error {
//	brew := NewBrewInstaller(h.shellExec)
//	if err := brew.installCask("kubernetes-helm"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (h *helmInstaller) Check() error {
//	if err := h.verify(); err != nil {
//		logInstallHeader("Helm")
//		in := input.New()
//		if result, err := in.WaitForInput("Helm is not installed, install?"); err != nil || !result {
//			return errors.Errorf("helm is missing, must be installed, cannot proceed.")
//		} else {
//			return h.install()
//		}
//	}
//	return nil
//}
//
//func NewHelmlInstaller(shellExec shell.Shell) Installer {
//	return &helmInstaller{
//		baseInstaller: baseInstaller{
//			shellExec: shellExec,
//		},
//	}
//}
//
//// endregion
//
//// region Env Substituter Installer
//type envsubstInstaller struct {
//	baseInstaller
//}
//
//func (e *envsubstInstaller) verify() error {
//	if _, err := e.shellExec.ExecuteWithOutput("which envsubst"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (e *envsubstInstaller) install() error {
//	brew := NewBrewInstaller(e.shellExec)
//	if err := brew.installPackage("gettext"); err != nil {
//		return err
//	}
//	if err := brew.linkPackageFiles("gettext"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (e *envsubstInstaller) Check() error {
//	if err := e.verify(); err != nil {
//		logInstallHeader("envsubst")
//		logger.Info("envsubst is mandatory for ENV vars substitution on Dockerfiles/Kubernetes manifests")
//		return e.install()
//	}
//	return nil
//}
//
//func NewEnvsubstInstaller(shellExec shell.Shell) Installer {
//	return &envsubstInstaller{
//		baseInstaller: baseInstaller{
//			shellExec: shellExec,
//		},
//	}
//}
//
//// endregion
//
//// region Hostess Installer
//type hostessInstaller struct {
//	baseInstaller
//}
//
//func (e *hostessInstaller) verify() error {
//	if _, err := e.shellExec.ExecuteWithOutput("which hostess"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (e *hostessInstaller) install() error {
//	brew := NewBrewInstaller(e.shellExec)
//	if err := brew.installPackage("hostess"); err != nil {
//		return err
//	}
//	if err := brew.linkPackageFiles("hostess"); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (e *hostessInstaller) Check() error {
//	if err := e.verify(); err != nil {
//		logInstallHeader("hostess")
//		logger.Infof("hostess is mandatory for managing your /etc/hosts to create a %v DNS record", "")
//		return e.install()
//	}
//	return nil
//}
//
//func NewHostessInstaller(shellExec shell.Shell) Installer {
//	return &hostessInstaller{
//		baseInstaller: baseInstaller{
//			shellExec: shellExec,
//		},
//	}
//}
//
//// endregion
