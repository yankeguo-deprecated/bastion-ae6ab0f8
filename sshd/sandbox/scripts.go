/**
 * sandbox/scripts.go
 *
 * Copyright (c) 2018 Yanke Guo <guoyk.cn@gmail.com>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package sandbox

import (
	"bytes"
	"log"
	"text/template"
)

const scriptGenerateSSHKey = `#!/bin/bash
# write README
echo "注意事项:
1. 沙箱环境互相隔离，可以自由使用 root 权限
2. 系统自动将 id_rsa.pub 公钥文件同步到数据库，并自动更新 .ssh/config 文件
3. /root 为持久目录，存放在其他位置的文件不保证可以持久保存
4. /shared 为共享目录，与其他用户共享访问
5. 建议使用 tmux 等会话保持工具
" > /root/README

# restore .bashrc .profile
cp -f /etc/skel/.bashrc /etc/skel/.profile /root/

# create /root/.ssh
mkdir -p /root/.ssh
chmod 700 /root/.ssh
cd /root/.ssh

# create id_rsa
ssh-keygen -f /root/.ssh/id_rsa -t rsa -N ''

# write README
echo "id_rsa 和 id_rsa.pub 受 Bunker 管理，请勿修改" > README
`

const tplSSHConfig = `#!/bin/bash
# remove .ssh/config
rm -f /root/.ssh/config

# create new .ssh/config
{{if .Entries}}
{{range .Entries}}
echo "Host {{.Name}}" >> /root/.ssh/config
echo "HostName {{.Host}}" >> /root/.ssh/config
echo "Port {{.Port}}" >> /root/.ssh/config
echo "User {{.User}}" >> /root/.ssh/config
echo "" >> /root/.ssh/config
{{end}}
{{else}}
echo "" > /root/.ssh/config
{{end}}
`

// SSHEntry a entry in ssh_config
type SSHEntry struct {
	Name string
	Host string
	Port uint
	User string
}

func createScript(name string, tmpl string, data map[string]interface{}) string {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}
	buf := &bytes.Buffer{}
	t.Execute(buf, data)
	return buf.String()
}

// ScriptSeedSSHConfig create a script for seeding .ssh/config
func ScriptSeedSSHConfig(entries []SSHEntry) string {
	return createScript(
		"seed-ssh-config",
		tplSSHConfig,
		map[string]interface{}{
			"Entries": entries,
		},
	)
}
