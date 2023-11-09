package iscsi

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type iscsiLUNInfo struct {
	host string
	fqdn string
}

type ISCSIPlugin struct {
	lunInfo []iscsiLUNInfo
}

const CmdNotFound = "Command Not Found"

func ExecuteCommand(command string, args ...string) (string, string) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var errMsg bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errMsg
	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", CmdNotFound
		}
		return "", errMsg.String()
	}
	return out.String(), errMsg.String()
}

func NewISCSIPlugin() ISCSIPlugin {
	var lun []iscsiLUNInfo
	iscsiPlugin := ISCSIPlugin{lun}
	return iscsiPlugin
}

func (plugin *ISCSIPlugin) CheckIscsiSupport() bool {
	//Check if "iscsiadm" is installed
	_, err := ExecuteCommand("iscsiadm")
	if strings.Contains(err, CmdNotFound) {
		return false
	}
	return true
}

// iscsiadm -m discovery -t sendtargets -p <IP | Target>

func (plugin *ISCSIPlugin) DiscoverLUNs(host string) error {

	if len(host) == 0 {
		err := fmt.Errorf("IP or Hostname is expected")
		return err
	}
	out, errMsg := ExecuteCommand("iscsiadm",
		"-m",
		"discovery",
		"-t",
		"sendtargets",
		"-p",
		host)

	if len(out) > 0 {
		lineArray := strings.Split(out, "\n")
		for _, line := range lineArray {
			if len(line) == 0 {
				break
			}
			token := strings.Split(line, ",")
			var lun iscsiLUNInfo
			lun.host = strings.TrimSpace(token[0])
			// Split again to get only fqdn name.
			fqdn := strings.Split(token[1], " ")
			lun.fqdn = strings.TrimSpace(fqdn[1])
			fmt.Println(lun.host)
			fmt.Println(lun.fqdn)
			plugin.lunInfo = append(plugin.lunInfo, lun)
		}
	}
	if len(errMsg) > 0 {
		err := fmt.Errorf("Unable to Discover: %s", errMsg)
		return err
	}

	return nil
}

// iscsiadm -m node -o show  (Shows discovered list)
func (plugin *ISCSIPlugin) ListVolumes() error {

	out, errMsg := ExecuteCommand("iscsiadm",
		"-m",
		"node",
		"-o",
		"show")

	log.Println(out)
	if len(errMsg) > 0 {
		err := fmt.Errorf("Unable to fetch List: %s", errMsg)
		return err
	}
	return nil
}

// Login: iscsiadm -m node --login (login on all discovered nodes.)
// iscsiadm -m node -T <Complete Target Name> -l -p <Group IP>:3260

func (plugin *ISCSIPlugin) LoginTarget(target string, group string) error {
	var out, errMsg string

	if len(target) == 0 {
		out, errMsg = ExecuteCommand("iscsiadm",
			"-m",
			"node",
			"-l")
	} else {
		if len(group) == 0 {
			err := fmt.Errorf("group IP for target is missing!!")
			return err
		}
		out, errMsg = ExecuteCommand("iscsiadm",
			"-m",
			"node",
			"-T",
			target,
			"-l",
			"-p",
			group)

	}

	log.Println(out)
	if len(errMsg) > 0 {
		err := fmt.Errorf("Unable to Login: %s", errMsg)
		return err
	}
	return nil

}

// iscsiadm -m node -u
// iscsiadm -m node -u -T <Complete Target Name> -p <Group IP address>:3260

func (plugin *ISCSIPlugin) LogoutTarget(target string, group string) error {
	var out, errMsg string
	if len(target) == 0 {
		out, errMsg = ExecuteCommand("iscsiadm",
			"-m",
			"node",
			"-u")
	} else {
		out, errMsg = ExecuteCommand("iscsiadm",
			"-m",
			"node",
			"-u",
			"-T",
			target,
			"-p",
			group)
	}

	log.Println(out)
	if len(errMsg) > 0 {
		err := fmt.Errorf("Unable to Logout: %s", errMsg)
		return err
	}

	return nil

}
