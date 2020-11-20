package helpers

import (
	"io/ioutil"
	"log"
	"net"
	"os"
)

// LimanUser : Just in case if liman user changed to something else.
const LimanUser = "liman"

// DefaultShell : Default sh shell
const DefaultShell = "/bin/bash"

// ResolvPath : Dns server' configuration path.
const ResolvPath = "/etc/resolv.conf"

// DNSOptions : Options to have multiple dns servers
const DNSOptions = "options rotate timeout:1 retries:1"

func FixExtensionKeys(extensionID string) bool {
	_, err := ExecuteCommand("chmod -R 700 " + KeysPath + extensionID)
	if err != nil {
		return false
	}

	_, err = ExecuteCommand("chown -R " + extensionID + ":" + LimanUser + " " + KeysPath + extensionID)
	if err == nil {
		return true
	}
	return false
}

func AddUser(extensionID string) bool {
	log.Println("Adding System User : " + extensionID)
	_, err := ExecuteCommand("useradd -r -s " + DefaultShell + " " + extensionID)
	if err == nil {
		log.Println("System User Added : " + extensionID)
		return true
	}
	log.Println(err)
	return false
}

func RemoveUser(extensionID string) bool {
	log.Println("Removing System User : " + extensionID)
	_, err := ExecuteCommand("userdel " + extensionID)
	if err == nil {
		log.Println("System User Removed : " + extensionID)
		return true
	}
	log.Println(err)
	return false
}

func FixExtensionPermissions(extensionID string, extensionName string) bool {
	_, err := ExecuteCommand("chmod -R 770 " + ExtensionsPath + extensionName + " 2>&1")
	log.Println("Fixing Extension Permissions")
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = ExecuteCommand("chown -R " + extensionID + ":" + LimanUser + " " + ExtensionsPath + extensionName + " 2>&1")
	if err == nil {
		log.Println("Extension Permissions Fixed")
		return true
	}
	log.Println(err)
	return false
}

func AddSystemCertificate(tmpPath string, targetName string) bool {
	certPath, certUpdateCommand := GetCertificateStrings()
	log.Println("Adding System Certificate")
	_, err := ExecuteCommand("mv " + tmpPath + " " + certPath + "/" + targetName + ".crt")
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = ExecuteCommand(certUpdateCommand)
	if err == nil {
		log.Println("System Certificate Added")
		return true
	}
	log.Println(err)
	return false
}

func RemoveSystemCertificate(targetName string) bool {
	log.Println("Removing System Certificate")
	certPath, certUpdateCommand := GetCertificateStrings()
	_, err := ExecuteCommand("rm " + certPath + "/" + targetName + ".crt")
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = ExecuteCommand(certUpdateCommand)
	if err == nil {
		log.Println("System Certificate Removed")
		return true
	}
	log.Println(err)
	return false
}

func GetCertificateStrings() (string, string) {
	certPath := "/usr/local/share/ca-certificates/"
	certUpdateCommand := "update-ca-certificates"
	if IsCentOs() == true {
		certPath = "/etc/pki/ca-trust/source/anchors/"
		certUpdateCommand = "sudo update-ca-trust"
	}
	return certPath, certUpdateCommand
}

func SetDNSServers(server1 string, server2 string, server3 string) bool {
	_, err := ExecuteCommand("chattr -i " + ResolvPath)
	log.Println("Updating DNS Servers")
	if err != nil {
		log.Println(err)
		return false
	}
	newData := DNSOptions + "\n"
	if server1 != "" {
		newData += "nameserver " + server1 + "\n"
	}

	if server2 != "" {
		newData += "nameserver " + server2 + "\n"
	}

	if server3 != "" {
		newData += "nameserver " + server3 + "\n"
	}

	err = ioutil.WriteFile(ResolvPath, []byte(newData), 0644)

	if err != nil {
		log.Println(err)
		return false
	}

	_, err = ExecuteCommand("chattr +i " + ResolvPath)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("DNS Servers Updated")
	return true
}

func IsCentOs() bool {
	_, err := os.Stat("/etc/redhat-release")
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetLocalIP() string {
	if CurrentIP != "" {
		return CurrentIP
	}

	conn, err := net.Dial("tcp", DBHost+":"+DBPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr)

	return localAddr.IP.String()
}
