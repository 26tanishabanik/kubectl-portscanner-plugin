# kubectl-portscanner-plugin

This is a simple plugin written in Golang for checking if a port is open and listening in a particular pod.

## Usage

`kubectl portScanner "pod name" "port number" "namespace"`

## Examples

### Closed Port

![alt text](https://github.com/26tanishabanik/kubectl-portscanner-plugin/blob/main/assets/kubectl-port-Closed-cropped.png?raw=true)

### Open Port

![alt text](https://github.com/26tanishabanik/kubectl-portscanner-plugin/blob/main/assets/kubectl-port-Open-cropped.png?raw=true)