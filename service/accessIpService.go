package service

import (
	"errors"
	"sort"
	"strings"

	aDto "hy.juck.com/go-publisher-server/dto/auth"
	"hy.juck.com/go-publisher-server/model/auth"
	"hy.juck.com/go-publisher-server/model/server"
)

type AccessIpService struct {
}

func NewAccessIpService() *AccessIpService {
	return &AccessIpService{}
}

// 校验认证信息
func (obj *AccessIpService) CheckAuthInfo(authDto aDto.AuthRequest) error {
	// 主要校验项目id和服务器唯一标识是否正确，并且服务器是启用状态才行
	var serverInfo server.ServerInfo
	G.DB.Debug().Where("uid = ? and project_id = ? and state = 1", authDto.Uid, authDto.ProjectId).First(&serverInfo)
	if serverInfo.ID == 0 {
		return errors.New("服务器信息校验失败，请重试")
	}
	// 校验该ip是否已经存在不包含当前serverUid和projectList的ip列表中
	var accessIpWhiteList []auth.AccessIpWhite
	G.DB.Debug().Where("project_id != ? and server_uid != ?", authDto.ProjectId, authDto.Uid).Find(&accessIpWhiteList)
	var ipMap = make(map[string]bool, 1)
	for _, accessIpWhite := range accessIpWhiteList {
		ipList := strings.Split(accessIpWhite.IpList, ",")
		for _, ip := range ipList {
			ipMap[ip] = true
		}
	}
	for _, serverIp := range authDto.ServerIpList {
		if _, ok := ipMap[serverIp]; ok {
			return errors.New("服务器信息校验失败，部分或全部ip重复上报")
		}
	}
	return nil
}

func (obj *AccessIpService) SaveAuthIp(authDto aDto.AuthRequest) error {
	// 没有则存储到数据库中，有则更新ip列表
	var accessIpWhite auth.AccessIpWhite
	G.DB.Debug().Where("server_uid = ? and project_id = ?", authDto.Uid, authDto.ProjectId).First(&accessIpWhite)
	// 对上传的ip列表进行去重
	ipMap := make(map[string]bool, 1)
	var ipArr []string
	for _, ip := range authDto.ServerIpList {
		if !ipMap[ip] {
			ipMap[ip] = true
			ipArr = append(ipArr, ip)
		}
	}
	sort.Strings(ipArr)
	ipListStr := strings.Join(ipArr, ",")
	if accessIpWhite.ID == 0 {
		var insertAccessIpWhite = &auth.AccessIpWhite{
			ProjectId: authDto.ProjectId,
			ServerUId: authDto.Uid,
			IpList:    ipListStr,
		}
		if err := G.DB.Debug().Create(insertAccessIpWhite).Error; err != nil {
			return errors.New("插入IP数据失败")
		}
	} else {
		// 直接更新
		if err := G.DB.Debug().Model(&auth.AccessIpWhite{}).Where("id=?", accessIpWhite.ID).Update("ip_list", ipListStr).Error; err != nil {
			return errors.New("更新IP数据失败")
		}
	}
	return nil
}
