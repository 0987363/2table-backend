package middleware

import (
	"github.com/0987363/2table-backend/models"
	"errors"
	"net"

	"github.com/gin-gonic/gin"
)

func ClientInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := models.ReadRemoteAddress(c)
		geoInfo := &models.GeoInfo{IP: ip, IsPublicIP: isPublicIP(ip)}
		c.Set(models.MiddwareKeyGeoInfo, geoInfo)

		c.Next()
	}
}

func isPublicIP(ip string) bool {
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return false
	}

	if isPrivateAddress := models.IsPrivateAddress(ipAddress); isPrivateAddress {
		return false
	}

	return true
}

func GetClientInfo(c *gin.Context) *models.GeoInfo {
	if i, ok := c.Get(models.MiddwareKeyGeoInfo); ok {
		return i.(*models.GeoInfo)
	}

	return &models.GeoInfo{}
}

func GetClientInfoWithGeo(c *gin.Context) (*models.GeoInfo, error) {
	var geoInfo *models.GeoInfo
	if i, ok := c.Get(models.MiddwareKeyGeoInfo); ok {
		geoInfo = i.(*models.GeoInfo)
	} else {
		return nil, errors.New("Middleware get geoInfo failed")
	}

	if geoInfo.CountryISOCode != "" {
		return geoInfo, nil
	}
	if !geoInfo.IsPublicIP {
		return nil, errors.New("Invalid IP address")
	}

	geoInfo, err := models.GetGeoInfoWithIP(c, geoInfo.IP, GetGeoUrl(), GetRedisHandler())
	if err != nil {
		return nil, err
	}

	c.Set(models.MiddwareKeyGeoInfo, geoInfo)
	return geoInfo, nil
}
