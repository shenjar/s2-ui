package service

import (
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shenaba/2s-ui/logger"
	"github.com/shenaba/2s-ui/util/common"
)

// 证书申请统一走系统 acme.sh,证书装到与 s-ui.sh 脚本一致的 /root/cert/{域名}/,
// 便于脚本 / nginx / sing-box 按固定文件名(fullchain.pem / privkey.pem)复用。
const (
	certBaseDir   = "/root/cert"      // 与脚本完全一致
	acmeIssueTO   = 180 * time.Second // 申请/安装超时
	acmeInstallTO = 120 * time.Second // acme.sh / socat 安装超时
	cmdDetectTO   = 5 * time.Second   // 检测类命令超时
	// systemd 服务的 PATH 可能不全(且不继承登录 shell),补一个兜底,
	// 确保 exec 调用能定位 nginx/socat/systemctl/apt 等。
	fallbackPath = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
)

// acmeIssuing 保证同一时刻只有一个证书申请在跑,避免前端重复点击导致并发申请、
// 撞 Let's Encrypt 限速。
var acmeIssuing sync.Mutex

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

// withHome 返回把 HOME 固定为指定值、并补全 PATH 兜底的环境变量(去重)。
// systemd 服务即使以 root 运行也往往不设 HOME、PATH 也可能不全。
func withHome(home string) []string {
	env := os.Environ()
	out := make([]string, 0, len(env)+2)
	hasPath := false
	for _, e := range env {
		switch {
		case strings.HasPrefix(e, "HOME="):
			continue
		case strings.HasPrefix(e, "PATH="):
			hasPath = true
			e = e + ":" + fallbackPath // 合并兜底路径,重复项无害
		}
		out = append(out, e)
	}
	out = append(out, "HOME="+home)
	if !hasPath {
		out = append(out, "PATH="+fallbackPath)
	}
	return out
}

// runCmd 执行外部命令(HOME 固定为 home),合并 stdout/stderr,超时或非零退出码
// 都包成 error,并把输出原文附在错误里回传前端,便于排查(80 端口被占、域名未解析等)。
func runCmd(timeout time.Duration, home, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = withHome(home)
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

// resolveAcmeSh 在常见位置查找已安装的 acme.sh,返回可执行路径及其对应的 HOME。
// systemd 服务下 HOME 常为空,acme.sh 会装到 /.acme.sh;s-ui.sh 脚本则装到
// /root/.acme.sh。都纳入探测,避免硬编码单一路径导致"找不到"。
func resolveAcmeSh() (bin, home string) {
	candidates := make([]string, 0, 3)
	if h := os.Getenv("HOME"); h != "" {
		candidates = append(candidates, filepath.Join(h, ".acme.sh", "acme.sh"))
	}
	candidates = append(candidates, "/root/.acme.sh/acme.sh", "/.acme.sh/acme.sh")
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			// home = 可执行文件上两级目录:/x/.acme.sh/acme.sh -> /x
			return p, filepath.Dir(filepath.Dir(p))
		}
	}
	return "", ""
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
	if out, err := runCmd(cmdDetectTO, "/root", "systemctl", "is-active", "nginx"); err == nil && strings.TrimSpace(out) == "active" {
		status.Active = true
		status.Installed = true
	}
	return status
}

