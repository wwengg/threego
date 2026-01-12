// @Title
// @Description
// @Author  Wangwengang  2023/12/24 23:03
// @Update  Wangwengang  2023/12/24 23:03
package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/wwengg/threego/tool/threegoctl/tpl"
)

// Project contains name, license and paths to projects.
type Project struct {
	// v2
	PkgName      string
	Copyright    string
	AbsolutePath string
	Viper        bool
	AppName      string
}

type Command struct {
	CmdName   string
	CmdParent string
	*Project
}

func (c *Command) Create() error {
	if _, err := os.Stat(fmt.Sprintf("%s/proto/pb%s", c.AbsolutePath, c.CmdName)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/proto/pb%s", c.AbsolutePath, c.CmdName), 0751))
	}
	protoFile, err := os.Create(fmt.Sprintf("%s/proto/pb%s/pb%s.proto", c.AbsolutePath, c.CmdName, c.CmdName))
	if err != nil {
		return err
	}
	defer protoFile.Close()

	protoTemplate := template.Must(template.New("proto").Parse(string(tpl.NewProtoTemplate())))
	err = protoTemplate.Execute(protoFile, c)
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}

	// create {{Appname}}.yaml
	configFile, err := os.Create(fmt.Sprintf("%s/%s.yaml", p.AbsolutePath, p.AppName))
	if err != nil {
		return err
	}
	defer configFile.Close()

	configTemplate := template.Must(template.New("configYaml").Parse(string(tpl.ConfigYamlTemplate(p.AppName))))
	err = configTemplate.Execute(configFile, p)
	if err != nil {
		return err
	}

	// create cmd/root.go
	if _, err = os.Stat(fmt.Sprintf("%s/cmd", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/cmd", p.AbsolutePath), 0751))
	}
	rootFile, err := os.Create(fmt.Sprintf("%s/cmd/root.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer rootFile.Close()

	rootTemplate := template.Must(template.New("root").Parse(string(tpl.RootTemplate())))
	err = rootTemplate.Execute(rootFile, p)
	if err != nil {
		return err
	}

	// create global/global.go
	if _, err = os.Stat(fmt.Sprintf("%s/global", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/global", p.AbsolutePath), 0751))
	}
	globalFile, err := os.Create(fmt.Sprintf("%s/global/global.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer globalFile.Close()

	globalTemplate := template.Must(template.New("global").Parse(string(tpl.GlobalTemplate())))
	err = globalTemplate.Execute(globalFile, p)
	if err != nil {
		return err
	}

	// create global/config.go
	if _, err = os.Stat(fmt.Sprintf("%s/global", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/global", p.AbsolutePath), 0751))
	}
	globalConfigFile, err := os.Create(fmt.Sprintf("%s/global/config.go", p.AbsolutePath))
	if err != nil {
		return err

	}
	defer globalConfigFile.Close()

	globalConfigTemplate := template.Must(template.New("globalConfig").Parse(string(tpl.GlobalConfigTemplate())))
	err = globalConfigTemplate.Execute(globalConfigFile, p)
	if err != nil {
		return err
	}

	// create model
	if _, err = os.Stat(fmt.Sprintf("%s/model", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/model", p.AbsolutePath), 0751))
	}

	// create proto
	if _, err = os.Stat(fmt.Sprintf("%s/proto", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/proto", p.AbsolutePath), 0751))
	}

	// create proto/pbcommon
	if _, err = os.Stat(fmt.Sprintf("%s/proto/pbcommon", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/proto/pbcommon", p.AbsolutePath), 0751))
	}
	pbcommonFile, err := os.Create(fmt.Sprintf("%s/proto/pbcommon/pbcommon.proto", p.AbsolutePath))
	if err != nil {
		return err

	}
	defer pbcommonFile.Close()

	commonProtoTemplate := template.Must(template.New("globalConfig").Parse(string(tpl.CommonProtoTemplate())))
	err = commonProtoTemplate.Execute(pbcommonFile, p)
	if err != nil {
		return err
	}

	// create service
	if _, err = os.Stat(fmt.Sprintf("%s/service", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/service", p.AbsolutePath), 0751))
	}

	// create service/impl
	if _, err = os.Stat(fmt.Sprintf("%s/service/impl", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/service/impl", p.AbsolutePath), 0751))
	}

	// create license
	return nil
}
