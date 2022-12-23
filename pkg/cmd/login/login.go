package login

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olekukonko/tablewriter"
)

const subdomain string = "clidom.es"

// NewLogin login to the selected instance.
func NewLogin(name, region, user, profile string, silent, ssh, pushKey bool) {

	client, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
		Config:            aws.Config{Region: aws.String(region)},
	})
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	rawInstanceList := filterInstances(region, name, silent, client)
	instances := getInstancesInfo(rawInstanceList)

	if silent {
		showIPsList(instances)
	} else {
		showInstanceList(instances, user)
	}

	selectedInstance := selectInstanceIndex(instances)

	if ssh {
		launchSSH(instances[selectedInstance].IP, user)
	} else if pushKey {
		pushTempKeyPair(instances[selectedInstance].ID, instances[selectedInstance].AZ, instances[selectedInstance].IP, user)
	} else {
		launchSSM(instances[selectedInstance].ID, profile)
	}
}

func showIPsList(instanceList []Instance) {
	if len(instanceList) == 0 {
		os.Exit(0)
	}
	for _, inst := range instanceList {
		fmt.Printf("%s\n", inst.IP)
	}
	os.Exit(0)
}

func showInstanceList(instanceList []Instance, user string) {
	if len(instanceList) == 0 {
		fmt.Printf("There are no instances matching your request.\n")
		os.Exit(0)
	}

	tableData := [][]string{}

	for idx, inst := range instanceList {
		dataInstance := []string{
			strconv.Itoa(idx),
			inst.Name,
			convertNameToDNS(inst.Name),
			inst.IP,
			inst.ID,
			inst.Size,
			inst.LaunchTime.Format("2006-01-02 15:04:05"),
		}
		tableData = append(tableData, dataInstance)
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Id", "Name", "DNS", "IP", "Instance Id", "Size", "LaunchTime"})
	table.SetAutoWrapText(false)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(tableData)
	table.Render()
}

func selectInstanceIndex(instanceList []Instance) int {
	var input string
	var index int
	var err error
	if len(instanceList) == 1 {
		fmt.Printf("\n\n")
		index = 0
	} else {
		fmt.Println("\nWhich one do you want to ssh in?")
		fmt.Scanln(&input)
		index, err = strconv.Atoi(input)
		if err != nil || index > len(instanceList)-1 || index < 0 {
			fmt.Println("Please enter a valid number.")
			index = selectInstanceIndex(instanceList)
		}
	}
	return index
}

func convertNameToDNS(name string) string {
	if !strings.Contains(name, "bastion") {
		return name
	}
	replaceBastion := strings.ReplaceAll(name, "bastion-", "bastion.")
	replaceBastionV2 := strings.ReplaceAll(replaceBastion, "v2-", "")
	replaceProd := strings.ReplaceAll(replaceBastionV2, "-prod", "")
	replaceStaging := strings.ReplaceAll(replaceProd, "-staging", ".staging")
	dnsSplit := strings.Split(replaceStaging, ".")
	if strings.Contains(name, "staging") {
		return dnsSplit[1] + "." + dnsSplit[0] + "." + dnsSplit[2] + "." + subdomain
	}

	return dnsSplit[1] + "." + dnsSplit[0] + "." + subdomain
}
