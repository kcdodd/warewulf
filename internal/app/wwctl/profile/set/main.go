package set

import (
	"fmt"
	"os"
	"strings"

	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/warewulfd"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	var err error

	nodeDB, err := node.New()
	if err != nil {
		return errors.Wrap(err, "Could not open node configuration")
	}

	profiles, err := nodeDB.FindAllProfiles()
	if err != nil {
		return err
	}

	if len(profiles) == 0 {
		return errors.New("No profiles exist.")
	}

	if SetAll {
		wwlog.Warn("This command will modify all profiles!")

	} else if len(args) > 0 {
		profiles = node.FilterByName(profiles, args)

		if len(profiles) == 0 {
			return errors.Errorf("No profiles found for patterns: %v", args)
		}

	} else {
		wwlog.Warn("No profile specified, selecting the 'default' profile")
		profiles = node.FilterByName(profiles, []string{"default"})

		if len(profiles) == 0 {
			return errors.New("No 'default' profile found.")
		}
	}

	for _, p := range profiles {
		wwlog.Verbose("Modifying profile: %s", p.Id.Get())

		if SetComment != "" {
			wwlog.Verbose("Profile: %s, Setting comment to: %s", p.Id.Get(), SetComment)
			p.Comment.Set(SetComment)
		}

		if SetClusterName != "" {
			wwlog.Verbose("Profile: %s, Setting cluster name to: %s", p.Id.Get(), SetClusterName)
			p.ClusterName.Set(SetClusterName)
		}

		if SetContainer != "" {
			wwlog.Verbose("Profile: %s, Setting Container name to: %s", p.Id.Get(), SetContainer)
			p.ContainerName.Set(SetContainer)
		}

		if SetInit != "" {
			wwlog.Verbose("Profile: %s, Setting init command to: %s", p.Id.Get(), SetInit)
			p.Init.Set(SetInit)
		}

		if SetRoot != "" {
			wwlog.Verbose("Profile: %s, Setting root to: %s", p.Id.Get(), SetRoot)
			p.Root.Set(SetRoot)
		}

		if SetAssetKey != "" {
			wwlog.Verbose("Profile: %s, Setting asset key to: %s", p.Id.Get(), SetAssetKey)
			p.AssetKey.Set(SetAssetKey)
		}

		if SetKernelOverride != "" {
			wwlog.Verbose("Profile: %s, Setting Kernel override version to: %s", p.Id.Get(), SetKernelOverride)
			p.Kernel.Override.Set(SetKernelOverride)
		}

		if SetKernelArgs != "" {
			wwlog.Verbose("Profile: %s, Setting Kernel args to: %s", p.Id.Get(), SetKernelArgs)
			p.Kernel.Args.Set(SetKernelArgs)
		}

		if SetIpxe != "" {
			wwlog.Verbose("Profile: %s, Setting iPXE template to: %s", p.Id.Get(), SetIpxe)
			p.Ipxe.Set(SetIpxe)
		}

		if len(SetRuntimeOverlay) != 0 {
			wwlog.Verbose("Profile: %s, Setting runtime overlay to: %s", p.Id.Get(), SetRuntimeOverlay)
			p.RuntimeOverlay.SetSlice(SetRuntimeOverlay)
		}

		if len(SetSystemOverlay) != 0 {
			wwlog.Verbose("Profile: %s, Setting system overlay to: %s", p.Id.Get(), SetSystemOverlay)
			p.SystemOverlay.SetSlice(SetSystemOverlay)
		}

		if SetIpmiNetmask != "" {
			wwlog.Verbose("Profile: %s, Setting IPMI netmask to: %s", p.Id.Get(), SetIpmiNetmask)
			p.Ipmi.Netmask.Set(SetIpmiNetmask)
		}

		if SetIpmiPort != "" {
			wwlog.Verbose("Profile: %s, Setting IPMI port to: %s", p.Id.Get(), SetIpmiPort)
			p.Ipmi.Port.Set(SetIpmiPort)
		}

		if SetIpmiGateway != "" {
			wwlog.Verbose("Profile: %s, Setting IPMI gateway to: %s", p.Id.Get(), SetIpmiGateway)
			p.Ipmi.Gateway.Set(SetIpmiGateway)
		}

		if SetIpmiUsername != "" {
			wwlog.Verbose("Profile: %s, Setting IPMI username to: %s", p.Id.Get(), SetIpmiUsername)
			p.Ipmi.UserName.Set(SetIpmiUsername)
		}

		if SetIpmiPassword != "" {
			wwlog.Verbose("Profile: %s, Setting IPMI password to: %s", p.Id.Get(), SetIpmiPassword)
			p.Ipmi.Password.Set(SetIpmiPassword)
		}

		if SetIpmiInterface != "" {
			wwlog.Verbose("Profile: %s, Setting IPMI interface to: %s", p.Id.Get(), SetIpmiInterface)
			p.Ipmi.Interface.Set(SetIpmiInterface)
		}

		if SetIpmiWrite == "yes" || SetNetOnBoot == "y" || SetNetOnBoot == "1" || SetNetOnBoot == "true" {
			wwlog.Verbose("Node: %s, Setting Ipmiwrite to %s", p.Id.Get(), SetIpmiWrite)
			p.Ipmi.Write.SetB(true)
		} else {
			wwlog.Verbose("Node: %s, Setting Ipmiwrite to %s", p.Id.Get(), SetIpmiWrite)
			p.Ipmi.Write.SetB(false)
		}

		if SetDiscoverable {
			wwlog.Verbose("Profile: %s, Setting all nodes to discoverable", p.Id.Get())
			p.Discoverable.SetB(true)
		}

		if SetUndiscoverable {
			wwlog.Verbose("Profile: %s, Setting all nodes to undiscoverable", p.Id.Get())
			p.Discoverable.SetB(false)
		}

		if SetNetName != "" {
			if _, ok := p.NetDevs[SetNetName]; !ok {
				var nd node.NetDevEntry
				nd.Tags = make(map[string]*node.Entry)
				p.NetDevs[SetNetName] = &nd
			}
		}

		if SetNetDev != "" {
			if SetNetName == "" {
				return errors.New("You must include the '--netname' option")
			}

			wwlog.Verbose("Node: %s:%s, Setting net Device to: %s", p.Id.Get(), SetNetName, SetNetDev)
			p.NetDevs[SetNetName].Device.Set(SetNetDev)
		}

		if SetNetmask != "" {
			if SetNetName == "" {
				return errors.New("You must include the '--netname' option")
			}

			wwlog.Verbose("Profile '%s': Setting netmask to: %s", p.Id.Get(), SetNetName)
			p.NetDevs[SetNetName].Netmask.Set(SetNetmask)
		}

		if SetGateway != "" {
			if SetNetName == "" {
				wwlog.Error("You must include the '--netname' option")
				os.Exit(1)
			}

			wwlog.Verbose("Profile '%s': Setting gateway to: %s", p.Id.Get(), SetNetName)
			p.NetDevs[SetNetName].Gateway.Set(SetGateway)
		}

		if SetType != "" {
			if SetNetName == "" {
				return errors.New("You must include the '--netname' option")
			}

			wwlog.Verbose("Profile '%s': Setting HW address to: %s:%s", p.Id.Get(), SetNetName, SetType)
			p.NetDevs[SetNetName].Type.Set(SetType)
		}

		if SetNetOnBoot != "" {
			if SetNetName == "" {
				return errors.New("You must include the '--netname' option")
			}

			if SetNetOnBoot == "yes" || SetNetOnBoot == "y" || SetNetOnBoot == "1" || SetNetOnBoot == "true" {
				wwlog.Verbose("Profile: %s:%s, Setting ONBOOT", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].OnBoot.SetB(true)
			} else {
				wwlog.Verbose("Profile: %s:%s, Unsetting ONBOOT", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].OnBoot.SetB(false)
			}
		}

		if SetNetPrimary != "" {
			if SetNetName == "" {
				return errors.New("You must include the '--netname' option")
			}

			if SetNetPrimary == "yes" || SetNetPrimary == "y" || SetNetPrimary == "1" || SetNetPrimary == "true" {

				// Set all other networks to non-default
				for _, n := range p.NetDevs {
					n.Primary.SetB(false)
				}

				wwlog.Verbose("Profile: %s:%s, Setting PRIMARY", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].Primary.SetB(true)
			} else {
				wwlog.Verbose("Profile: %s:%s, Unsetting PRIMARY", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].Primary.SetB(false)
			}
		}

		if SetNetDevDel {
			if SetNetName == "" {
				return errors.New("You must include the '--netname' option")
			}

			if _, ok := p.NetDevs[SetNetName]; !ok {
				return errors.Errorf("Profile '%s': network name doesn't exist: %s", p.Id.Get(), SetNetName)
			}

			wwlog.Verbose("Profile %s: Deleting network: %s", p.Id.Get(), SetNetName)
			delete(p.NetDevs, SetNetName)
		}

		if len(SetTags) > 0 {
			for _, t := range SetTags {
				keyval := strings.SplitN(t, "=", 2)
				key := keyval[0]
				val := keyval[1]

				if _, ok := p.Tags[key]; !ok {
					var nd node.Entry
					p.Tags[key] = &nd
				}

				wwlog.Verbose("Profile: %s, Setting Tag '%s'='%s'", p.Id.Get(), key, val)
				p.Tags[key].Set(val)
			}
		}
		if len(SetDelTags) > 0 {
			for _, t := range SetDelTags {
				keyval := strings.SplitN(t, "=", 1)
				key := keyval[0]

				if _, ok := p.Tags[key]; !ok {
					return errors.Errorf("Profile: %s, tag does not exist: %s", p.Id.Get(), key)
				}

				wwlog.Verbose("Profile: %s, Deleting tag: %s", p.Id.Get(), key)
				delete(p.Tags, key)
			}
		}
		if len(SetNetTags) > 0 {
			for _, t := range SetNetTags {
				keyval := strings.SplitN(t, "=", 2)
				key := keyval[0]
				val := keyval[1]
				if _, ok := p.NetDevs[SetNetName].Tags[key]; !ok {
					var nd node.Entry
					p.NetDevs[SetNetName].Tags[key] = &nd
				}

				wwlog.Verbose("Profile: %s:%s, Setting NETTAG '%s'='%s'", p.Id.Get(), SetNetName, key, val)
				p.NetDevs[SetNetName].Tags[key].Set(val)
			}

		}
		if len(SetNetDelTags) > 0 {
			for _, t := range SetNetDelTags {
				keyval := strings.SplitN(t, "=", 1)
				key := keyval[0]
				if _, ok := p.NetDevs[SetNetName].Tags[key]; !ok {
					return errors.Errorf("Profile: %s, %s, NETTAG does not exist: %s", p.Id.Get(), SetNetName, key)
				}

				wwlog.Verbose("Profile: %s,%s Deleting NETTAG: %s", p.Id.Get(), SetNetName, key)
				delete(p.NetDevs[SetNetName].Tags, key)
			}
		}

		err := nodeDB.ProfileUpdate(p)
		if err != nil {
			return err
		}
	}


	if SetYes {
		err := nodeDB.Persist()
		if err != nil {
			return errors.Wrap(err, "failed to persist nodedb")
		}

		err = warewulfd.DaemonReload()
		if err != nil {
			// NOTE: don't fail if the daemon is not running.
			wwlog.WarnExc(err, "failed to reload warewulf daemon")
		}
	} else {
		q := fmt.Sprintf("Are you sure you want to modify %d profile(s)", len(profiles))

		prompt := promptui.Prompt{
			Label:     q,
			IsConfirm: true,
		}

		result, _ := prompt.Run()

		if result == "y" || result == "yes" {
			err := nodeDB.Persist()
			if err != nil {
				return errors.Wrap(err, "failed to persist nodedb")
			}

			err = warewulfd.DaemonReload()
			if err != nil {
				// NOTE: don't fail if the daemon is not running.
				wwlog.WarnExc(err, "failed to reload warewulf daemon")
			}
		}else{
			wwlog.Warn("Profile update cancelled by user.")
		}
	}

	return nil
}
