package service

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shenaba/2s-ui/logger"
	"github.com/shenaba/2s-ui/util/common"
)

// 证书申请统一走系统 acme.sh,证书装到与 s-ui.sh 脚本一致的 /root/cert/{域名}/,
// 便于脚本 / nginx / sing-box 按固定文件名(fullchain.pem / privkey.pem)复用。
const (
	acmeShBin     = "/root/.acme.sh/acme.sh" // 与 s-ui.sh 安装路径一致(面板以 root 运行)
	certBaseDir   = "/root/cert"             // 与脚本完全一致
	acmeIssueTO   = 180 * time.Second        // 申请/安装超时
	acmeInstallTO = 120 * time.Second        // acme.sh / socat 安装超时
	cmdDetectTO   = 5 * time.Second          // 检测类命令超时
)

// AcmeService 是无状态工具,不嵌入 SettingService(避免与 ApiService 已嵌入的
// SettingService 产生方法集二义性)。所有入参由调用方传入,不直接读写数据库。
type AcmeService struct{}

type NginxStatus struct {
	Installed bool `json:"installed"`
	Active    bool `json:"active"`
}

type IssueResult struct {
	CertFile  string `json:"certFile"`
	KeyFile   string `json:"keyFile"`
	NginxMode bool   `json:"nginxMode"`
}

// runCmd 执行外部命令,合并 stdout/stderr,超时或非零退出码都包成 error
// 并把输出原文附在错误里回传前端,便于排查(如 80 端口被占、域名未解析等)。
func runCmd(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if ctx.Err() == context.DeadlineExceeded {
		return output, common.NewErrorf("命令超时(%v): %s %s\n%s", timeout, name, strings.Join(args, " "), output)
	}
	if err != nil {
		return output, common.NewErrorf("命令执行失败: %s %s\n%s\n%v", name, strings.Join(args, " "), output, err)
	}
	return output, nil
}

// DetectNginx 检测 nginx 是否安装并运行。Windows 直接返回未安装。
func (a *AcmeService) DetectNginx() NginxStatus {
	status := NginxStatus{}
	if runtime.GOOS == "windows" {
		return status
	}
	if _, err := exec.LookPath("nginx"); err == nil {
		status.Installed = true
	}
	// systemctl is-active 在运行时退出码为 0、输出 "active"
	if out, err := runCmd(cmdDetectTO, "systemctl", "is-active", "nginx"); err == nil && strings.TrimSpace(out) == "active" {
		status.Active = true
		status.Installed = true
	}
	return status
}

// ensureAcmeSh 确保 acme.sh 可用,缺失时自动安装(与 s-ui.sh 一致)。
func ensureAcmeSh() error {
	if _, err := os.Stat(acmeShBin); err == nil {
		return nil
	}
	logger.Info("acme.sh 未安装,开始自动安装...")
	out, err := runCmd(acmeInstallTO, "sh", "-c", "curl https://get.acme.sh | sh")
	if err != nil {
		return common.NewErrorf("自动安装 acme.sh 失败,请在服务器手动安装(s-ui 脚本 SSL 菜单)。详情:\n%s", out)
	}
	if _, err := os.Stat(acmeShBin); err != nil {
		return common.NewErrorf("acme.sh 安装后仍未找到 %s,请手动检查。安装输出:\n%s", acmeShBin, out)
	}
	logger.Info("acme.sh 安装成功")
	return nil
}

// ensureSocat 仅 standalone 申请需要,best-effort 安装,失败不致命(由后续申请报真实错)。
func ensureSocat() {
	if _, err := exec.LookPath("socat"); err == nil {
		return
	}
	logger.Info("socat 未安装,尝试自动安装(standalone 申请需要)...")
	managers := [][]string{
		{"apt", "-y", "install", "socat"},
		{"yum", "-y", "install", "socat"},
		{"dnf", "-y", "install", "socat"},
		{"pacman", "-Sy", "--noconfirm", "socat"},
	}
	for _, m := range managers {
		if _, err := exec.LookPath(m[0]); err == nil {
			if _, err := runCmd(acmeInstallTO, m[0], m[1:]...); err == nil {
				return
			}
		}
	}
	logger.Warning("socat 自动安装失败,若 standalone 申请报错请手动安装 socat")
}

