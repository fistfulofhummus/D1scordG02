DIRECTORY=bin
MAC=macos-agent
LINUX=linux-agent
WIN=windows-agent.exe
WINOBF=windows-agentObf.exe
RASP=rasp
BSD=bsd-agent
FLAGS=-ldflags "-s -w"
WIN-FLAGS=-ldflags -H=windowsgui

all: clean create-directory agent-mac agent-linux agent-windows agent-rasp agent-fuckbsd agent-WindowsObf

create-directory:
	mkdir ${DIRECTORY}

agent-mac:
	echo "Compiling macos binary"
	env GOOS=darwin GOARCH=amd64 go build ${FLAGS} -o ${DIRECTORY}/${MAC} cmd/agent/main.go

agent-linux:
	echo "Compiling Linux binary"
	env GOOS=linux GOARCH=amd64 go build -o ${DIRECTORY}/${LINUX} cmd/agent/main.go

agent-windows:
	echo "Compiling Windows binary"
	env GOOS=windows GOARCH=amd64 go build -o ${DIRECTORY}/${WIN} cmd/agentWinDev/main.go

agent-WindowsObf:
	echo "Compiling Obfuscated Windows binary with Garble"
	env CGO_ENABLE=1 GOOS=windows GOARCH=amd64 ~/go/bin/garble -literals -tiny build -trimpath -o ${DIRECTORY}/${WINOBF} cmd/agentWinDev/main.go

agent-rasp:
	echo "Compiling RASPI binary"
	env GOOS=linux GOARCH=arm GOARM=7 go build ${FLAGS} -o ${DIRECTORY}/${RASP} cmd/agent/main.go

agent-fuckbsd:
	echo "Compiling FUCKBSD binary"
	env GOOS=freebsd GOARCH=amd64 go build ${FLAGS} -o ${DIRECTORY}/${BSD} cmd/agent/main.go

clean:
	rm -rf ${DIRECTORY}
