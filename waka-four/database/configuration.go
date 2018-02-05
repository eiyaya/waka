package database

import (
	"sync"

	"github.com/liuhan907/waka/waka-four/proto"
)

// 系统配置
type Configuration struct {
	// 主键
	Id int32 `gorm:"index;primary_key;AUTO_INCREMENT"`
	// 类型
	// ext 附加配置
	//   val1 类型
	//        ios           ios 审核
	//   val2 公告内容
	// notice 公告
	//   val1 类型
	//        gapp          健康游戏公告
	//        four_roll     四张滚动公告
	//   val2 公告内容
	// url 链接
	//   val1 类型
	//        recharge      充值链接
	//   val2 链接内容
	// customer_service 客服
	//   val1 客服名称
	//   val2 客服微信
	Type string
	// 值
	Value1 string `gorm:"column:val1;type:text"`
	Value2 string `gorm:"column:val2;type:text"`
	Value3 string `gorm:"column:val3;type:text"`
	Value4 string `gorm:"column:val4;type:text"`
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	lock             sync.RWMutex
	customerServices []*four_proto.FourWelcome_Customer
	exts             map[string]string
	notices          map[string]string
	urls             map[string]string
)

// 获取附加配置
func GetExts() map[string]string {
	lock.RLock()
	defer lock.RUnlock()
	return exts
}

// 获取公告
func GetNotices() map[string]string {
	lock.RLock()
	defer lock.RUnlock()
	return notices
}

// 获取链接配置
func GetUrls() map[string]string {
	lock.RLock()
	defer lock.RUnlock()
	return urls
}

// 获取客服信息
func GetCustomerServices() []*four_proto.FourWelcome_Customer {
	lock.RLock()
	defer lock.RUnlock()
	return customerServices
}

// 刷新配置
func RefreshConfiguration() error {
	v1, err := getExt()
	if err != nil {
		return err
	}

	v2, err := getNotices()
	if err != nil {
		return err
	}

	v3, err := getUrls()
	if err != nil {
		return err
	}

	v4, err := getCustomerServices()
	if err != nil {
		return err
	}

	lock.Lock()

	exts = v1
	notices = v2
	urls = v3
	customerServices = v4

	lock.Unlock()

	return nil
}

func getExt() (map[string]string, error) {
	return getMap("ext")
}

func getNotices() (map[string]string, error) {
	return getMap("notice")
}

func getUrls() (map[string]string, error) {
	return getMap("url")
}

func getCustomerServices() ([]*four_proto.FourWelcome_Customer, error) {
	var vals []*Configuration
	if err := mysql.Where("type = ?", "customer_service").Find(&vals).Error; err != nil {
		return nil, err
	}
	var result []*four_proto.FourWelcome_Customer
	for _, val := range vals {
		result = append(result, &four_proto.FourWelcome_Customer{
			Name:   val.Value1,
			Wechat: val.Value2,
		})
	}
	return result, nil
}

func getMap(mapType string) (map[string]string, error) {
	var vals []*Configuration
	if err := mysql.Where("type = ?", mapType).Find(&vals).Error; err != nil {
		return nil, err
	}
	r := make(map[string]string, len(vals))
	for _, val := range vals {
		if val.Value1 == "" || val.Value2 == "" {
			continue
		}
		r[val.Value1] = val.Value2
	}
	return r, nil
}