// IssueWeb 为面板申请证书并安装到 /root/cert/{域名}/。
//   - useNginx=false:standalone 模式占用 80 端口申请(需 socat)。
//   - useNginx=true :nginx 模式申请,不抢 80 端口,证书供 nginx 使用。
func (a *AcmeService) IssueWeb(domain, email string, useNginx bool) (*IssueResult, error) {
	if runtime.GOOS == "windows" {
		return nil, common.NewError("Windows 不支持 acme.sh 申请证书")
	}
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return nil, common.NewError("域名不能为空")
	}
	if err := ensureAcmeSh(); err != nil {
		return nil, err
	}

	// 默认 CA 设为 Let's Encrypt(与脚本一致)
	if out, err := runCmd(cmdDetectTO, acmeShBin, "--set-default-ca", "--server", "letsencrypt"); err != nil {
		return nil, common.NewErrorf("设置默认 CA 失败:\n%s", out)
	}

	// 申请证书
	issueArgs := []string{"--issue", "-d", domain}
	if email != "" {
		issueArgs = append(issueArgs, "--accountemail", email)
	}
	if useNginx {
		issueArgs = append(issueArgs, "--nginx")
	} else {
		ensureSocat()
		issueArgs = append(issueArgs, "--standalone", "--httpport", "80")
	}
	if out, err := runCmd(acmeIssueTO, acmeShBin, issueArgs...); err != nil {
		return nil, common.NewErrorf("证书申请失败(若域名已有未到期证书可忽略,或用 --force 强制续期):\n%s", out)
	}

	// 安装证书到 /root/cert/{域名}/,并配置续期后自动重载
	certDir := filepath.Join(certBaseDir, domain)
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return nil, common.NewErrorf("创建证书目录失败 %s: %v", certDir, err)
	}
	keyFile := filepath.Join(certDir, "privkey.pem")
	certFile := filepath.Join(certDir, "fullchain.pem")
	installArgs := []string{
		"--installcert", "-d", domain,
		"--key-file", keyFile,
		"--fullchain-file", certFile,
		"--reloadcmd", a.buildReloadCmd(useNginx),
	}
	if out, err := runCmd(acmeIssueTO, acmeShBin, installArgs...); err != nil {
		return nil, common.NewErrorf("安装证书失败:\n%s", out)
	}

	// 启用 acme.sh 自带 cron 自动续期(失败不影响本次证书)
	if _, err := runCmd(acmeInstallTO, acmeShBin, "--upgrade", "--auto-upgrade"); err != nil {
		logger.Warning("启用 acme.sh 自动续期失败(不影响本次证书):", err)
	}

	return &IssueResult{CertFile: certFile, KeyFile: keyFile, NginxMode: useNginx}, nil
}

// buildReloadCmd 决定续期成功后的重载命令。
//   - 有 nginx:reload nginx 即可让其重读证书。
//   - 无 nginx:面板进程内的证书是启动时一次性加载,需重启面板重新读取。
func (a *AcmeService) buildReloadCmd(useNginx bool) string {
	if useNginx {
		return "systemctl reload nginx"
	}
	return "systemctl restart s-ui"
}

// ListCerts 返回 acme.sh 已管理的证书列表,供前端展示、避免重复申请触发 LE 限速。
func (a *AcmeService) ListCerts() (string, error) {
	if runtime.GOOS == "windows" {
		return "", common.NewError("Windows 不支持 acme.sh")
	}
	if _, err := os.Stat(acmeShBin); err != nil {
		return "", nil
	}
	out, _ := runCmd(cmdDetectTO, acmeShBin, "--list")
	return out, nil
}
