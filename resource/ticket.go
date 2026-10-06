package resource

// ticket.go — bili_ticket 获取（2026-10-06 -352 三层对抗补丁第三层）
// 协议: bilibili-API-collect docs/login/ticket.md (HMAC-SHA256 流程, key 公开)
// 直播间接口自 2023-10 起强制校验 bili_ticket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ticketGenURL    = "https://api.bilibili.com/bapis/bilibili.api.ticket.v1.Ticket/GenWebTicket"
	ticketKeyIDNew  = "ed04" // 新 cookie (含 bili_ticket 获取能力)
	ticketHMACKeyV2 = "XgOnmc5C1RqIumIKvD47Bd1rLqtklnREBsgknkNE"
)

type ticketRsp struct {
	Code int `json:"code"`
	Data struct {
		Ticket    string `json:"ticket"`
		CreatedAt int64  `json:"created_at"`
		N_expires int64  `json:"expire"`
	} `json:"data"`
}

// GetBiliTicket 获取 bili_ticket（带 cookie 时 90 天有效，匿名 3 小时）
func (a *API) GetBiliTicket() (string, error) {
	if a.biliTicket != "" && time.Since(a.ticketAt) < 12*time.Hour {
		return a.biliTicket, nil // 缓存 12h 内复用
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	payload := url.Values{}
	payload.Set("key_id", ticketKeyIDNew)
	payload.Set("hmac_str", hmacSign(now, ticketHMACKeyV2))
	payload.Set("context_str", "")
	payload.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))

	req, err := http.NewRequest("POST", ticketGenURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	ck := a.BuvidCookieHeader()
	if a.cookie != "" {
		ck = a.cookie + "; " + ck
	}
	req.Header.Set("Cookie", ck)
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var r ticketRsp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if r.Code != 0 || r.Data.Ticket == "" {
		return "", fmt.Errorf("bili_ticket gen failed: code=%d", r.Code)
	}
	a.biliTicket = r.Data.Ticket
	a.ticketAt = time.Now()
	return a.biliTicket, nil
}

// hmacSign: HMAC-SHA256( chr(0)+start+end+nonce, key ) hex
func hmacSign(now, key string) string {
	full := "\x00" + now + "\x00" + now + "\x00" + "ffffffff"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(full))
	return hex.EncodeToString(mac.Sum(nil))
}
