package gaode

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
)

// 高德API配置（替换成你的Key）
const (
	amapAPIKey    = "你的高德API Key"
	amapGeocodeURL = "https://restapi.amap.com/v3/geocode/geo"    // 正向地理编码
	amapRegeoURL  = "https://restapi.amap.com/v3/geocode/regeo"   // 逆地理编码
)

// 坐标类型定义
type CoordType string

const (
	CoordTypeWGS84  CoordType = "wgs84"  // GPS原始坐标
	CoordTypeGCJ02  CoordType = "gcj02"  // 高德/谷歌中国坐标（火星坐标）
	CoordTypeBD09   CoordType = "bd09"   // 百度坐标
)

// 地理编码结果结构体
type GeocodeResult struct {
	Status    int      `json:"status"`    // 状态码 1=成功 0=失败
	Province  string   `json:"province"`  // 省份
	City      string   `json:"city"`      // 城市
	District  string   `json:"district"`  // 区县
	Address   string   `json:"address"`   // 标准化地址
	Lng       float64  `json:"lng"`       // 经度
	Lat       float64  `json:"lat"`       // 纬度
	Formatted string   `json:"formatted"` // 格式化完整地址
}

// 逆地理编码结果结构体
type RegeoResult struct {
	Status    int      `json:"status"`    // 状态码 1=成功 0=失败
	Province  string   `json:"province"`  // 省份
	City      string   `json:"city"`      // 城市
	District  string   `json:"district"`  // 区县
	Street    string   `json:"street"`    // 街道
	Number    string   `json:"number"`    // 门牌号
	Business  string   `json:"business"`  // 商圈
	Address   string   `json:"address"`   // 完整地址
}

// -------------------------- 核心功能：正向地理编码（地址→经纬度） --------------------------
func Geocode(address string) (*GeocodeResult, error) {
	// 1. 地址预处理（简单标准化）
	standardAddr := standardizeAddress(address)

	// 2. 构造请求参数
	params := url.Values{}
	params.Set("key", amapAPIKey)
	params.Set("address", standardAddr)
	params.Set("output", "json")

	// 3. 发送请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(fmt.Sprintf("%s?%s", amapGeocodeURL, params.Encode()))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("请求高德API失败: %v", err)
	}

	// 4. 解析响应
	respBody := string(resp.Body())
	result := &GeocodeResult{}

	// 检查状态
	status := gjson.Get(respBody, "status").Int()
	if status != 1 {
		errMsg := gjson.Get(respBody, "info").String()
		return nil, fmt.Errorf("地理编码失败: %s", errMsg)
	}

	// 提取核心数据
	geo := gjson.Get(respBody, "geocodes.0")
	if !geo.Exists() {
		return nil, fmt.Errorf("未找到该地址的坐标信息")
	}

	result.Status = 1
	result.Province = geo.Get("province").String()
	result.City = geo.Get("city").String()
	result.District = geo.Get("district").String()
	result.Address = geo.Get("formatted_address").String()
	result.Formatted = geo.Get("formatted_address").String()
	
	// 解析经纬度
	location := strings.Split(geo.Get("location").String(), ",")
	if len(location) != 2 {
		return nil, fmt.Errorf("经纬度解析失败")
	}
	lng, _ := fmt.ParseFloat(location[0], 64)
	lat, _ := fmt.ParseFloat(location[1], 64)
	result.Lng = lng
	result.Lat = lat

	return result, nil
}

