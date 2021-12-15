# GIG Initializer windows service

The GIG Initializer windows service allows you to run powershell scripts to initialize your windows vm. At boot time the GIG Initializer service scans the `C:\gig\init` directory for powershell scripts (`*.ps1`). It runs the scripts one by one. After running them, it creates a `<<script name>>.ps1.executed` file to indicate that this particular init script has been executed. Scripts are never executed twice, even across several reboots. When all scripts are executed, the GIG Initializer service stops itself.

## How to dynamically add scripts to windows VM instances

The GIG cloud makes use of [cloud-init](https://cloudinit.readthedocs.io/en/latest/) for initialization of virtual machines. For windows vms the compatible cloudbase-init software is used. Using the [Userdata](https://cloudbase-init.readthedocs.io/en/latest/userdata.html#cloud-config) field of the machine create api, you can insert scripts in your windows vm. Note that you need to put the in `C:\gig\init` for them to be picked up by the GIG Initializer service.

# Installation

- Download the executable for your platform from the [releases page](https://github.com/gig-tech/windows-init/releases)
- Then register the program as a Windows service by running it with the `install` flag. Eg
```powershell
mkdir "C:\Program Files\GIG.tech"
cd "C:\Program Files\GIG.tech"
mkdir "GIG Init"
cd "GIG Init"
Invoke-WebRequest https://github.com/gig-tech/windows-init/releases/download/v1.0.0-beta2/giginit-x64.exe -OutFile giginit.exe
giginit.exe install
```
