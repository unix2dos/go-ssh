# go-ssh

This is a simple ssh client

#### usage:

```
func main() {
   client := ssh.NewSSHClient("172.24.120.46", 22, "root", "password", "")
   res, err := client.RunCommand("ls -al")
   if err != nil {
      fmt.Println("ssh error: %v", err)
      return
   }
   fmt.Println(res)
}
```