// -------------------------- 核心功能：逆地理编码（经纬度→地址） --------------------------
func Regeocode(lng, lat float64, coordType CoordType) (*RegeoResult, error) {
	// 1. 构造请求参数
	params := url.Values{}
	params.Set("key", amapAPIKey)
	params.Set("location", fmt.Sprintf("%f,%f", lng, lat))
	params.Set("output", "json")
	params.Set("coordtype", string(coordType)) // 指定输入坐标类型

	// 2. 发送请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(fmt.Sprintf("%s?%s", amapRegeoURL, params.Encode()))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("请求高德API失败: %v", err)
	}

	// 3. 解析响应
	respBody := string(resp.Body())
	result := &RegeoResult{}

	// 检查状态
	status := gjson.Get(respBody, "status").Int()
	if status != 1 {
		errMsg := gjson.Get(respBody, "info").String()
		return nil, fmt.Errorf("逆地理编码失败: %s", errMsg)
	}

	// 提取核心数据
	regeo := gjson.Get(respBody, "regeocode")
	if !regeo.Exists() {
		return nil, fmt.Errorf("未找到该坐标的地址信息")
	}

	result.Status = 1
	result.Address = regeo.Get("formatted_address").String()
	
	// 行政区域信息
	addressComponent := regeo.Get("addressComponent")
	result.Province = addressComponent.Get("province").String()
	result.City = addressComponent.Get("city").String()
	result.District = addressComponent.Get("district").String()
	result.Street = addressComponent.Get("street").String()
	result.Number = addressComponent.Get("number").String()
	
	// 商圈信息
	businessAreas := regeo.Get("addressComponent.businessAreas.0.name").String()
	result.Business = businessAreas

	return result, nil
}

// -------------------------- 核心功能：坐标体系转换（WGS84 ↔ GCJ02 ↔ BD09） --------------------------
// 注：国内地图平台默认返回GCJ02，这里实现常用的转换算法（火星坐标/百度坐标纠偏）
const (
	xPi    = 3.14159265358979324 * 3000.0 / 180.0
	pi     = 3.1415926535897932384626
	a      = 6378245.0
	ee     = 0.00669342162296594323
)

// WGS84转GCJ02（GPS坐标转火星坐标）
func WGS84ToGCJ02(wgsLng, wgsLat float64) (gcjLng, gcjLat float64) {
	if outOfChina(wgsLng, wgsLat) {
		return wgsLng, wgsLat
	}
	dLat := transformLat(wgsLng-105.0, wgsLat-35.0)
	dLng := transformLng(wgsLng-105.0, wgsLat-35.0)
	radLat := wgsLat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * pi)
	dLng = (dLng * 180.0) / (a / sqrtMagic * math.Cos(radLat) * pi)
	gcjLat = wgsLat + dLat
	gcjLng = wgsLng + dLng
	return
}

// GCJ02转BD09（火星坐标转百度坐标）
func GCJ02ToBD09(gcjLng, gcjLat float64) (bdLng, bdLat float64) {
	z := math.Sqrt(gcjLng*gcjLng + gcjLat*gcjLat) + 0.00002*math.Sin(gcjLat*xPi)
	theta := math.Atan2(gcjLat, gcjLng) + 0.000003*math.Cos(gcjLng*xPi)
	bdLng = z*math.Cos(theta) + 0.0065
	bdLat = z*math.Sin(theta) + 0.006
	return
}

// 辅助函数：判断是否在中国境内（境外无需纠偏）
func outOfChina(lng, lat float64) bool {
	return !(lng > 73.66 && lng < 135.05 && lat > 3.86 && lat < 53.55)
}

// 辅助函数：纬度纠偏
func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*pi) + 40.0*math.Sin(y/3.0*pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*pi) + 320*math.Sin(y*pi/30.0)) * 2.0 / 3.0
	return ret
}

// 辅助函数：经度纠偏
func transformLng(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*pi) + 40.0*math.Sin(x/3.0*pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*pi) + 300.0*math.Sin(x/30.0*pi)) * 2.0 / 3.0
	return ret
}

// -------------------------- 核心功能：地址标准化/纠错/补全 --------------------------
func standardizeAddress(address string) string {
	// 1. 去除多余空格、特殊字符
	address = strings.TrimSpace(address)
	reg := regexp.MustCompile(`\s+`)
	address = reg.ReplaceAllString(address, "")

	// 2. 常见错别字/简写纠正
	correctMap := map[string]string{
		"北京市": "北京市",
		"上海":   "上海市",
		"广东省": "广东省",
		"深圳":   "深圳市",
		"杭州市": "杭州市",
		"XX区":   "", // 可扩展更多纠错规则
		" Rd":   "路",
		" St":   "街",
	}
	for wrong, right := range correctMap {
		address = strings.ReplaceAll(address, wrong, right)
	}

	// 3. 补全省份前缀（示例：如果地址以"朝阳区"开头，补全"北京市"）
	if strings.HasPrefix(address, "朝阳区") && !strings.Contains(address, "北京市") {
		address = "北京市" + address
	}

	return address
}