// ensureAcmeSh 确保 acme.sh 可用,返回其可执行路径与对应 HOME;缺失时自动安装。
// 安装时在 shell 内显式 export HOME=/root,确保装到 /root/.acme.sh(systemd 下
// HOME 常为空,否则会装到 /.acme.sh);并 curl/wget 自适应(最小化系统可能只有其一)。
func ensureAcmeSh() (bin, home string, err error) {
	if bin, home = resolveAcmeSh(); bin != "" {
		return bin, home, nil
	}
	logger.Info("acme.sh 未安装,开始自动安装...")
	installScript := "export HOME=/root; " +
		"if command -v curl >/dev/null 2>&1; then curl https://get.acme.sh | sh; " +
		"else wget -O - https://get.acme.sh | sh; fi"
	out, e := runCmd(acmeInstallTO, "/root", "sh", "-c", installScript)
	if e != nil {
		return "", "", common.NewErrorf("自动安装 acme.sh 失败,请在服务器手动安装(s-ui 脚本 SSL 菜单)。详情:\n%s", out)
	}
	if bin, home = resolveAcmeSh(); bin != "" {
		logger.Info("acme.sh 安装成功:", bin)
		return bin, home, nil
	}
	return "", "", common.NewErrorf("acme.sh 安装后仍未找到(已查 $HOME/.acme.sh、/root/.acme.sh、/.acme.sh)。安装输出:\n%s", out)
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
			if _, err := runCmd(acmeInstallTO, "/root", m[0], m[1:]...); err == nil {
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
	// 同一时刻只允许一个申请,避免重复点击并发申请撞 Let's Encrypt 限速
	if !acmeIssuing.TryLock() {
		return nil, common.NewError("已有证书申请正在进行,请等当前申请完成后再试")
	}
	defer acmeIssuing.Unlock()

	bin, home, err := ensureAcmeSh()
	if err != nil {
		return nil, err
	}

	// 默认 CA 设为 Let's Encrypt(与脚本一致)
	if out, err := runCmd(cmdDetectTO, home, bin, "--set-default-ca", "--server", "letsencrypt"); err != nil {
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
		// standalone 需独占 80;先预检,避免被 nginx 等占用时把底层 bind 报错丢给用户
		if ln, e := net.Listen("tcp", ":80"); e != nil {
			return nil, common.NewErrorf("80 端口被占用,无法用 standalone 申请。若已部署 nginx 请改选「已部署 Nginx」;否则请先停止占用 80 端口的服务。(%v)", e)
		} else {
			_ = ln.Close()
		}
		ensureSocat()
		issueArgs = append(issueArgs, "--standalone", "--httpport", "80")
	}
	if out, err := runCmd(acmeIssueTO, home, bin, issueArgs...); err != nil {
		return nil, common.NewErrorf("证书申请失败(若域名已有未到期证书可加 --force 强制续期):\n%s", out)
	}

	// 安装证书到 /root/cert/{域名}/
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
	}
	// 仅在 reloadcmd 非空时添加。无 nginx 模式【不配】reloadcmd:acme.sh 在首次
	// --installcert 时会立即执行一次 reloadcmd,若是 restart s-ui 会把正在处理本次
	// 申请请求的面板进程杀掉,导致前端拿不到结果(即便证书已申请成功)。无 nginx 续期
	// 后需重启面板加载新证书(前端流程本就有"重启面板"步骤)。
	if rc := a.buildReloadCmd(useNginx); rc != "" {
		installArgs = append(installArgs, "--reloadcmd", rc)
	}
	if out, err := runCmd(acmeIssueTO, home, bin, installArgs...); err != nil {
		return nil, common.NewErrorf("安装证书失败:\n%s", out)
	}

	// 启用 acme.sh 自带 cron 自动续期(失败不影响本次证书)
	if _, err := runCmd(acmeInstallTO, home, bin, "--upgrade", "--auto-upgrade"); err != nil {
		logger.Warning("启用 acme.sh 自动续期失败(不影响本次证书):", err)
	}

	return &IssueResult{CertFile: certFile, KeyFile: keyFile, NginxMode: useNginx}, nil
}

// buildReloadCmd 决定续期成功后的重载命令。
//   - 有 nginx:reload nginx 让其重读证书,不影响面板进程,可安全作为 reloadcmd。
//   - 无 nginx:返回空——面板证书是启动时一次性加载,只能靠重启进程生效,但 reloadcmd
//     在首次申请时会内联触发(见 IssueWeb 注释),会杀掉当前请求,故不配;续期后需手动重启。
func (a *AcmeService) buildReloadCmd(useNginx bool) string {
	if useNginx {
		return "systemctl reload nginx"
	}
	return ""
}

// ListCerts 返回 acme.sh 已管理的证书列表,供前端展示、避免重复申请触发 LE 限速。
func (a *AcmeService) ListCerts() (string, error) {
	if runtime.GOOS == "windows" {
		return "", common.NewError("Windows 不支持 acme.sh")
	}
	bin, home := resolveAcmeSh()
	if bin == "" {
		return "", nil
	}
	out, _ := runCmd(cmdDetectTO, home, bin, "--list")
	return out, nil
}
