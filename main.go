package main

import (
	"strconv"
	"fmt"
	"os"
	"io"
	corev1 "k8s.io/api/core/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/remotecommand"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"bytes"
	"bufio"
	"log"

)

func readLine(reader *bufio.Reader) (strLine string, err error) {
	buffer := new(bytes.Buffer)
	for {
		var line []byte
		var isPrefix bool
		line, isPrefix, err = reader.ReadLine()
		if err != nil && err != io.EOF {
			return "", err
		}
		buffer.Write(line)
		if !isPrefix {
			break
		}
	}
	return buffer.String(), err
}

func PodCommand(podName ,port, namespace string) {
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
	)

	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
			panic(err)
	}

	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
			panic(err)
	}
	cmd := []string{
        "sh",
        "-c",
        "netstat -tulpn | grep LISTEN | grep "+port,
    }
	req := coreclient.RESTClient().
			Post().
			Namespace(namespace).
			Resource("pods").
			Name(podName).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Command: cmd,
				Stdin:   true,
				Stdout:  true,
				Stderr:  true,
				TTY:     true,
			}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
	if err != nil {
			fmt.Println(err)
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
			fmt.Println(err)
	}
	defer terminal.Restore(0, oldState)
	stdout := new(bytes.Buffer)
   	stderr := new(bytes.Buffer)
	err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: stdout,
			Stderr: stderr,
			Tty:    true,
	})
	if err != nil {
			fmt.Println(err)
	}
	
	i := 0
	flag := 0
	reader := bufio.NewReader(stdout)
	for {
		strline, err := readLine(reader)
		if err != nil && err != io.EOF {
			log.Println(err)
		}
		if len(strline) > 0 {
			words:= strings.Fields(strline)
			portShown:= strings.Split(words[3], ":")
			if i >0 {
				givenPort,_ := strconv.ParseInt(port, 10,0)
				shownPort, _ := strconv.ParseInt(portShown[1], 10, 0)
				if givenPort == shownPort {
					flag = 1
					break
				}
			}	
		}
		i += 1
		if err == io.EOF {
			break
		}
	}
	if flag == 1{
		fmt.Printf("%s is open and listening", port)
	}else{
		fmt.Printf("%s is closed", port)
	}
}


func main(){
	if len(os.Args) < 3{
		fmt.Println("First argument is for pod name, second argument is for port number and third argument is for namespace")
	}else {
		PodCommand(os.Args[1],os.Args[2], os.Args[3])
	}
}
